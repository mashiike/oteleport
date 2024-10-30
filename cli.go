package oteleport

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/mashiike/oteleport/pkg/client"
	"github.com/mashiike/slogutils"
)

type ServerCLIOptions struct {
	ConfigPath string            `name:"config" help:"config file path" default:"oteleport.jsonnet" env:"OTELPORT_CONFIG"`
	ExtStr     map[string]string `help:"external string values for Jsonnet" env:"OTELEPORT_EXTSTR"`
	ExtCode    map[string]string `help:"external code values for Jsonnet" env:"OTELEPORT_EXTCODE"`

	LogLevel string `help:"log level (debug, info, warn, error)" default:"info" enum:"debug,info,warn,error" env:"OTELPORT_LOG_LEVEL"`
	Color    bool   `help:"enable colored output" default:"false" env:"OTELPORT_COLOR"`

	Serve   struct{} `cmd:"" help:"start oteleport server" default:"1"`
	Version struct{} `cmd:"version" help:"show version"`
}

type ServerCLIParseFunc func([]string) (string, *ServerCLIOptions, func(), error)

func ParseServerCLI(args []string) (string, *ServerCLIOptions, func(), error) {

	var opts ServerCLIOptions
	parser, err := kong.New(&opts,
		kong.Name("oteleport"),
		kong.Description("oteleport is a OpenTelemetry signals receiver and REST API server."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{"version": Version},
	)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to new kong: %w", err)
	}
	c, err := parser.Parse(args)
	if err != nil {
		parser.FatalIfErrorf(err)
		return "", nil, nil, fmt.Errorf("failed to parse args: %w", err)
	}
	sub := strings.Fields(c.Command())[0]
	return sub, &opts, func() { c.PrintUsage(true) }, nil
}

func ServerCLI(ctx context.Context, parse ServerCLIParseFunc) (int, error) {
	sub, opts, usage, err := parse(os.Args[1:])
	if err != nil {
		return 1, err
	}
	if err := setupLogger(opts.LogLevel, opts.Color); err != nil {
		return 1, err
	}
	if err := dispatchServerCLI(ctx, sub, usage, opts); err != nil {
		return 1, err
	}
	return 0, nil
}

func dispatchServerCLI(ctx context.Context, sub string, usage func(), opts *ServerCLIOptions) error {
	switch sub {
	case "version", "":
		fmt.Println("oteleport-server", Version)
		return nil
	}

	switch sub {
	case "serve":
		cfg := DefaultServerConfig()
		if err := cfg.Load(opts.ConfigPath, &LoadOptions{
			ExtVars:  opts.ExtStr,
			ExtCodes: opts.ExtCode,
		}); err != nil {
			return err
		}
		s, err := NewServer(cfg)
		if err != nil {
			return err
		}
		return s.Run(ctx)
	default:
		usage()
	}
	return nil
}

type ClientCLIOptions struct {
	LogLevel string `help:"log level (debug, info, warn, error)" default:"info" enum:"debug,info,warn,error" env:"OTELPORT_LOG_LEVEL"`
	Color    bool   `help:"enable colored output" default:"false" env:"OTELPORT_COLOR"`

	ProfilePath string            `help:"oteleport client profile" default:"" env:"OTELPORT_PROFILE"`
	ExtStr      map[string]string `help:"external string values for Jsonnet" env:"OTELEPORT_EXTSTR"`
	ExtCode     map[string]string `help:"external code values for Jsonnet" env:"OTELEPORT_EXTCODE"`

	Endpoint        string `help:"oteleport server endpoint" default:"http://localhost:8080" env:"OTELPORT_ENDPOINT"`
	AccessKey       string `help:"oteleport server access key" env:"OTELPORT_ACCESS_KEY"`
	AccessKeyHeader string `help:"oteleport server access key header" default:"Oteleport-Access-Key" env:"OTELEPORT_ACCESS_KEY_HEADER"`
	ClientSignalOutputOptions

	Version struct{}                    `cmd:"version" help:"show version"`
	Traces  ClientTracesCommandOptions  `cmd:"traces" help:"traces subcommand"`
	Metrics ClientMetricsCommandOptions `cmd:"metrics" help:"metrics subcommand"`
	Logs    ClientLogsCommandOptions    `cmd:"logs" help:"logs subcommand"`
}

type ClientTracesCommandOptions struct {
	ClientTimeRangeOptions
}

type ClientMetricsCommandOptions struct {
	ClientTimeRangeOptions
}

type ClientLogsCommandOptions struct {
	ClientTimeRangeOptions
}

type ClientTimeRangeOptions struct {
	StartTime *time.Time `help:"return Otel Signals newer than this time. RFC3339 format" env:"OTELPORT_START_TIME" format:"2006-01-02T15:04:05Z"`
	EndTime   *time.Time `help:"return Otel Signals older than this time. RFC3339 format" env:"OTELPORT_END_TIME" format:"2006-01-02T15:04:05Z"`
	Since     string     `help:"return Otel Signals newer than a relative duration. like 52, 2m, or 3h (default: 5m)" env:"OTELPORT_SINCE" default:"5m"`
	Until     string     `help:"return Otel Signals older than a relative duration. like 52, 2m, or 3h" env:"OTELPORT_UNTIL"`
}

func (o *ClientTimeRangeOptions) TimeRangeUnixNano() (int64, int64) {
	var start, end time.Time
	if o.StartTime != nil {
		start = *o.StartTime
	}
	if o.EndTime != nil {
		end = *o.EndTime
	}
	if o.Since != "" {
		d, err := time.ParseDuration(o.Since)
		if err != nil {
			return 0, 0
		}
		start = time.Now().Add(-d)
	}
	if o.Until != "" {
		d, err := time.ParseDuration(o.Until)
		if err != nil {
			return 0, 0
		}
		end = time.Now().Add(-d)
	}
	if end.IsZero() {
		return start.UnixNano(), 0
	}
	return start.UnixNano(), end.UnixNano()
}

type ClientSignalOutputOptions struct {
	OtelExporterOTLPEndpoint string `help:"exporter endpoint: if not set,signal output to stdout" default:"" env:"OTEL_EXPORTER_OTLP_ENDPOINT" group:"OpenTelemetry Exporter Parameters"`
	OtelExporterOTLPProtocol string `help:"exporter protocol" default:"grpc" enum:"grpc,http" env:"OTEL_EXPORTER_OTLP_PROTOCOL" group:"OpenTelemetry Exporter Parameters"`
}

type ClientCLIParseFunc func([]string) (string, *ClientCLIOptions, func(), error)

func ParseClientCLI(args []string) (string, *ClientCLIOptions, func(), error) {
	var opts ClientCLIOptions
	parser, err := kong.New(
		&opts,
		kong.Name("oteleport-client"),
		kong.Description("oteleport-client is a CLI tool for oteleport server."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{"version": Version},
	)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to new kong: %w", err)
	}
	c, err := parser.Parse(args)
	if err != nil {
		parser.FatalIfErrorf(err)
		return "", nil, nil, fmt.Errorf("failed to parse args: %w", err)
	}
	sub := strings.Fields(c.Command())[0]
	return sub, &opts, func() { c.PrintUsage(true) }, nil
}

func ClientCLI(ctx context.Context, parse ClientCLIParseFunc) (int, error) {
	sub, opts, usage, err := parse(os.Args[1:])
	if err != nil {
		return 1, err
	}
	if err := setupLogger(opts.LogLevel, opts.Color); err != nil {
		return 1, err
	}

	if err := dispatchClientCLI(ctx, sub, usage, opts); err != nil {
		return 1, err
	}
	return 0, nil
}

func dispatchClientCLI(ctx context.Context, sub string, usage func(), opts *ClientCLIOptions) error {
	switch sub {
	case "version", "":
		fmt.Println("oteleport-client", Version)
		return nil
	}
	profile := &Profile{
		Profile: client.DefaultProfile(),
		Output:  opts.ClientSignalOutputOptions,
	}
	if opts.Endpoint != "" {
		profile.Endpoint = opts.Endpoint
	}
	if opts.AccessKey != "" {
		profile.AccessKey = opts.AccessKey
		if opts.AccessKeyHeader != "" {
			profile.AccessKeyHeader = opts.AccessKeyHeader
		}
	}
	if opts.ProfilePath != "" {
		if err := profile.Load(opts.ProfilePath, &LoadOptions{
			ExtVars:  opts.ExtStr,
			ExtCodes: opts.ExtCode,
		}); err != nil {
			return err
		}
	}
	app, err := NewClientApp(profile)
	if err != nil {
		return err
	}

	switch sub {
	case "traces":
		return app.FetchTracesData(ctx, &opts.Traces)
	case "metrics":
		return app.FetchMetricsData(ctx, &opts.Metrics)
	case "logs":
		return app.FetchLogsData(ctx, &opts.Logs)
	default:
		usage()
	}
	return nil
}

func setupLogger(l string, c bool) error {
	var level slog.Level
	if err := level.UnmarshalText([]byte(l)); err != nil {
		return fmt.Errorf("failed to unmarshal log level: %w", err)
	}
	color.NoColor = c
	logMiddleware := slogutils.NewMiddleware(
		slog.NewJSONHandler,
		slogutils.MiddlewareOptions{
			ModifierFuncs: map[slog.Level]slogutils.ModifierFunc{
				slog.LevelDebug: slogutils.Color(color.FgHiBlack),
				slog.LevelInfo:  nil,
				slog.LevelWarn:  slogutils.Color(color.FgYellow),
				slog.LevelError: slogutils.Color(color.FgRed),
			},
			Writer: os.Stderr,
			HandlerOptions: &slog.HandlerOptions{
				Level: level,
			},
		},
	)
	slog.SetDefault(slog.New(logMiddleware))
	return nil
}
