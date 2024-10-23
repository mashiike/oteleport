package oteleport

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
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
	parser, err := kong.New(&opts, kong.Vars{"version": Version})
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to new kong: %w", err)
	}
	c, err := parser.Parse(args)
	if err != nil {
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
