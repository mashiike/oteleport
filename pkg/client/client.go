package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"

	"github.com/mashiike/go-otlp-helper/otlp"
	oteleportpb "github.com/mashiike/oteleport/proto"
	"github.com/samber/oops"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Profile struct {
	Endpoint        string `json:"endpoint"`
	AccessKey       string `json:"access_key"`
	AccessKeyHeader string `json:"access_key_header"`
}

func (p *Profile) Validate() error {
	if p.Endpoint == "" {
		return oops.Errorf("endpoint is required")
	}
	if p.AccessKey != "" {
		if p.AccessKeyHeader == "" {
			return oops.Errorf("access_key_header is required")
		}
	}
	return nil
}

func DefaultProfile() *Profile {
	return &Profile{
		Endpoint:        "http://localhost:8080",
		AccessKey:       "",
		AccessKeyHeader: "Oteleport-Access-Key",
	}
}

type Client struct {
	p           *Profile
	endpointURL *url.URL
	httpClient  *http.Client
}

func New(p *Profile) (*Client, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	u, err := url.Parse(p.Endpoint)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to parse endpoint url")
	}
	return &Client{
		p:           p,
		endpointURL: u,
		httpClient:  http.DefaultClient,
	}, nil
}

func (c *Client) FetchTracesData(ctx context.Context, req *oteleportpb.FetchTracesDataRequest) (*oteleportpb.FetchTracesDataResponse, error) {
	var resp oteleportpb.FetchTracesDataResponse
	if err := c.do(ctx, "/api/traces/fetch", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) FetchMetricsData(ctx context.Context, req *oteleportpb.FetchMetricsDataRequest) (*oteleportpb.FetchMetricsDataResponse, error) {
	var resp oteleportpb.FetchMetricsDataResponse
	if err := c.do(ctx, "/api/metrics/fetch", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) FetchLogsData(ctx context.Context, req *oteleportpb.FetchLogsDataRequest) (*oteleportpb.FetchLogsDataResponse, error) {
	var resp oteleportpb.FetchLogsDataResponse
	if err := c.do(ctx, "/api/logs/fetch", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type FetchTracesDataPagenator struct {
	c           *Client
	req         *oteleportpb.FetchTracesDataRequest
	firstPage   bool
	hasMorePage bool
}

func NewFetchTracesDataPagenator(c *Client, req *oteleportpb.FetchTracesDataRequest) *FetchTracesDataPagenator {
	if req == nil {
		req = &oteleportpb.FetchTracesDataRequest{}
	}
	return &FetchTracesDataPagenator{
		c:         c,
		req:       req,
		firstPage: true,
	}
}

func (p *FetchTracesDataPagenator) HasMorePages() bool {
	return p.firstPage || p.hasMorePage
}

func (p *FetchTracesDataPagenator) NextPage(ctx context.Context) (*oteleportpb.FetchTracesDataResponse, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	resp, err := p.c.FetchTracesData(ctx, p.req)
	if err != nil {
		return nil, err
	}
	p.firstPage = false
	nextCursor := resp.GetNextCursor()
	if nextCursor == "" {
		slog.DebugContext(ctx, "no more pages available")
		p.hasMorePage = false
		p.req.Cursor = ""
	} else {
		slog.DebugContext(ctx, "more pages available", "next_cursor", nextCursor)
		p.hasMorePage = true
		p.req.Cursor = nextCursor
	}
	return resp, nil
}

type FetchMetricsDataPagenator struct {
	c           *Client
	req         *oteleportpb.FetchMetricsDataRequest
	firstPage   bool
	hasMorePage bool
}

func NewFetchMetricsDataPagenator(c *Client, req *oteleportpb.FetchMetricsDataRequest) *FetchMetricsDataPagenator {
	if req == nil {
		req = &oteleportpb.FetchMetricsDataRequest{}
	}
	return &FetchMetricsDataPagenator{
		c:         c,
		req:       req,
		firstPage: true,
	}
}

func (p *FetchMetricsDataPagenator) HasMorePages() bool {
	return p.firstPage || p.hasMorePage
}

func (p *FetchMetricsDataPagenator) NextPage(ctx context.Context) (*oteleportpb.FetchMetricsDataResponse, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	resp, err := p.c.FetchMetricsData(ctx, p.req)
	if err != nil {
		return nil, err
	}
	p.firstPage = false
	nextCursor := resp.GetNextCursor()
	if nextCursor == "" {
		slog.DebugContext(ctx, "no more pages available")
		p.hasMorePage = false
		p.req.Cursor = ""
	} else {
		slog.DebugContext(ctx, "more pages available", "next_cursor", nextCursor)
		p.hasMorePage = true
		p.req.Cursor = nextCursor
	}
	return resp, nil
}

type FetchLogsDataPagenator struct {
	c           *Client
	req         *oteleportpb.FetchLogsDataRequest
	firstPage   bool
	hasMorePage bool
}

func NewFetchLogsDataPagenator(c *Client, req *oteleportpb.FetchLogsDataRequest) *FetchLogsDataPagenator {
	if req == nil {
		req = &oteleportpb.FetchLogsDataRequest{}
	}
	return &FetchLogsDataPagenator{
		c:         c,
		req:       req,
		firstPage: true,
	}
}

func (p *FetchLogsDataPagenator) HasMorePages() bool {
	return p.firstPage || p.hasMorePage
}

func (p *FetchLogsDataPagenator) NextPage(ctx context.Context) (*oteleportpb.FetchLogsDataResponse, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	resp, err := p.c.FetchLogsData(ctx, p.req)
	if err != nil {
		return nil, err
	}
	p.firstPage = false
	nextCursor := resp.GetNextCursor()
	if nextCursor == "" {
		slog.DebugContext(ctx, "no more pages available")
		p.hasMorePage = false
		p.req.Cursor = ""
	} else {
		slog.DebugContext(ctx, "more pages available", "next_cursor", nextCursor)
		p.hasMorePage = true
		p.req.Cursor = nextCursor
	}
	return resp, nil
}

func (c *Client) do(ctx context.Context, path string, protoBody proto.Message, respBody proto.Message) error {
	body, err := proto.Marshal(protoBody)
	if err != nil {
		return oops.Wrapf(err, "failed to marshal request")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpointURL.JoinPath(path).String(), bytes.NewReader(body))
	if err != nil {
		return oops.Wrapf(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	if c.p.AccessKey != "" {
		req.Header.Set(c.p.AccessKeyHeader, c.p.AccessKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return oops.Wrapf(err, "failed to do request")
	}
	defer resp.Body.Close()
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return oops.Wrapf(err, "failed to read response body")
	}
	if resp.StatusCode != http.StatusOK {
		var protoStatus spb.Status
		if unmarshalErr := unmarshalBody(resp.Header.Get("Content-Type"), respBodyBytes, &protoStatus); unmarshalErr != nil {
			slog.WarnContext(ctx, "failed to unmarshal response body", "message", unmarshalErr.Error())
			return oops.Errorf("failed fetch traces data: status code %d", resp.StatusCode)
		}
		st := status.FromProto(&protoStatus)
		return oops.Errorf("faild fetch traces data: code=%s, message=%s", st.Code(), st.Message())
	}
	if err := unmarshalBody(resp.Header.Get("Content-Type"), respBodyBytes, respBody); err != nil {
		return oops.Wrapf(err, "failed to unmarshal response body")
	}
	return nil
}

func unmarshalBody(contentType string, data []byte, pb proto.Message) error {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return oops.Wrapf(err, "failed to parse media type")
	}
	slog.Debug("unmarshalBody", "mediaType", mediaType, "contentType", contentType)
	switch mediaType {
	case "application/x-protobuf", "application/protobuf":
		return proto.Unmarshal(data, pb)
	case "application/json":
		return otlp.UnmarshalJSON(data, pb)
	default:
		return oops.Errorf("unsupported content type %s", contentType)
	}
}
