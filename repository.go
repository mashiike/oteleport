package oteleport

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mashiike/go-otlp-helper/otlp"
	oteleportpb "github.com/mashiike/oteleport/proto"
	"github.com/samber/oops"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	metricspb "go.opentelemetry.io/proto/otlp/metrics/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SignalRepository interface {
	PushTracesData(ctx context.Context, data *oteleportpb.TracesData) error
	PushMetricsData(ctx context.Context, data *oteleportpb.MetricsData) error
	PushLogsData(ctx context.Context, data *oteleportpb.LogsData) error
	FetchTracesData(ctx context.Context, input *oteleportpb.FetchTracesDataRequest) (*oteleportpb.FetchTracesDataResponse, error)
	FetchMetricsData(ctx context.Context, input *oteleportpb.FetchMetricsDataRequest) (*oteleportpb.FetchMetricsDataResponse, error)
	FetchLogsData(ctx context.Context, input *oteleportpb.FetchLogsDataRequest) (*oteleportpb.FetchLogsDataResponse, error)
}

type S3SignalRepository struct {
	bucketName          string
	objectPathPrefix    string
	client              *s3.Client
	gzip                bool
	cursorEncryptionKey []byte
	uploader            *manager.Uploader
	downloader          *manager.Downloader
}

func NewSignalRepository(cfg *StorageConfig) (SignalRepository, error) {
	switch cfg.locationURL.Scheme {
	case "s3":
		return newS3SignalRepository(cfg), nil
	default:
		return nil, oops.Errorf("unsupported location scheme %s", cfg.locationURL.Scheme)
	}
}

func newS3SignalRepository(cfg *StorageConfig) *S3SignalRepository {
	s3Opts := []func(*s3.Options){}
	if cfg.AWS.Endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.AWS.Endpoint)
		})
	}
	if cfg.AWS.UseS3PathStyle {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}
	client := s3.NewFromConfig(cfg.AWS.awsConfig, s3Opts...)
	return &S3SignalRepository{
		cursorEncryptionKey: adjustKey(cfg.CursorEncryptionKey, 32),
		gzip:                cfg.GZip != nil && *cfg.GZip,
		bucketName:          cfg.locationURL.Host,
		objectPathPrefix:    strings.TrimPrefix(cfg.locationURL.Path, "/"),
		client:              client,
		uploader:            manager.NewUploader(client),
		downloader:          manager.NewDownloader(client),
	}
}

var randReader = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[randReader.Intn(len(chars))]
	}
	return string(b)
}

const (
	partitionForamt = "2006/01/02/15"
)

func (r *S3SignalRepository) PushTracesData(ctx context.Context, data *oteleportpb.TracesData) error {
	partitionBy := otlp.PartitionResourceSpans(data.GetResourceSpans(), func(rs *tracepb.ResourceSpans) string {
		if str := otlp.PartitionBySpanStartTime(partitionForamt, time.Local)(rs); str != "" {
			return str
		}
		if str := otlp.PartitionBySpanEndTime(partitionForamt, time.Local)(rs); str != "" {
			return str
		}
		return time.Now().Format(partitionForamt)
	})
	for partition, spans := range partitionBy {
		bs, err := otlp.MarshalJSON(&oteleportpb.TracesData{
			ResourceSpans: spans,
			SignalType:    data.GetSignalType(),
		})
		if err != nil {
			return oops.Wrapf(err, "failed to marshal json")
		}
		spansCount := otlp.TotalSpans(spans)
		slog.DebugContext(ctx, "push traces data", "partition", partition, "spans", spansCount)
		objectKeySuffix := fmt.Sprintf("traces/%s/spans-%s-%s.json", partition, time.Now().Format("20060102150405"), RandomString(8))
		if err := r.putObject(ctx, objectKeySuffix, strings.NewReader(string(bs))); err != nil {
			return oops.Wrapf(err, "failed to put object")
		}
	}
	return nil
}

var zeroTimeStr = time.Unix(0, 0).In(time.Local).Format(partitionForamt)

func (r *S3SignalRepository) PushMetricsData(ctx context.Context, data *oteleportpb.MetricsData) error {
	partitionBy := otlp.PartitionResourceMetrics(data.GetResourceMetrics(), func(rm *metricspb.ResourceMetrics) string {
		if str := otlp.PartitionByMetricStartTime(partitionForamt, time.Local)(rm); str != "" && str != zeroTimeStr {
			return str
		}
		if str := otlp.PartitionByMetricTime(partitionForamt, time.Local)(rm); str != "" && str != zeroTimeStr {
			return str
		}
		return time.Now().Format(partitionForamt)
	})
	for partition, metrics := range partitionBy {
		bs, err := otlp.MarshalJSON(&oteleportpb.MetricsData{
			ResourceMetrics: metrics,
			SignalType:      data.GetSignalType(),
		})
		if err != nil {
			return oops.Wrapf(err, "failed to marshal json")
		}
		objectKeySuffix := fmt.Sprintf("metrics/%s/data-points-%s-%s.json", partition, time.Now().Format("20060102150405"), RandomString(8))
		if err := r.putObject(ctx, objectKeySuffix, strings.NewReader(string(bs))); err != nil {
			return oops.Wrapf(err, "failed to put object")
		}
	}
	return nil
}

func (r *S3SignalRepository) PushLogsData(ctx context.Context, data *oteleportpb.LogsData) error {
	partitionBy := otlp.PartitionResourceLogs(data.GetResourceLogs(), func(rl *logspb.ResourceLogs) string {
		if str := otlp.PartitionByLogTime(partitionForamt, time.Local)(rl); str != "" {
			return str
		}
		if str := otlp.PartitionByLogObservedTime(partitionForamt, time.Local)(rl); str != "" {
			return str
		}
		return time.Now().Format(partitionForamt)
	})
	for partition, logs := range partitionBy {
		bs, err := otlp.MarshalJSON(&oteleportpb.LogsData{
			ResourceLogs: logs,
			SignalType:   data.GetSignalType(),
		})
		if err != nil {
			return oops.Wrapf(err, "failed to marshal json")
		}
		objectKeySuffix := fmt.Sprintf("logs/%s/records-%s-%s.json", partition, time.Now().Format("20060102150405"), RandomString(8))
		if err := r.putObject(ctx, objectKeySuffix, strings.NewReader(string(bs))); err != nil {
			return oops.Wrapf(err, "failed to put object")
		}
	}
	return nil
}

func (r *S3SignalRepository) putObject(ctx context.Context, objectKeySuffix string, body io.Reader) error {
	objKey := filepath.Join(r.objectPathPrefix, objectKeySuffix)
	var contentEncoding *string
	if r.gzip {
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		if _, err := io.Copy(gzipWriter, body); err != nil {
			return oops.Wrapf(err, "failed to write gzip")
		}
		if err := gzipWriter.Close(); err != nil {
			return oops.Wrapf(err, "failed to close gzip")
		}
		body = &buf
		objKey += ".gz"
		contentEncoding = aws.String("gzip")
	}
	output, err := r.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:          aws.String(r.bucketName),
		Key:             aws.String(objKey),
		Body:            body,
		ContentType:     aws.String("application/json"),
		ContentEncoding: contentEncoding,
	})
	if err != nil {
		return oops.Wrapf(err, "failed to put object")
	}
	slog.InfoContext(ctx, "put object", "s3_url", output.Location, "etag", output.ETag, "version_id", output.VersionID)
	return nil
}

func (r *S3SignalRepository) walkObjects(
	ctx context.Context,
	startTime time.Time, endTime time.Time,
	startAfter *string,
	getObjectKeyPrefixFunc func(time.Time) string,
	f func(context.Context, time.Time, types.Object) (bool, error),
) (bool, error) {
	currentTime := startTime.Truncate(time.Hour)
	slog.DebugContext(ctx, "start walk objects", "start_time", startTime, "end_time", endTime, "current_time", currentTime, "is_equal", currentTime.Equal(endTime), "is_before", currentTime.Before(endTime))
	for currentTime.Before(endTime) || currentTime.Equal(endTime) {
		objectKeyPrefix := getObjectKeyPrefixFunc(currentTime)
		slog.DebugContext(ctx, "list objects", "prefix", objectKeyPrefix, "start_after", startAfter)
		paginator := s3.NewListObjectsV2Paginator(r.client, &s3.ListObjectsV2Input{
			Bucket:     aws.String(r.bucketName),
			Prefix:     aws.String(objectKeyPrefix),
			StartAfter: startAfter,
		})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return false, oops.Wrapf(err, "failed to list objects")
			}
			for _, obj := range page.Contents {
				ok, err := f(ctx, currentTime, obj)
				if err != nil {
					return false, err
				}
				if !ok {
					return false, nil
				}
			}
		}
		currentTime = currentTime.Add(time.Hour)
		slog.DebugContext(ctx, "next walk", "start_time", startTime, "end_time", endTime, "current_time", currentTime, "is_equal", currentTime.Equal(endTime), "is_before", currentTime.Before(endTime))
	}
	slog.DebugContext(ctx, "end walk objects", "start_time", startTime, "end_time", endTime, "current_time", currentTime)
	return true, nil
}

func (r *S3SignalRepository) getObjectBody(ctx context.Context, obj types.Object) ([]byte, error) {
	var buf = make([]byte, 1024*1024*5) //5MB
	w := manager.NewWriteAtBuffer(buf)
	n, err := r.downloader.Download(ctx, w, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(*obj.Key),
	})
	if err != nil {
		return nil, oops.Wrapf(err, "failed to get object")
	}
	body := w.Bytes()[:n]
	gzipReader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		if err == gzip.ErrHeader {
			return body, nil
		}
		return nil, oops.Wrapf(err, "failed to create gzip reader")
	}
	body, err = io.ReadAll(gzipReader)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to read gzip")
	}
	return body, nil
}

type s3Cursor struct {
	CurrentTime      time.Time `json:"ct"`
	CurrentObjectKey *string   `json:"ck"`
	Offset           int       `json:"o"`
}

func (c *s3Cursor) encrypt(key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", oops.Wrapf(err, "failed to create cipher")
	}
	plainText, err := json.Marshal(c)
	if err != nil {
		return "", oops.Wrapf(err, "failed to marshal json")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return "", oops.Wrapf(err, "failed to read random")
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (c *s3Cursor) decrypt(encryptedCursor string, key []byte) error {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedCursor)
	if err != nil {
		return oops.Wrapf(err, "failed to decode base64")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return oops.Wrapf(err, "failed to create cipher")
	}

	if len(ciphertext) < aes.BlockSize {
		return fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	if err := json.Unmarshal(ciphertext, c); err != nil {
		return oops.Wrapf(err, "failed to unmarshal json")
	}
	return nil
}

func adjustKey(key []byte, size int) []byte {
	if len(key) > size {
		return key[:size]
	}
	paddedKey := make([]byte, size)
	copy(paddedKey, key)
	return paddedKey
}

func (r *S3SignalRepository) FetchTracesData(ctx context.Context, input *oteleportpb.FetchTracesDataRequest) (*oteleportpb.FetchTracesDataResponse, error) {
	startTime, endTime, limit, err := validateRequest(input.GetStartTimeUnixNano(), input.GetEndTimeUnixNano(), input.GetLimit())
	if err != nil {
		return nil, err
	}
	cursor := input.GetCursor()
	slog.InfoContext(ctx, "fetch traces data", "start_time", startTime, "end_time", endTime, "cursor", cursor, "limit", limit)
	resp := &oteleportpb.FetchTracesDataResponse{}
	num := 0
	cursorObj := &s3Cursor{}
	if cursor != "" {
		if err := cursorObj.decrypt(cursor, r.cursorEncryptionKey); err != nil {
			errID := RandomString(8)
			slog.ErrorContext(ctx, "failed to unmarshal cursor", "error_id", errID, "error", err.Error())
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid cursor: err_id=%s", errID))
		}
		if !cursorObj.CurrentTime.IsZero() {
			startTime = cursorObj.CurrentTime
		}
		slog.DebugContext(ctx, "cursor", "current_time", cursorObj.CurrentTime, "current_object_key", cursorObj.CurrentObjectKey, "offset", cursorObj.Offset, "start_time", startTime)
	}
	noMore, err := r.walkObjects(
		ctx,
		startTime,
		endTime,
		cursorObj.CurrentObjectKey,
		func(t time.Time) string {
			key := fmt.Sprintf("traces/%s/", t.Format(partitionForamt))
			if r.objectPathPrefix != "" {
				key = filepath.Join(r.objectPathPrefix, key)
			}
			return key
		},
		func(ctx context.Context, t time.Time, obj types.Object) (bool, error) {
			slog.DebugContext(ctx, "fetch object", "key", *obj.Key)
			body, err := r.getObjectBody(ctx, obj)
			if err != nil {
				return false, oops.Wrapf(err, "failed to get object %q", *obj.Key)
			}
			var data oteleportpb.TracesData
			if err := otlp.UnmarshalJSON(body, &data); err != nil {
				return false, oops.Wrapf(err, "failed to unmarshal json")
			}
			resourceSpans := otlp.FilterResourceSpans(
				data.GetResourceSpans(),
				otlp.SpanInTimeRangeFilter(startTime, endTime),
			)
			dataLen := len(resourceSpans)
			slog.DebugContext(ctx, "restore spans", "spans", dataLen, "key", *obj.Key, "current", num, "limit", limit)
			if cursorObj.Offset != 0 && cursorObj.Offset < len(resourceSpans) {
				resourceSpans = resourceSpans[cursorObj.Offset:]
				slog.DebugContext(ctx, "skip spans", "offset", cursorObj.Offset, "key", *obj.Key, "spans", len(resourceSpans))
			}
			if num+dataLen > int(limit) {
				dataLen = int(limit) - num
				resourceSpans = resourceSpans[:dataLen]
				cursorObj.Offset += dataLen
				slog.DebugContext(ctx, "limit over in one object", "current", num, "limit", limit, "data_len", dataLen, "offset", cursorObj.Offset, "key", *obj.Key, "spans", len(resourceSpans))
			} else {
				cursorObj.CurrentTime = t
				cursorObj.CurrentObjectKey = obj.Key
				cursorObj.Offset = 0
				slog.DebugContext(ctx, "all data in one object", "current", num, "limit", limit, "data_len", dataLen, "offset", cursorObj.Offset, "key", *obj.Key, "spans", len(resourceSpans))
			}
			resp.ResourceSpans = otlp.AppendResourceSpans(resp.GetResourceSpans(), resourceSpans...)
			num += dataLen
			if int64(num) >= limit {
				slog.DebugContext(ctx, "limit exceeded", "current", num, "limit", limit, "cursor", cursorObj)
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		errID := RandomString(8)
		slog.ErrorContext(ctx, "failed to fetch traces data", "error_id", errID, "error", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch traces data: err_id=%s", errID))
	}
	slog.InfoContext(ctx, "fetched traces data", "num", num, "limit", limit, "no_more", noMore)
	if noMore {
		return resp, nil
	}
	resp.HasMore = true
	resp.NextCursor, err = cursorObj.encrypt(r.cursorEncryptionKey)
	if err != nil {
		errID := RandomString(8)
		slog.ErrorContext(ctx, "failed to marshal cursor", "error_id", errID, "error", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to marshal cursor: err_id=%s", errID))
	}
	return resp, nil
}

func (r *S3SignalRepository) FetchMetricsData(ctx context.Context, input *oteleportpb.FetchMetricsDataRequest) (*oteleportpb.FetchMetricsDataResponse, error) {
	startTime, endTime, limit, err := validateRequest(input.GetStartTimeUnixNano(), input.GetEndTimeUnixNano(), input.GetLimit())
	if err != nil {
		return nil, err
	}
	cursor := input.GetCursor()
	slog.InfoContext(ctx, "fetch metrics data", "start_time", startTime, "end_time", endTime, "cursor", cursor, "limit", limit)
	cursorObj := &s3Cursor{}
	if cursor != "" {
		if err := cursorObj.decrypt(cursor, r.cursorEncryptionKey); err != nil {
			errID := RandomString(8)
			slog.ErrorContext(ctx, "failed to unmarshal cursor", "error_id", errID, "error", err.Error())
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid cursor: %s", errID))
		}
		if !cursorObj.CurrentTime.IsZero() {
			startTime = cursorObj.CurrentTime
		}
		slog.DebugContext(ctx, "cursor", "current_time", cursorObj.CurrentTime, "current_object_key", cursorObj.CurrentObjectKey, "offset", cursorObj.Offset, "start_time", startTime)
	}
	resp := &oteleportpb.FetchMetricsDataResponse{}
	num := 0
	noMore, err := r.walkObjects(
		ctx,
		startTime,
		endTime,
		cursorObj.CurrentObjectKey,
		func(t time.Time) string {
			key := fmt.Sprintf("metrics/%s/", t.Format(partitionForamt))
			if r.objectPathPrefix != "" {
				key = filepath.Join(r.objectPathPrefix, key)
			}
			return key
		},
		func(ctx context.Context, t time.Time, obj types.Object) (bool, error) {
			slog.DebugContext(ctx, "fetch object", "key", *obj.Key)
			body, err := r.getObjectBody(ctx, obj)
			if err != nil {
				return false, oops.Wrapf(err, "failed to get object %q", *obj.Key)
			}
			var data oteleportpb.MetricsData
			if err := otlp.UnmarshalJSON(body, &data); err != nil {
				return false, oops.Wrapf(err, "failed to unmarshal json")
			}
			resourceMetrics := otlp.FilterResourceMetrics(
				data.GetResourceMetrics(),
				otlp.MetricDataPointInTimeRangeFilter(startTime, endTime),
			)
			dataLen := len(resourceMetrics)
			slog.DebugContext(ctx, "restore metrics", "metrics", dataLen, "key", *obj.Key)
			if cursorObj.Offset != 0 && cursorObj.Offset < len(resourceMetrics) {
				resourceMetrics = resourceMetrics[cursorObj.Offset:]
				slog.DebugContext(ctx, "skip metrics", "offset", cursorObj.Offset, "key", *obj.Key, "metrics", len(resourceMetrics))
			}
			if cursorObj.Offset+dataLen > int(limit) {
				dataLen = int(limit) - cursorObj.Offset
				resourceMetrics = resourceMetrics[:dataLen]
				cursorObj.Offset += dataLen
				slog.DebugContext(ctx, "limit over in one object", "current", cursorObj.Offset, "limit", limit, "data_len", dataLen, "offset", cursorObj.Offset, "key", *obj.Key, "metrics", len(resourceMetrics))
			} else {
				cursorObj.CurrentTime = t
				cursorObj.CurrentObjectKey = obj.Key
				cursorObj.Offset = 0
			}
			resp.ResourceMetrics = otlp.AppendResourceMetrics(resp.GetResourceMetrics(), resourceMetrics...)
			num += dataLen
			if int64(num) >= limit {
				slog.DebugContext(ctx, "limit exceeded", "current", num, "limit", limit, "cursor", cursorObj)
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		errID := RandomString(8)
		slog.ErrorContext(ctx, "failed to fetch metrics data", "error_id", errID, "error", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch metrics data: %s", errID))
	}
	slog.InfoContext(ctx, "fetch metrics data", "num", num, "limit", limit, "cursor", cursorObj)
	if noMore {
		return resp, nil
	}
	resp.HasMore = true
	resp.NextCursor, err = cursorObj.encrypt(r.cursorEncryptionKey)
	if err != nil {
		errID := RandomString(8)
		slog.ErrorContext(ctx, "failed to marshal cursor", "error_id", errID, "error", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to marshal cursor: %s", errID))
	}
	return resp, nil
}

func (r *S3SignalRepository) FetchLogsData(ctx context.Context, input *oteleportpb.FetchLogsDataRequest) (*oteleportpb.FetchLogsDataResponse, error) {
	startTime, endTime, limit, err := validateRequest(input.GetStartTimeUnixNano(), input.GetEndTimeUnixNano(), input.GetLimit())
	if err != nil {
		return nil, err
	}
	cursor := input.GetCursor()
	slog.InfoContext(ctx, "fetch logs data", "start_time", startTime, "end_time", endTime, "cursor", cursor, "limit", limit)
	cursorObj := &s3Cursor{}
	if cursor != "" {
		if err := cursorObj.decrypt(cursor, r.cursorEncryptionKey); err != nil {
			errID := RandomString(8)
			slog.ErrorContext(ctx, "failed to unmarshal cursor", "error_id", errID, "error", err.Error())
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid cursor: %s", errID))
		}
		if !cursorObj.CurrentTime.IsZero() {
			startTime = cursorObj.CurrentTime
		}
		slog.DebugContext(ctx, "cursor", "current_time", cursorObj.CurrentTime, "current_object_key", cursorObj.CurrentObjectKey, "offset", cursorObj.Offset, "start_time", startTime)
	}
	resp := &oteleportpb.FetchLogsDataResponse{}
	num := 0
	noMore, err := r.walkObjects(
		ctx,
		startTime,
		endTime,
		cursorObj.CurrentObjectKey,
		func(t time.Time) string {
			key := fmt.Sprintf("logs/%s/", t.Format(partitionForamt))
			if r.objectPathPrefix != "" {
				key = filepath.Join(r.objectPathPrefix, key)
			}
			return key
		},
		func(ctx context.Context, t time.Time, obj types.Object) (bool, error) {
			slog.DebugContext(ctx, "fetch object", "key", *obj.Key)
			body, err := r.getObjectBody(ctx, obj)
			if err != nil {
				return false, oops.Wrapf(err, "failed to get object %q", *obj.Key)
			}
			var data oteleportpb.LogsData
			if err := otlp.UnmarshalJSON(body, &data); err != nil {
				return false, oops.Wrapf(err, "failed to unmarshal json")
			}
			resourceLogs := otlp.FilterResourceLogs(
				data.GetResourceLogs(),
				otlp.LogRecordInTimeRangeFilter(startTime, endTime),
			)
			dataLen := len(resourceLogs)
			slog.DebugContext(ctx, "restore logs", "logs", dataLen, "key", *obj.Key)
			if cursorObj.Offset != 0 && cursorObj.Offset < len(resourceLogs) {
				resourceLogs = resourceLogs[cursorObj.Offset:]
				slog.DebugContext(ctx, "skip logs", "offset", cursorObj.Offset, "key", *obj.Key, "logs", len(resourceLogs))
			}
			if cursorObj.Offset+dataLen > int(limit) {
				dataLen = int(limit) - cursorObj.Offset
				resourceLogs = resourceLogs[:dataLen]
				cursorObj.Offset += dataLen
			} else {
				cursorObj.CurrentTime = t
				cursorObj.CurrentObjectKey = obj.Key
				cursorObj.Offset = 0
			}
			resp.ResourceLogs = otlp.AppendResourceLogs(resp.GetResourceLogs(), resourceLogs...)
			num += dataLen
			if int64(num) >= limit {
				slog.DebugContext(ctx, "limit exceeded", "current", num, "limit", limit, "cursor", cursorObj)
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		errID := RandomString(8)
		slog.ErrorContext(ctx, "failed to fetch logs data", "error_id", errID, "error", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to fetch logs data: %s", errID))
	}
	slog.InfoContext(ctx, "fetched logs data", "num", num, "limit", limit, "no_more", noMore)
	if noMore {
		return resp, nil
	}
	resp.HasMore = true
	resp.NextCursor, err = cursorObj.encrypt(r.cursorEncryptionKey)
	if err != nil {
		errID := RandomString(8)
		slog.ErrorContext(ctx, "failed to marshal cursor", "error_id", errID, "error", err.Error())
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to marshal cursor: %s", errID))
	}
	return resp, nil
}

func validateRequest(startTimeUnixNano uint64, endTimeUnixNano uint64, limit int64) (time.Time, time.Time, int64, error) {
	if startTimeUnixNano == 0 {
		return time.Time{}, time.Time{}, 0, status.Error(codes.InvalidArgument, "start time is required")
	}
	startTime := time.Unix(0, int64(startTimeUnixNano)).In(time.Local)
	endTime := time.Unix(0, int64(endTimeUnixNano))
	if endTimeUnixNano == 0 {
		endTime = time.Now()
	}
	endTime = endTime.In(time.Local)
	if startTime.After(endTime) {
		return time.Time{}, time.Time{}, 0, status.Error(codes.InvalidArgument, "start time is after end time")
	}
	if limit < 0 {
		return time.Time{}, time.Time{}, 0, status.Error(codes.InvalidArgument, "limit is negative")
	}
	if limit > 10000 {
		return time.Time{}, time.Time{}, 0, status.Error(codes.InvalidArgument, "limit is too large")
	}
	if limit == 0 {
		limit = 10000
	}
	return startTime, endTime, limit, nil
}
