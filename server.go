package oteleport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fujiwara/ridge"
	"github.com/gorilla/mux"
	"github.com/mashiike/go-otlp-helper/otlp"
	oteleportpb "github.com/mashiike/oteleport/proto"
	"github.com/samber/oops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	otlpMux     *otlp.ServerMux
	apiMux      *mux.Router
	cfg         *ServerConfig
	signalRepo  SignalRepository
	TermHandler func()
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	s := &Server{
		otlpMux: otlp.NewServerMux(),
		apiMux:  mux.NewRouter(),
		cfg:     cfg,
	}
	repo, err := NewSignalRepository(&cfg.Storage)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to create signal repository")
	}
	s.signalRepo = repo
	s.setupOTLP()
	s.setupAPI()
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	if ridge.AsLambdaHandler() {
		return s.runAsLambdaHandler(ctx)
	}
	return s.runOnLocalServer(ctx)
}

func (s *Server) setupOTLP() {
	s.otlpMux.Trace().HandleFunc(s.handleTraces)
	s.otlpMux.Metrics().HandleFunc(s.handleMetrics)
	s.otlpMux.Logs().HandleFunc(s.handleLogs)
	s.otlpMux.Use(
		func(next otlp.ProtoHandlerFunc) otlp.ProtoHandlerFunc {
			return func(ctx context.Context, req proto.Message) (proto.Message, error) {
				var signalType string
				switch req.(type) {
				case *otlp.TraceRequest:
					signalType = "trace"
				case *otlp.MetricsRequest:
					signalType = "metrics"
				case *otlp.LogsRequest:
					signalType = "logs"
				}
				slog.Info("received otlp telemetry", "type", signalType)
				return next(ctx, req)
			}
		},
	)
	if s.cfg.EnableAuth() {
		s.otlpMux.Use(func(next otlp.ProtoHandlerFunc) otlp.ProtoHandlerFunc {
			return func(ctx context.Context, req proto.Message) (proto.Message, error) {
				header, ok := otlp.HeadersFromContext(ctx)
				if !ok {
					return nil, status.Error(codes.Unauthenticated, "no metadata found")
				}
				accessKey := header.Get(s.cfg.AccessKeyHeader)
				if accessKey == "" {
					slog.InfoContext(ctx, "access denided", "reason", "no access key found")
					return nil, status.Error(codes.Unauthenticated, "no access key found")
				}
				for _, key := range s.cfg.AccessKeys {
					if key.SecretKey == accessKey {
						slog.InfoContext(ctx, "authenticated", "key_id", key.KeyID)
						return next(ctx, req)
					}
				}
				slog.InfoContext(ctx, "access denided", "reason", "access key mismatch")
				return nil, status.Error(codes.PermissionDenied, "access denied")
			}
		})
	}
}

const (
	apiPathPrefix    = "/api"
	fetchTracesPath  = "/traces/fetch"
	fetchMetricsPath = "/metrics/fetch"
	fetchLogsPath    = "/logs/fetch"
)

func (s *Server) setupAPI() {
	base := s.apiMux.PathPrefix(apiPathPrefix).Subrouter()
	s.apiMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	base.HandleFunc(fetchTracesPath, s.serveFetchTraces)
	base.HandleFunc(fetchMetricsPath, s.serveFetchMetrics)
	base.HandleFunc(fetchLogsPath, s.serveFetchLogs)
	base.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("accept api request", "method", r.Method, "path", r.URL.Path, "content_type", r.Header.Get("Content-Type"))
			next.ServeHTTP(w, r)
		})
	})
	if s.cfg.EnableAuth() {
		base.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				header := r.Header.Get(s.cfg.AccessKeyHeader)
				if header == "" {
					st := status.New(codes.Unauthenticated, "no access key found")
					writeError(w, r, st, http.StatusUnsupportedMediaType)
					return
				}
				for _, key := range s.cfg.AccessKeys {
					if key.SecretKey == header {
						slog.Info("authenticated", "key_id", key.KeyID)
						next.ServeHTTP(w, r)
						return
					}
				}
				st := status.New(codes.PermissionDenied, "access denied")
				writeError(w, r, st, http.StatusForbidden)
			})
		})
	}
}

func (s *Server) runAsLambdaHandler(ctx context.Context) error {
	httpMux := mux.NewRouter()
	httpMux.PathPrefix("/v1").Handler(s.otlpMux)
	httpMux.PathPrefix("/api").Handler(s.apiMux)
	handler := func(ctx context.Context, event json.RawMessage) (interface{}, error) {
		req, err := ridge.NewRequest(event)
		if err != nil {
			slog.Error("failed to build request", "err", err)
			return nil, err
		}
		req = req.WithContext(ctx)
		w := ridge.NewResponseWriter()
		httpMux.ServeHTTP(w, req)
		return w.Response(), nil
	}
	opts := []lambda.Option{lambda.WithContext(ctx)}
	if s.TermHandler != nil {
		opts = append(opts, lambda.WithEnableSIGTERM(s.TermHandler))
	}
	lambda.StartWithOptions(handler, opts...)
	return nil
}

func valueOrDefault[T any](v *T, d T) T {
	if v == nil {
		return d
	}
	return *v
}

func (s *Server) runOnLocalServer(ctx context.Context) error {
	var wg sync.WaitGroup
	cleanups := make([]func(context.Context), 0)
	var onceCleanup sync.Once
	cleanup := func() {
		onceCleanup.Do(func() {
			for _, f := range cleanups {
				f(ctx)
			}
		})
	}
	defer cleanup()
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	if valueOrDefault(s.cfg.OTLP.GRPC.Enable, false) {
		grpcServer := grpc.NewServer()
		s.otlpMux.Register(grpcServer)
		reflection.Register(grpcServer)
		grpcListener := s.cfg.OTLP.GRPC.Listener
		if grpcListener == nil {
			var err error
			grpcListener, err = net.Listen("tcp", s.cfg.OTLP.GRPC.Address)
			if err != nil {
				return oops.Wrapf(err, "failed to listen to %s", s.cfg.OTLP.GRPC.Address)
			}
		}
		cleanups = append(cleanups, startGRPCServer(&wg, ctx, cancel, grpcServer, grpcListener, "otlp"))
	}
	if valueOrDefault(s.cfg.OTLP.HTTP.Enable, false) {
		httpMux := http.NewServeMux()
		httpMux.Handle("/", s.otlpMux)
		httpServer := &http.Server{
			Addr:    s.cfg.OTLP.HTTP.Address,
			Handler: httpMux,
		}
		httpListener := s.cfg.OTLP.HTTP.Listener
		if httpListener == nil {
			var err error
			httpListener, err = net.Listen("tcp", s.cfg.OTLP.HTTP.Address)
			if err != nil {
				return oops.Wrapf(err, "failed to listen to %s", s.cfg.OTLP.HTTP.Address)
			}
		}
		cleanups = append(cleanups, startHTTPServer(&wg, ctx, cancel, httpServer, httpListener, "otlp"))
	}

	if valueOrDefault(s.cfg.API.HTTP.Enable, false) {
		httpMux := http.NewServeMux()
		httpMux.Handle("/", s.apiMux)
		server := &http.Server{
			Addr:    s.cfg.API.HTTP.Address,
			Handler: httpMux,
		}
		httpListener := s.cfg.API.HTTP.Listener
		if httpListener == nil {
			var err error
			httpListener, err = net.Listen("tcp", s.cfg.API.HTTP.Address)
			if err != nil {
				return oops.Wrapf(err, "failed to listen to %s", s.cfg.API.HTTP.Address)
			}
		}
		cleanups = append(cleanups, startHTTPServer(&wg, ctx, cancel, server, httpListener, "api"))
	}
	wg.Add(1)
	go func() {
		<-ctx.Done()
		cleanup()
		cancel(ctx.Err())
		wg.Done()
	}()
	wg.Wait()
	return context.Cause(ctx)
}

func startGRPCServer(wg *sync.WaitGroup, ctx context.Context, cancel context.CancelCauseFunc, server *grpc.Server, listener net.Listener, purpuse string) func(ctx context.Context) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.InfoContext(ctx, fmt.Sprintf("starting %s grpc server", purpuse), "addr", listener.Addr().String())
		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			slog.DebugContext(ctx, "failed to serve", "err", err)
			cancel(fmt.Errorf("failed to serve: %w", err))
		}
	}()
	return func(ctx context.Context) {
		slog.InfoContext(ctx, fmt.Sprintf("shutting down %s grpc server", purpuse), "addr", listener.Addr().String())
		server.Stop()
		if err := listener.Close(); err != nil {
			slog.DebugContext(ctx, "failed to close listener", "err", err.Error())
		}
	}
}

func startHTTPServer(wg *sync.WaitGroup, ctx context.Context, cancel context.CancelCauseFunc, server *http.Server, listener net.Listener, purpuse string) func(ctx context.Context) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.InfoContext(ctx, fmt.Sprintf("starting %s http server", purpuse), "addr", server.Addr)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			slog.DebugContext(ctx, "failed to serve", "err", err)
			cancel(fmt.Errorf("failed to serve: %w", err))
		}
	}()
	return func(ctx context.Context) {
		slog.InfoContext(ctx, fmt.Sprintf("shutting down %s http server", purpuse), "addr", server.Addr)
		sCtx, sCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer sCancel()
		if err := server.Shutdown(sCtx); err != nil {
			slog.DebugContext(ctx, "failed to shutdown server", "err", err.Error())
		}
		if err := listener.Close(); err != nil {
			slog.DebugContext(ctx, "failed to close listener", "err", err.Error())
		}
	}
}

func (s *Server) handleTraces(ctx context.Context, req *otlp.TraceRequest) (*otlp.TraceResponse, error) {
	resourceSpans := req.GetResourceSpans()
	slog.Info("received otlp trace", "total_spans", otlp.TotalSpans(resourceSpans))
	if err := s.signalRepo.PushTracesData(ctx, &oteleportpb.TracesData{
		ResourceSpans: resourceSpans,
		SignalType:    "traces",
	}); err != nil {
		errID := RandomString(16)
		slog.Error("failed to put resource spans", "err_id", errID, "details", err.Error())
		return nil, fmt.Errorf("failed to put resource spans: error_id=%s", errID)
	}
	return &otlp.TraceResponse{}, nil
}

func (s *Server) handleMetrics(ctx context.Context, req *otlp.MetricsRequest) (*otlp.MetricsResponse, error) {
	resourceMetrics := req.GetResourceMetrics()
	slog.Info("received otlp metrics", "total_data_points", otlp.TotalDataPoints(resourceMetrics))
	if err := s.signalRepo.PushMetricsData(ctx, &oteleportpb.MetricsData{
		ResourceMetrics: resourceMetrics,
		SignalType:      "metrics",
	}); err != nil {
		errID := RandomString(16)
		slog.Error("failed to put resource metrics", "err_id", errID, "details", err.Error())
		return nil, fmt.Errorf("failed to put resource metrics: error_id=%s", errID)
	}
	return &otlp.MetricsResponse{}, nil
}

func (s *Server) handleLogs(ctx context.Context, req *otlp.LogsRequest) (*otlp.LogsResponse, error) {
	resourceLogs := req.GetResourceLogs()
	slog.Info("received otlp logs", "total_log_records", otlp.TotalLogRecords(resourceLogs))
	if err := s.signalRepo.PushLogsData(ctx, &oteleportpb.LogsData{
		ResourceLogs: resourceLogs,
		SignalType:   "logs",
	}); err != nil {
		errID := RandomString(16)
		slog.Error("failed to put resource logs", "err_id", errID, "details", err.Error())
		return nil, fmt.Errorf("failed to put resource logs: error_id=%s", errID)
	}
	return &otlp.LogsResponse{}, nil
}

func parseRequest[T proto.Message](r *http.Request, v T) error {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err := otlp.UnmarshalJSON(bs, v); err != nil {
			return err
		}
	case "application/protobuf", "application/x-protobuf":
		if err := proto.Unmarshal(bs, v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported content type: %s", r.Header.Get("Content-Type"))
	}
	return nil
}

func writeError(w http.ResponseWriter, r *http.Request, st *status.Status, code int) {
	switch r.Header.Get("Accept") {
	case "application/json":
		slog.Debug("response content type", "content_type", "application/json")
		writeErrorJSON(w, st, code)
	case "application/protobuf", "application/x-protobuf":
		slog.Debug("response content type", "content_type", "application/protobuf")
		writeErrorProto(w, st, code)
	case "*/*", "":
		slog.Debug("response content type", "content_type", "*/*")
		if r.Header.Get("Content-Type") == "application/json" {
			writeErrorJSON(w, st, code)
		} else {
			writeErrorProto(w, st, code)
		}
	default:
		slog.Debug("response content type", "content_type", r.Header.Get("Accept"))
		writeErrorJSON(w, st, code)
	}
}

func writeErrorProto(w http.ResponseWriter, st *status.Status, code int) {
	bs, err := proto.Marshal(st.Proto())
	if err != nil {
		http.Error(w, http.StatusText(code), code)
	}
	w.Header().Set("Content-Type", "application/protobuf")
	w.WriteHeader(code)
	if _, err := w.Write(bs); err != nil {
		slog.Debug("failed to write response", "error", err.Error())
	}
}

func writeErrorJSON(w http.ResponseWriter, st *status.Status, code int) {
	bs, err := otlp.MarshalJSON(st.Proto())
	if err != nil {
		http.Error(w, http.StatusText(code), code)
	}
	bs = append(bs, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(bs); err != nil {
		slog.Debug("failed to write response", "error", err.Error())
	}
}

func writeResponse(w http.ResponseWriter, r *http.Request, v proto.Message) {
	switch r.Header.Get("Accept") {
	case "application/json":
		slog.Debug("response content type", "content_type", "application/json")
		writeResponseJSON(w, v)
	case "application/protobuf", "application/x-protobuf":
		slog.Debug("response content type", "content_type", "application/protobuf")
		writeResponseProto(w, v)
	case "*/*", "":
		slog.Debug("response content type", "content_type", "*/*")
		if r.Header.Get("Content-Type") == "application/json" {
			writeResponseJSON(w, v)
		} else {
			writeResponseProto(w, v)
		}
	default:
		slog.Debug("response content type", "content_type", r.Header.Get("Accept"))
		writeResponseJSON(w, v)
	}
}

func writeResponseProto(w http.ResponseWriter, v proto.Message) {
	bs, err := proto.Marshal(v)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		writeError(w, nil, st, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/protobuf")
	if _, err := w.Write(bs); err != nil {
		slog.Debug("failed to write response", "error", err.Error())
	}
}

func writeResponseJSON(w http.ResponseWriter, v proto.Message) {
	bs, err := otlp.MarshalJSON(v)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		writeError(w, nil, st, http.StatusInternalServerError)
		return
	}
	bs = append(bs, '\n')
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bs); err != nil {
		slog.Debug("failed to write response", "error", err.Error())
	}

}

var allowedContentTypes = []string{
	"application/json",
	"application/protobuf",
	"application/x-protobuf",
}

func (s *Server) serveFetchTraces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		st := status.New(codes.Unimplemented, "method not allowed")
		writeError(w, r, st, http.StatusMethodNotAllowed)
		return
	}
	if !slices.Contains(allowedContentTypes, r.Header.Get("Content-Type")) {
		st := status.New(codes.InvalidArgument, "unsupported content type")
		writeError(w, r, st, http.StatusUnsupportedMediaType)
		return
	}
	req := &oteleportpb.FetchTracesDataRequest{}
	if err := parseRequest(r, req); err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		writeError(w, r, st, http.StatusBadRequest)
		return
	}
	resp, err := s.signalRepo.FetchTracesData(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			st = status.New(codes.Internal, err.Error())
		}
		writeError(w, r, st, http.StatusInternalServerError)
		return
	}
	writeResponse(w, r, resp)
}

func (s *Server) serveFetchMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		st := status.New(codes.Unimplemented, "method not allowed")
		writeError(w, r, st, http.StatusMethodNotAllowed)
		return
	}
	if !slices.Contains(allowedContentTypes, r.Header.Get("Content-Type")) {
		st := status.New(codes.InvalidArgument, "unsupported content type")
		writeError(w, r, st, http.StatusUnsupportedMediaType)
		return
	}
	req := &oteleportpb.FetchMetricsDataRequest{}
	if err := parseRequest(r, req); err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		writeError(w, r, st, http.StatusBadRequest)
		return
	}
	resp, err := s.signalRepo.FetchMetricsData(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			st = status.New(codes.Internal, err.Error())
		}
		writeError(w, r, st, http.StatusInternalServerError)
		return
	}
	writeResponse(w, r, resp)
}

func (s *Server) serveFetchLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		st := status.New(codes.Unimplemented, "method not allowed")
		writeError(w, r, st, http.StatusMethodNotAllowed)
		return
	}
	if !slices.Contains(allowedContentTypes, r.Header.Get("Content-Type")) {
		st := status.New(codes.InvalidArgument, "unsupported content type")
		writeError(w, r, st, http.StatusUnsupportedMediaType)
		return
	}
	req := &oteleportpb.FetchLogsDataRequest{}
	if err := parseRequest(r, req); err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		writeError(w, r, st, http.StatusBadRequest)
		return
	}
	resp, err := s.signalRepo.FetchLogsData(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			st = status.New(codes.Internal, err.Error())
		}
		writeError(w, r, st, http.StatusInternalServerError)
		return
	}
	writeResponse(w, r, resp)
}
