package proto_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mashiike/go-otlp-helper/otlp"
	oteleportpb "github.com/mashiike/oteleport/proto"
	"github.com/stretchr/testify/require"
	metricspb "go.opentelemetry.io/proto/otlp/metrics/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestConvertFlattenSpans(t *testing.T) {
	bs, err := os.ReadFile("testdata/trace.json")
	require.NoError(t, err)
	var td tracepb.TracesData
	require.NoError(t, otlp.UnmarshalJSON(bs, &td))
	actual := oteleportpb.ConvertToFlattenSpans(td.GetResourceSpans())
	require.NotNil(t, actual)
	require.Len(t, actual, 1)
	actualJSON, err := otlp.MarshalJSON(actual[0])
	require.NoError(t, err)
	expected, err := os.ReadFile("testdata/flatten_span.json")
	require.NoError(t, err)
	t.Log("actual:", string(actualJSON))
	t.Log("expected:", string(expected))
	require.JSONEq(t, string(expected), string(actualJSON))

	restoreAcutal := oteleportpb.ConvertFromFlattenSpans(actual)
	require.NotNil(t, restoreAcutal)
	restoreActualJSON, err := otlp.MarshalJSON(&tracepb.TracesData{
		ResourceSpans: restoreAcutal,
	})
	require.NoError(t, err)
	require.JSONEq(t, string(bs), string(restoreActualJSON))
}

func TestConvertFlattenDataPoints(t *testing.T) {
	bs, err := os.ReadFile("testdata/metrics.json")
	require.NoError(t, err)
	var md metricspb.MetricsData
	require.NoError(t, otlp.UnmarshalJSON(bs, &md))
	actual := oteleportpb.ConvertToFlattenDataPoints(md.GetResourceMetrics())
	require.NotNil(t, actual)
	require.Len(t, actual, 4)
	for i := 0; i < 4; i++ {
		t.Run(fmt.Sprintf("data_points[%d]", i), func(t *testing.T) {
			actualJSON, err := otlp.MarshalJSON(actual[i])
			require.NoError(t, err)
			t.Log("actual:", string(actualJSON))
			expected, err := os.ReadFile(fmt.Sprintf("testdata/flatten_data_point_%d.json", i))
			require.NoError(t, err)
			t.Log("expected:", string(expected))
			require.JSONEq(t, string(expected), string(actualJSON))
		})
	}
	restoreAcutal := oteleportpb.ConvertFromFlattenDataPoints(actual)
	require.NotNil(t, restoreAcutal)
	restoreActualJSON, err := otlp.MarshalJSON(&metricspb.MetricsData{
		ResourceMetrics: restoreAcutal,
	})
	require.NoError(t, err)
	t.Log("restore:", string(restoreActualJSON))
	restoreExpected, err := otlp.MarshalJSON(&md)
	require.NoError(t, err)
	require.JSONEq(t, string(restoreExpected), string(restoreActualJSON))
}

func TestConvertFlattenLogRecords(t *testing.T) {
	bs, err := os.ReadFile("testdata/logs.json")
	require.NoError(t, err)
	var ld oteleportpb.LogsData
	require.NoError(t, otlp.UnmarshalJSON(bs, &ld))
	actual := oteleportpb.ConvertToFlattenLogRecords(ld.GetResourceLogs())
	require.NotNil(t, actual)
	require.Len(t, actual, 1)

	actualJSON, err := otlp.MarshalJSON(actual[0])
	require.NoError(t, err)
	t.Log("actual:", string(actualJSON))
	expected, err := os.ReadFile("testdata/flatten_log_record.json")
	require.NoError(t, err)
	t.Log("expected:", string(expected))
	require.JSONEq(t, string(expected), string(actualJSON))

	restoreAcutal := oteleportpb.ConvertFromFlattenLogRecords(actual)
	require.NotNil(t, restoreAcutal)
	restoreActualJSON, err := otlp.MarshalJSON(&oteleportpb.LogsData{
		ResourceLogs: restoreAcutal,
	})
	require.NoError(t, err)
	t.Log("restore:", string(restoreActualJSON))
	restoreExpected, err := otlp.MarshalJSON(&ld)
	require.NoError(t, err)
	require.JSONEq(t, string(restoreExpected), string(restoreActualJSON))
}
