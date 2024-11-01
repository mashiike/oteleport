package oteleport

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/samber/oops"
)

// Overall Configuration structure
type ServerConfig struct {
	AccessKeyHeader string             `json:"access_key_header"`
	AccessKeys      []*AccessKeyConfig `json:"access_keys"`
	Storage         StorageConfig      `json:"storage"`
	OTLP            OTLPConfig         `json:"otlp"`
	API             APIConfig          `json:"api"`
}

type AccessKeyConfig struct {
	KeyID     string `json:"key_id"`
	SecretKey string `json:"secret_key"`
}

type StorageConfig struct {
	CursorEncryptionKey []byte           `json:"cursor_encryption_key"`
	GZip                *bool            `json:"gzip,omitempty"`
	Flatten             *bool            `json:"flatten,omitempty"`
	Location            string           `json:"location"`
	locationURL         *url.URL         `json:"-"`
	AWS                 StorageAWSConfig `json:"aws,omitempty"`
}

type StorageAWSConfig struct {
	Region         string                       `json:"region"`
	Endpoint       string                       `json:"endpoint"`
	UseS3PathStyle bool                         `json:"use_s3_path_style"`
	Credentials    *StorageAWSCredentialsConfig `json:"credentials"`
	awsConfig      aws.Config                   `json:"-"`
}

type StorageAWSCredentialsConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
}

// OTLP gRPC and HTTP configuration
type OTLPConfig struct {
	Enable *bool          `json:"enable,omitempty"`
	GRPC   OTLPGRPCConfig `json:"grpc"`
	HTTP   OTLPHTTPConfig `json:"http"`
}

// gRPC configuration
type OTLPGRPCConfig struct {
	Enable   *bool        `json:"enable,omitempty"`
	Address  string       `json:"address"`
	Listener net.Listener `json:"-"`
}

// HTTP configuration
type OTLPHTTPConfig struct {
	Enable   *bool        `json:"enable,omitempty"`
	Address  string       `json:"address"`
	Listener net.Listener `json:"-"`
}

// API configuration
type APIConfig struct {
	Enable *bool         `json:"enable,omitempty"`
	HTTP   APIHTTPConfig `json:"http"`
	GRPC   APIGRPCConfig `json:"grpc"`
}

// API HTTP configuration
type APIHTTPConfig struct {
	Enable   *bool        `json:"enable,omitempty"`
	Address  string       `json:"address"`
	Listener net.Listener `json:"-"`
}

type APIGRPCConfig struct {
	Enable   *bool        `json:"enable,omitempty"`
	Address  string       `json:"address"`
	Listener net.Listener `json:"-"`
}

func Pointer[T any](v T) *T {
	return &v
}

func Coalasce[T any](v ...*T) *T {
	for _, p := range v {
		if p != nil {
			return p
		}
	}
	return nil
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		AccessKeyHeader: "Oteleport-Access-Key",
		OTLP: OTLPConfig{
			GRPC: OTLPGRPCConfig{
				Address: ":4317",
			},
			HTTP: OTLPHTTPConfig{
				Address: ":4318",
			},
		},
		API: APIConfig{
			HTTP: APIHTTPConfig{
				Address: ":8080",
			},
		},
	}
}

func (c *ServerConfig) EnableAuth() bool {
	return len(c.AccessKeys) > 0
}

// Validate function to check the configuration for validity
func (c *ServerConfig) Validate() error {
	if err := c.Storage.Validate(c); err != nil {
		return oops.Wrapf(err, "storage")
	}
	if err := c.OTLP.Validate(); err != nil {
		return oops.Wrapf(err, "otlp")
	}
	if err := c.API.Validate(); err != nil {
		return oops.Wrapf(err, "api")
	}
	keyIDs := make(map[string]int)
	for index, keyCfg := range c.AccessKeys {
		if keyCfg.KeyID == "" {
			keyCfg.KeyID = fmt.Sprintf("key%d", index)
		}
		if duplicateIndex, ok := keyIDs[keyCfg.KeyID]; ok {
			return oops.Errorf("duplicate access key id: index %d and %d", duplicateIndex, index)
		}
		keyIDs[keyCfg.KeyID] = index
		if keyCfg.SecretKey == "" {
			return oops.Errorf("access secret key index=%d is empty", index)
		}
	}
	return nil
}

type LoadOptions struct {
	ExtVars  map[string]string
	ExtCodes map[string]string
}

func (c *ServerConfig) Load(path string, opts *LoadOptions) error {
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
	if err := dec.Decode(c); err != nil {
		return oops.Wrapf(err, "failed to decode jsonnet file %s", path)
	}
	return c.Validate()
}

func (c *StorageConfig) Validate(parent *ServerConfig) error {
	if c.CursorEncryptionKey == nil {
		return oops.Errorf("cursor_encryption_key is required")
	}
	if c.GZip == nil {
		c.GZip = Coalasce(parent.Storage.GZip, Pointer(true))
	}
	if c.Flatten == nil {
		c.Flatten = Coalasce(parent.Storage.Flatten, Pointer(false))
	}
	if c.Location == "" {
		return oops.Errorf("location is required")
	}
	u, err := url.Parse(c.Location)
	if err != nil {
		return oops.Wrapf(err, "failed to parse location")
	}
	switch u.Scheme {
	case "s3":
		if u.Host == "" {
			return oops.Errorf("s3 bucket name is required")
		}
		if err := c.AWS.Validate(); err != nil {
			return oops.Wrapf(err, "aws")
		}
	default:
		return oops.Errorf("unsupported location scheme %s", u.Scheme)
	}
	c.locationURL = u
	return nil
}

func (c *StorageAWSConfig) Validate() error {
	loadConfigOptions := []func(*config.LoadOptions) error{}
	if c.Credentials != nil {
		loadConfigOptions = append(loadConfigOptions, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(c.Credentials.AccessKeyID, c.Credentials.SecretAccessKey, c.Credentials.SessionToken),
		))
	}
	if c.Region != "" {
		loadConfigOptions = append(loadConfigOptions, config.WithRegion(c.Region))
	}
	awsCfg, err := config.LoadDefaultConfig(context.Background(), loadConfigOptions...)
	if err != nil {
		return oops.Wrapf(err, "failed to load default aws config")
	}
	c.awsConfig = awsCfg
	return nil
}

func (c *OTLPConfig) Validate() error {
	if err := c.GRPC.Validate(c); err != nil {
		return oops.Wrapf(err, "grpc")
	}
	if err := c.HTTP.Validate(c); err != nil {
		return oops.Wrapf(err, "http")
	}
	return nil
}

func (c *OTLPGRPCConfig) Validate(parent *OTLPConfig) error {
	if c.Enable == nil {
		c.Enable = Coalasce(parent.Enable, Pointer(true))
	}
	if c.Listener != nil {
		c.Address = c.Listener.Addr().String()
	}
	if *c.Enable && c.Address == "" {
		return oops.Errorf("address is required")
	}
	return nil
}

func (c *OTLPHTTPConfig) Validate(parent *OTLPConfig) error {
	if c.Enable == nil {
		c.Enable = Coalasce(parent.Enable, Pointer(false))
	}
	if c.Listener != nil {
		c.Address = c.Listener.Addr().String()
	}
	if *c.Enable && c.Address == "" {
		return oops.Errorf("address is required")
	}
	return nil
}

func (c *APIConfig) Validate() error {
	if err := c.HTTP.Validate(c); err != nil {
		return oops.Wrapf(err, "http")
	}
	return nil
}

func (c *APIHTTPConfig) Validate(parent *APIConfig) error {
	if c.Enable == nil {
		c.Enable = Coalasce(parent.Enable, Pointer(true))
	}
	if c.Listener != nil {
		c.Address = c.Listener.Addr().String()
	}
	if *c.Enable && c.Address == "" {
		return oops.Errorf("address is required")
	}
	return nil
}

func (c *AccessKeyConfig) UnmarshalJSON(data []byte) error {
	type alias AccessKeyConfig
	aux := &struct {
		*alias
	}{
		alias: (*alias)(c),
	}
	structErr := json.Unmarshal(data, aux)
	if structErr == nil {
		return nil
	}
	var secretKey string
	if err := json.Unmarshal(data, &secretKey); err == nil {
		c.SecretKey = secretKey
		return nil
	}
	return structErr
}
