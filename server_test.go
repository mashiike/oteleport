package oteleport_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/mashiike/go-otlp-helper/otlp"
	"github.com/mashiike/oteleport"
	oteleportpb "github.com/mashiike/oteleport/proto"
	"github.com/stretchr/testify/require"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestServer__Trace(t *testing.T) {
	cfg := oteleport.DefaultServerConfig()
	err := cfg.Load("testdata/default.jsonnet", nil)
	require.NoError(t, err)
	grpcOTLPLis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	cfg.OTLP.GRPC.Listener = grpcOTLPLis
	httpOTLPLis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	cfg.OTLP.HTTP.Enable = oteleport.Pointer(false)
	cfg.API.HTTP.Listener = httpOTLPLis
	cfg.Storage.Location += oteleport.RandomString(12)
	err = cfg.Validate()
	require.NoError(t, err)

	s, err := oteleport.NewServer(cfg)
	require.NoError(t, err)
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer wg.Done()
		err := s.Run(ctx)
		require.ErrorIs(t, err, context.Canceled)
	}()

	// upload trace
	bs, err := os.ReadFile("testdata/trace.json")
	require.NoError(t, err)
	var traces tracepb.TracesData
	require.NoError(t, otlp.UnmarshalJSON(bs, &traces))
	client, err := otlp.NewClient("http://" + cfg.OTLP.GRPC.Address)
	require.NoError(t, err)
	err = client.Start(ctx)
	require.NoError(t, err)

	err = client.UploadTraces(ctx, traces.GetResourceSpans())
	require.NoError(t, err)

	err = client.Stop(ctx)
	require.NoError(t, err)

	// fetch trace
	reqBody := oteleportpb.FetchTracesDataRequest{
		StartTimeUnixNano: 1544712660000000000,
		EndTimeUnixNano:   1544712661000000000,
	}
	body, err := otlp.MarshalJSON(&reqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "http://"+cfg.API.HTTP.Address+"/api/traces/fetch", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var respData oteleportpb.FetchTracesDataResponse
	require.NoError(t, otlp.UnmarshalJSON(respBody, &respData))
	require.Equal(t, "", respData.GetNextCursor())
	require.False(t, respData.GetHasMore())
	actual := &tracepb.TracesData{
		ResourceSpans: respData.GetResourceSpans(),
	}
	acutalJSON, err := otlp.MarshalJSON(actual)
	require.NoError(t, err)
	expectedJSON, err := otlp.MarshalJSON(&traces)
	require.NoError(t, err)
	require.JSONEq(t, string(expectedJSON), string(acutalJSON))

	cancel()
	wg.Wait()
}

func TestServer__Metrics(t *testing.T) {
	cfg := oteleport.DefaultServerConfig()
	err := cfg.Load("testdata/default.jsonnet", nil)
	require.NoError(t, err)
	grpcOTLPLis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	cfg.OTLP.GRPC.Listener = grpcOTLPLis
	httpOTLPLis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	cfg.OTLP.HTTP.Enable = oteleport.Pointer(false)
	cfg.API.HTTP.Listener = httpOTLPLis
	cfg.Storage.Location += oteleport.RandomString(12)
	err = cfg.Validate()
	require.NoError(t, err)

	s, err := oteleport.NewServer(cfg)
	require.NoError(t, err)
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer wg.Done()
		err := s.Run(ctx)
		require.ErrorIs(t, err, context.Canceled)
	}()

	// upload metrics
	bs, err := os.ReadFile("testdata/metrics.json")
	require.NoError(t, err)
	var metrics oteleportpb.MetricsData
	require.NoError(t, otlp.UnmarshalJSON(bs, &metrics))
	client, err := otlp.NewClient("http://" + cfg.OTLP.GRPC.Address)
	require.NoError(t, err)
	err = client.Start(ctx)
	require.NoError(t, err)

	err = client.UploadMetrics(ctx, metrics.GetResourceMetrics())
	require.NoError(t, err)

	err = client.Stop(ctx)
	require.NoError(t, err)

	// fetch metrics
	reqBody := oteleportpb.FetchMetricsDataRequest{
		StartTimeUnixNano: 1544712660000000000,
		EndTimeUnixNano:   1544712661000000000,
	}
	body, err := otlp.MarshalJSON(&reqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "http://"+cfg.API.HTTP.Address+"/api/metrics/fetch", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var respData oteleportpb.FetchMetricsDataResponse
	require.NoError(t, otlp.UnmarshalJSON(respBody, &respData))
	require.Equal(t, "", respData.GetNextCursor())
	require.False(t, respData.GetHasMore())
	actual := &oteleportpb.MetricsData{
		ResourceMetrics: respData.GetResourceMetrics(),
	}
	acutalJSON, err := otlp.MarshalJSON(actual)
	require.NoError(t, err)
	expectedJSON, err := otlp.MarshalJSON(&metrics)
	require.NoError(t, err)
	require.JSONEq(t, string(expectedJSON), string(acutalJSON))

	cancel()
	wg.Wait()
}
