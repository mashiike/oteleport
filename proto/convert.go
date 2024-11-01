package proto

import (
	"github.com/mashiike/go-otlp-helper/otlp"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func ConvertToFlattenSpans(resourceSpans []*tracepb.ResourceSpans) []*FlattenSpan {
	if resourceSpans == nil {
		return nil
	}
	flattenSpans := make([]*FlattenSpan, 0, len(resourceSpans))
	for _, resourceSpan := range resourceSpans {
		resource := resourceSpan.GetResource()
		resourceAttrs := resource.GetAttributes()
		resourceDroppedAttrsCount := resource.GetDroppedAttributesCount()
		resourceSpanSchemaURL := resourceSpan.GetSchemaUrl()
		for _, scopeSpans := range resourceSpan.GetScopeSpans() {
			scope := scopeSpans.GetScope()
			scopeName := scope.GetName()
			scopeVersion := scope.GetVersion()
			scopeAttrs := scope.GetAttributes()
			scopeDroppedAttrsCount := scope.GetDroppedAttributesCount()
			scopeSpansSchemaURL := scopeSpans.GetSchemaUrl()
			for _, span := range scopeSpans.GetSpans() {
				flattenSpans = append(flattenSpans, &FlattenSpan{
					ResourceAttributes:             resourceAttrs,
					ResourceSpanSchemaUrl:          resourceSpanSchemaURL,
					DroppedResourceAttributesCount: resourceDroppedAttrsCount,
					ScopeName:                      scopeName,
					ScopeVersion:                   scopeVersion,
					ScopeAttributes:                scopeAttrs,
					DroppedScopeAttributesCount:    scopeDroppedAttrsCount,
					ScopeSpanSchemaUrl:             scopeSpansSchemaURL,
					TraceId:                        span.GetTraceId(),
					SpanId:                         span.GetSpanId(),
					ParentSpanId:                   span.GetParentSpanId(),
					Name:                           span.GetName(),
					Kind:                           span.GetKind(),
					StartTimeUnixNano:              span.GetStartTimeUnixNano(),
					EndTimeUnixNano:                span.GetEndTimeUnixNano(),
					Attributes:                     span.GetAttributes(),
					DroppedAttributesCount:         span.GetDroppedAttributesCount(),
					Events:                         span.GetEvents(),
					DroppedEventsCount:             span.GetDroppedEventsCount(),
					Status:                         span.GetStatus(),
					Links:                          span.GetLinks(),
					DroppedLinksCount:              span.GetDroppedLinksCount(),
					Flags:                          span.GetFlags(),
				})
			}
		}
	}
	return flattenSpans
}

func ConvertFromFlattenSpans(fs []*FlattenSpan) []*tracepb.ResourceSpans {
	if fs == nil {
		return nil
	}
	resourceSpans := make([]*tracepb.ResourceSpans, 0)
	for _, f := range fs {
		resource := &resourcepb.Resource{
			Attributes:             f.GetResourceAttributes(),
			DroppedAttributesCount: f.GetDroppedResourceAttributesCount(),
		}
		resourceSpan := &tracepb.ResourceSpans{
			Resource:  resource,
			SchemaUrl: f.GetResourceSpanSchemaUrl(),
		}
		scope := &commonpb.InstrumentationScope{
			Name:                   f.GetScopeName(),
			Version:                f.GetScopeVersion(),
			Attributes:             f.GetScopeAttributes(),
			DroppedAttributesCount: f.GetDroppedScopeAttributesCount(),
		}
		scopeSpan := &tracepb.ScopeSpans{
			Scope:     scope,
			SchemaUrl: f.GetScopeSpanSchemaUrl(),
		}
		span := &tracepb.Span{
			TraceId:                f.GetTraceId(),
			SpanId:                 f.GetSpanId(),
			ParentSpanId:           f.GetParentSpanId(),
			Name:                   f.GetName(),
			Kind:                   f.GetKind(),
			StartTimeUnixNano:      f.GetStartTimeUnixNano(),
			EndTimeUnixNano:        f.GetEndTimeUnixNano(),
			Attributes:             f.GetAttributes(),
			DroppedAttributesCount: f.GetDroppedAttributesCount(),
			Events:                 f.GetEvents(),
			DroppedEventsCount:     f.GetDroppedEventsCount(),
		}
		scopeSpan.Spans = append(scopeSpan.Spans, span)
		resourceSpan.ScopeSpans = append(resourceSpan.ScopeSpans, scopeSpan)
		resourceSpans = otlp.AppendResourceSpans(resourceSpans, resourceSpan)
	}
	return resourceSpans
}

func ConvertToFlattenDataPoints(resourceMetrics []*metricpb.ResourceMetrics) []*FlattenDataPoint {
	if resourceMetrics == nil {
		return nil
	}
	flattenDataPoints := make([]*FlattenDataPoint, 0, len(resourceMetrics))
	for _, resourceMetric := range resourceMetrics {
		resource := resourceMetric.GetResource()
		resourceAttrs := resource.GetAttributes()
		resourceDroppedAttrsCount := resource.GetDroppedAttributesCount()
		resourceMetricSchemaURL := resourceMetric.GetSchemaUrl()
		for _, scopeMetric := range resourceMetric.GetScopeMetrics() {
			scope := scopeMetric.GetScope()
			scopeName := scope.GetName()
			scopeVersion := scope.GetVersion()
			scopeAttrs := scope.GetAttributes()
			scopeDroppedAttrsCount := scope.GetDroppedAttributesCount()
			scopeMetricSchemaURL := scopeMetric.GetSchemaUrl()
			for _, metric := range scopeMetric.GetMetrics() {
				for _, fdp := range convertToFlattenDataPoints(metric) {
					fdp.ResourceAttributes = resourceAttrs
					fdp.ResourceMetricSchemaUrl = resourceMetricSchemaURL
					fdp.DroppedResourceAttributesCount = resourceDroppedAttrsCount
					fdp.ScopeName = scopeName
					fdp.ScopeVersion = scopeVersion
					fdp.ScopeAttributes = scopeAttrs
					fdp.DroppedScopeAttributesCount = scopeDroppedAttrsCount
					fdp.ScopeMetricSchemaUrl = scopeMetricSchemaURL
					flattenDataPoints = append(
						flattenDataPoints,
						fdp,
					)
				}
			}
		}
	}
	return flattenDataPoints
}

func convertToFlattenDataPoints(metric *metricpb.Metric) []*FlattenDataPoint {
	dataPoints := make([]*FlattenDataPoint, 0)
	switch data := metric.GetData().(type) {
	case *metricpb.Metric_Gauge:
		for _, dp := range data.Gauge.GetDataPoints() {
			dataPoints = append(dataPoints, &FlattenDataPoint{
				Name:              metric.GetName(),
				Description:       metric.GetDescription(),
				Unit:              metric.GetUnit(),
				Metadata:          metric.GetMetadata(),
				StartTimeUnixNano: dp.GetStartTimeUnixNano(),
				TimeUnixNano:      dp.GetTimeUnixNano(),
				Data: &FlattenDataPoint_Gauge{
					Gauge: &FlattenGuage{
						DataPoint: dp,
					},
				},
			})
		}
	case *metricpb.Metric_Sum:
		for _, dp := range data.Sum.GetDataPoints() {
			dataPoints = append(dataPoints, &FlattenDataPoint{
				Name:              metric.GetName(),
				Description:       metric.GetDescription(),
				Unit:              metric.GetUnit(),
				Metadata:          metric.GetMetadata(),
				StartTimeUnixNano: dp.GetStartTimeUnixNano(),
				TimeUnixNano:      dp.GetTimeUnixNano(),
				Data: &FlattenDataPoint_Sum{
					Sum: &FlattenSum{
						DataPoint:              dp,
						AggregationTemporality: data.Sum.GetAggregationTemporality(),
						IsMonotonic:            data.Sum.GetIsMonotonic(),
					},
				},
			})
		}
	case *metricpb.Metric_Histogram:
		for _, dp := range data.Histogram.GetDataPoints() {
			dataPoints = append(dataPoints, &FlattenDataPoint{
				Name:              metric.GetName(),
				Description:       metric.GetDescription(),
				Unit:              metric.GetUnit(),
				Metadata:          metric.GetMetadata(),
				StartTimeUnixNano: dp.GetStartTimeUnixNano(),
				TimeUnixNano:      dp.GetTimeUnixNano(),
				Data: &FlattenDataPoint_Histogram{
					Histogram: &FlattenHistogram{
						DataPoint:              dp,
						AggregationTemporality: data.Histogram.GetAggregationTemporality(),
					},
				},
			})
		}
	case *metricpb.Metric_ExponentialHistogram:
		for _, dp := range data.ExponentialHistogram.GetDataPoints() {
			dataPoints = append(dataPoints, &FlattenDataPoint{
				Name:              metric.GetName(),
				Description:       metric.GetDescription(),
				Unit:              metric.GetUnit(),
				Metadata:          metric.GetMetadata(),
				StartTimeUnixNano: dp.GetStartTimeUnixNano(),
				TimeUnixNano:      dp.GetTimeUnixNano(),
				Data: &FlattenDataPoint_ExponentialHistogram{
					ExponentialHistogram: &FlattenExponentialHistogram{
						DataPoint:              dp,
						AggregationTemporality: data.ExponentialHistogram.GetAggregationTemporality(),
					},
				},
			})
		}
	case *metricpb.Metric_Summary:
		for _, dp := range data.Summary.GetDataPoints() {
			dataPoints = append(dataPoints, &FlattenDataPoint{
				Name:              metric.GetName(),
				Description:       metric.GetDescription(),
				Unit:              metric.GetUnit(),
				Metadata:          metric.GetMetadata(),
				StartTimeUnixNano: dp.GetStartTimeUnixNano(),
				TimeUnixNano:      dp.GetTimeUnixNano(),
				Data: &FlattenDataPoint_Summary{
					Summary: &FlattenSummary{
						DataPoint: dp,
					},
				},
			})
		}

	}
	return dataPoints
}

func ConvertFromFlattenDataPoints(fdp []*FlattenDataPoint) []*metricpb.ResourceMetrics {
	if fdp == nil {
		return nil
	}
	resourceMetrics := make([]*metricpb.ResourceMetrics, 0)
	for _, f := range fdp {
		resource := &resourcepb.Resource{
			Attributes:             f.GetResourceAttributes(),
			DroppedAttributesCount: f.GetDroppedResourceAttributesCount(),
		}
		resourceMetric := &metricpb.ResourceMetrics{
			Resource:  resource,
			SchemaUrl: f.GetResourceMetricSchemaUrl(),
		}
		scope := &commonpb.InstrumentationScope{
			Name:                   f.GetScopeName(),
			Version:                f.GetScopeVersion(),
			Attributes:             f.GetScopeAttributes(),
			DroppedAttributesCount: f.GetDroppedScopeAttributesCount(),
		}
		scopeMetric := &metricpb.ScopeMetrics{
			Scope:     scope,
			SchemaUrl: f.GetScopeMetricSchemaUrl(),
		}
		metric := convertFromFlattenDataPoint(f)
		scopeMetric.Metrics = append(scopeMetric.Metrics, metric)
		resourceMetric.ScopeMetrics = append(resourceMetric.ScopeMetrics, scopeMetric)
		resourceMetrics = otlp.AppendResourceMetrics(resourceMetrics, resourceMetric)
	}
	return resourceMetrics
}

func convertFromFlattenDataPoint(fdp *FlattenDataPoint) *metricpb.Metric {
	switch data := fdp.GetData().(type) {
	case *FlattenDataPoint_Gauge:
		return &metricpb.Metric{
			Name:        fdp.GetName(),
			Description: fdp.GetDescription(),
			Unit:        fdp.GetUnit(),
			Metadata:    fdp.GetMetadata(),
			Data: &metricpb.Metric_Gauge{
				Gauge: &metricpb.Gauge{
					DataPoints: []*metricpb.NumberDataPoint{
						data.Gauge.GetDataPoint(),
					},
				},
			},
		}
	case *FlattenDataPoint_Sum:
		return &metricpb.Metric{
			Name:        fdp.GetName(),
			Description: fdp.GetDescription(),
			Unit:        fdp.GetUnit(),
			Metadata:    fdp.GetMetadata(),
			Data: &metricpb.Metric_Sum{
				Sum: &metricpb.Sum{
					AggregationTemporality: data.Sum.GetAggregationTemporality(),
					IsMonotonic:            data.Sum.GetIsMonotonic(),
					DataPoints: []*metricpb.NumberDataPoint{
						data.Sum.GetDataPoint(),
					},
				},
			},
		}
	case *FlattenDataPoint_Histogram:
		return &metricpb.Metric{
			Name:        fdp.GetName(),
			Description: fdp.GetDescription(),
			Unit:        fdp.GetUnit(),
			Metadata:    fdp.GetMetadata(),
			Data: &metricpb.Metric_Histogram{
				Histogram: &metricpb.Histogram{
					AggregationTemporality: data.Histogram.GetAggregationTemporality(),
					DataPoints: []*metricpb.HistogramDataPoint{
						data.Histogram.GetDataPoint(),
					},
				},
			},
		}
	case *FlattenDataPoint_ExponentialHistogram:
		return &metricpb.Metric{
			Name:        fdp.GetName(),
			Description: fdp.GetDescription(),
			Unit:        fdp.GetUnit(),
			Metadata:    fdp.GetMetadata(),
			Data: &metricpb.Metric_ExponentialHistogram{
				ExponentialHistogram: &metricpb.ExponentialHistogram{
					AggregationTemporality: data.ExponentialHistogram.GetAggregationTemporality(),
					DataPoints: []*metricpb.ExponentialHistogramDataPoint{
						data.ExponentialHistogram.GetDataPoint(),
					},
				},
			},
		}
	case *FlattenDataPoint_Summary:
		return &metricpb.Metric{
			Name:        fdp.GetName(),
			Description: fdp.GetDescription(),
			Unit:        fdp.GetUnit(),
			Metadata:    fdp.GetMetadata(),
			Data: &metricpb.Metric_Summary{
				Summary: &metricpb.Summary{
					DataPoints: []*metricpb.SummaryDataPoint{
						data.Summary.GetDataPoint(),
					},
				},
			},
		}
	}
	return nil
}

func ConvertToFlattenLogRecords(resourceLogs []*logspb.ResourceLogs) []*FlattenLogRecord {
	if resourceLogs == nil {
		return nil
	}
	flattenLogRecords := make([]*FlattenLogRecord, 0, len(resourceLogs))
	for _, resourceLog := range resourceLogs {
		resource := resourceLog.GetResource()
		resourceAttrs := resource.GetAttributes()
		resourceDroppedAttrsCount := resource.GetDroppedAttributesCount()
		resourceLogSchemaURL := resourceLog.GetSchemaUrl()
		for _, scopeLog := range resourceLog.GetScopeLogs() {
			scope := scopeLog.GetScope()
			scopeName := scope.GetName()
			scopeVersion := scope.GetVersion()
			scopeAttrs := scope.GetAttributes()
			scopeDroppedAttrsCount := scope.GetDroppedAttributesCount()
			scopeLogSchemaURL := scopeLog.GetSchemaUrl()
			for _, logRecord := range scopeLog.GetLogRecords() {
				flattenLogRecords = append(flattenLogRecords, &FlattenLogRecord{
					ResourceAttributes:             resourceAttrs,
					ResourceLogSchemaUrl:           resourceLogSchemaURL,
					DroppedResourceAttributesCount: resourceDroppedAttrsCount,
					ScopeName:                      scopeName,
					ScopeVersion:                   scopeVersion,
					ScopeAttributes:                scopeAttrs,
					DroppedScopeAttributesCount:    scopeDroppedAttrsCount,
					ScopeLogSchemaUrl:              scopeLogSchemaURL,
					TraceId:                        logRecord.GetTraceId(),
					SpanId:                         logRecord.GetSpanId(),
					TimeUnixNano:                   logRecord.GetTimeUnixNano(),
					SeverityNumber:                 logRecord.GetSeverityNumber(),
					SeverityText:                   logRecord.GetSeverityText(),
					ObservedTimeUnixNano:           logRecord.GetObservedTimeUnixNano(),
					Attributes:                     logRecord.GetAttributes(),
					DroppedAttributesCount:         logRecord.GetDroppedAttributesCount(),
					Body:                           logRecord.GetBody(),
					Flags:                          logRecord.GetFlags(),
				})
			}
		}
	}
	return flattenLogRecords
}

func ConvertFromFlattenLogRecords(flr []*FlattenLogRecord) []*logspb.ResourceLogs {
	if flr == nil {
		return nil
	}
	resourceLogs := make([]*logspb.ResourceLogs, 0)
	for _, f := range flr {
		resource := &resourcepb.Resource{
			Attributes:             f.GetResourceAttributes(),
			DroppedAttributesCount: f.GetDroppedResourceAttributesCount(),
		}
		resourceLog := &logspb.ResourceLogs{
			Resource:  resource,
			SchemaUrl: f.GetResourceLogSchemaUrl(),
		}
		scope := &commonpb.InstrumentationScope{
			Name:                   f.GetScopeName(),
			Version:                f.GetScopeVersion(),
			Attributes:             f.GetScopeAttributes(),
			DroppedAttributesCount: f.GetDroppedScopeAttributesCount(),
		}
		scopeLog := &logspb.ScopeLogs{
			Scope:     scope,
			SchemaUrl: f.GetScopeLogSchemaUrl(),
		}
		logRecord := &logspb.LogRecord{
			TraceId:                f.GetTraceId(),
			SpanId:                 f.GetSpanId(),
			TimeUnixNano:           f.GetTimeUnixNano(),
			SeverityNumber:         f.GetSeverityNumber(),
			SeverityText:           f.GetSeverityText(),
			ObservedTimeUnixNano:   f.GetObservedTimeUnixNano(),
			Attributes:             f.GetAttributes(),
			DroppedAttributesCount: f.GetDroppedAttributesCount(),
			Body:                   f.GetBody(),
			Flags:                  f.GetFlags(),
		}
		scopeLog.LogRecords = append(scopeLog.LogRecords, logRecord)
		resourceLog.ScopeLogs = append(resourceLog.ScopeLogs, scopeLog)
		resourceLogs = otlp.AppendResourceLogs(resourceLogs, resourceLog)
	}
	return resourceLogs
}
