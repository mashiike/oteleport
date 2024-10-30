package oteleport

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/mashiike/go-otlp-helper/otlp"
	"github.com/mashiike/oteleport/pkg/client"
	oteleportpb "github.com/mashiike/oteleport/proto"
	"github.com/samber/oops"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

type ClientApp struct {
	c          *client.Client
	outputOpts *ClientSignalOutputOptions
}

type Profile struct {
	*client.Profile
	Output ClientSignalOutputOptions `json:"output"`
}

func (p *Profile) Load(path string, opts *LoadOptions) error {
	vm, err := MakeVM(context.Background())
	if err != nil {
		return oops.Wrapf(err, "failed to make jsonnet vm")
	}
	if opts != nil {
		for k, v := range opts.ExtCodes {
			vm.ExtCode(k, v)
		}
		for k, v := range opts.ExtVars {
			vm.ExtVar(k, v)
		}
	}
	jsonStr, err := vm.EvaluateFile(path)
	if err != nil {
		return oops.Wrapf(err, "failed to evaluate jsonnet file %s", path)
	}
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	dec.DisallowUnknownFields()
	if err := dec.Decode(p); err != nil {
		return oops.Wrapf(err, "failed to decode jsonnet file %s", path)
	}
	return p.Validate()
}

func (p *Profile) Validate() error {
	if err := p.Profile.Validate(); err != nil {
		return err
	}
	return nil
}

func NewClientApp(p *Profile) (*ClientApp, error) {
	c, err := client.New(p.Profile)
	if err != nil {
		return nil, err
	}
	app := &ClientApp{
		c:          c,
		outputOpts: &p.Output,
	}
	return app, nil
}

var (
	followPollingInterval = 5 * time.Second
	fetchPollingInterval  = 200 * time.Millisecond
)

func (a *ClientApp) FetchTracesData(ctx context.Context, opts *ClientTracesCommandOptions) error {
	startTimeUnixNano, endTimeUnixNano := opts.TimeRangeUnixNano()
	var follow bool
	if endTimeUnixNano == 0 {
		follow = true
		endTimeUnixNano = time.Now().UnixNano()
	}
	for {
		slog.DebugContext(ctx, "create pagenator", "start_time", time.Unix(0, startTimeUnixNano), "end_time", time.Unix(0, endTimeUnixNano))
		p := client.NewFetchTracesDataPagenator(a.c, &oteleportpb.FetchTracesDataRequest{
			StartTimeUnixNano: uint64(startTimeUnixNano),
			EndTimeUnixNano:   uint64(endTimeUnixNano),
			Limit:             100,
		})
		for p.HasMorePages() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			resp, err := p.NextPage(ctx)
			if err != nil {
				return err
			}
			if otlp.TotalSpans(resp.GetResourceSpans()) == 0 {
				slog.DebugContext(ctx, "no more spans available")
				continue
			}
			if a.outputOpts.OtelExporterOTLPEndpoint != "" {
				return oops.Errorf("signal export to otel exporter is not implemented yet")
			}
			tracesData := &tracepb.TracesData{
				ResourceSpans: resp.GetResourceSpans(),
			}

			bs, err := otlp.MarshalJSON(tracesData)
			if err != nil {
				slog.WarnContext(ctx, "failed to marshal fetch traces data response", "message", err.Error())
				return nil
			}
			fmt.Println(string(bs))
			time.Sleep(fetchPollingInterval)
		}
		if !follow {
			break
		}
		slog.DebugContext(ctx, "wait for next fetch", "util", time.Now().Add(followPollingInterval))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(followPollingInterval):
		}
		startTimeUnixNano = endTimeUnixNano + 1
		endTimeUnixNano = time.Now().UnixNano()
	}
	return nil
}

func (a *ClientApp) FetchMetricsData(ctx context.Context, opts *ClientMetricsCommandOptions) error {
	startTimeUnixNano, endTimeUnixNano := opts.TimeRangeUnixNano()
	var follow bool
	if endTimeUnixNano == 0 {
		follow = true
		endTimeUnixNano = time.Now().UnixNano()
	}
	for {
		slog.DebugContext(ctx, "create pagenator", "start_time", time.Unix(0, startTimeUnixNano), "end_time", time.Unix(0, endTimeUnixNano))
		p := client.NewFetchMetricsDataPagenator(a.c, &oteleportpb.FetchMetricsDataRequest{
			StartTimeUnixNano: uint64(startTimeUnixNano),
			EndTimeUnixNano:   uint64(endTimeUnixNano),
			Limit:             100,
		})
		for p.HasMorePages() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			resp, err := p.NextPage(ctx)
			if err != nil {
				return err
			}
			if otlp.TotalDataPoints(resp.GetResourceMetrics()) == 0 {
				slog.DebugContext(ctx, "no more metrics available")
				continue
			}
			if a.outputOpts.OtelExporterOTLPEndpoint != "" {
				return oops.Errorf("signal export to otel exporter is not implemented yet")
			}
			metricsData := &oteleportpb.MetricsData{
				ResourceMetrics: resp.GetResourceMetrics(),
			}

			bs, err := otlp.MarshalJSON(metricsData)
			if err != nil {
				slog.WarnContext(ctx, "failed to marshal fetch metrics data response", "message", err.Error())
			}

			fmt.Println(string(bs))
			time.Sleep(fetchPollingInterval)
		}
		if !follow {
			break
		}
		slog.DebugContext(ctx, "wait for next fetch", "util", time.Now().Add(followPollingInterval))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(followPollingInterval):
		}
		startTimeUnixNano = endTimeUnixNano + 1
		endTimeUnixNano = time.Now().UnixNano()
	}
	return nil
}

func (a *ClientApp) FetchLogsData(ctx context.Context, opts *ClientLogsCommandOptions) error {
	startTimeUnixNano, endTimeUnixNano := opts.TimeRangeUnixNano()
	var follow bool
	if endTimeUnixNano == 0 {
		follow = true
		endTimeUnixNano = time.Now().UnixNano()
	}
	for {
		slog.DebugContext(ctx, "create pagenator", "start_time", time.Unix(0, startTimeUnixNano), "end_time", time.Unix(0, endTimeUnixNano))
		p := client.NewFetchLogsDataPagenator(a.c, &oteleportpb.FetchLogsDataRequest{
			StartTimeUnixNano: uint64(startTimeUnixNano),
			EndTimeUnixNano:   uint64(endTimeUnixNano),
			Limit:             100,
		})
		for p.HasMorePages() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			resp, err := p.NextPage(ctx)
			if err != nil {
				return err
			}
			if otlp.TotalLogRecords(resp.GetResourceLogs()) == 0 {
				slog.DebugContext(ctx, "no more logs available")
				continue
			}
			if a.outputOpts.OtelExporterOTLPEndpoint != "" {
				return oops.Errorf("signal export to otel exporter is not implemented yet")
			}
			logsData := &oteleportpb.LogsData{
				ResourceLogs: resp.GetResourceLogs(),
			}

			bs, err := otlp.MarshalJSON(logsData)
			if err != nil {
				slog.WarnContext(ctx, "failed to marshal fetch logs data response", "message", err.Error())
			}

			fmt.Println(string(bs))
			time.Sleep(fetchPollingInterval)
		}
		if !follow {
			break
		}
		slog.DebugContext(ctx, "wait for next fetch", "util", time.Now().Add(followPollingInterval))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(followPollingInterval):
		}
		startTimeUnixNano = endTimeUnixNano + 1
		endTimeUnixNano = time.Now().UnixNano()
	}
	return nil
}
