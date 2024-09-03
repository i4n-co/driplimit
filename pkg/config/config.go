package config

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/fatih/structtag"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/sethvargo/go-envconfig"
)

// Mode represents the mode of the service
type Mode string

const (
	// Authoritative mode is the default mode where the service is the source of truth
	Authoritative Mode = "authoritative"
	// AsyncAuthoritative mode is the authoritative mode with an asynchronous cache
	AsyncAuthoritative Mode = "async_authoritative"
	// Proxy mode is the mode where the service is a proxy to another service upstream
	Proxy Mode = "proxy"
)

func (m Mode) isValid() bool {
	switch m {
	case Authoritative, AsyncAuthoritative, Proxy:
		return true
	}
	return false
}

// Config represents the configuration of the service
type Config struct {
	Addr                 string        `env:"ADDR, default=127.0.0.1" description:"address to listen on"`
	CacheDuration        time.Duration `env:"CACHE_DURATION, default=30s" description:"cache entries time-to-live"`
	DatabaseName         string        `env:"DATABASE_NAME, default=driplimit.db" description:"database file name"`
	DataDir              string        `env:"DATA_DIR" description:"directory where the database file is stored"`
	GzipCompression      bool          `env:"GZIP_COMPRESSION, default=false" description:"enable gzip compression"`
	KeysCacheSize        int           `env:"KEYS_CACHE_SIZE, default=65536" description:"maximum number of keys in the cache"`
	LogFormat            string        `env:"LOG_FORMAT, default=text" description:"log format (text or json)"`
	LogSeverity          string        `env:"LOG_SEVERITY, default=info" description:"log severity level (debug, info, warn, error)"`
	Mode                 Mode          `env:"MODE, default=authoritative" description:"service mode (authoritative, async_authoritative, proxy)"`
	Port                 int           `env:"PORT, default=7131" description:"port to listen on"`
	ServiceKeysCacheSize int           `env:"SERVICE_KEYS_CACHE_SIZE, default=2048" description:"maximum number of service keys in the cache"`
	UpstreamTimeout      time.Duration `env:"UPSTREAM_TIMEOUT, default=5s" description:"timeout for upstream requests"`
	UpstreamURL          string        `env:"UPSTREAM_URL" description:"upstream URL for proxy mode or SDK client"`
	RootServiceKeyToken  string        `env:"ROOT_SERVICE_KEY_TOKEN" description:"create a root service key at startup with this token"`

	logger *slog.Logger
}

// FromEnv creates a new configuration by extracting environment variables
func FromEnv(ctx context.Context) (*Config, error) {
	cfg := new(Config)
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// FromEnvFile creates a new configuration by using variables from a file
func FromEnvFile(ctx context.Context, envfile io.Reader) (*Config, error) {
	envconfigmap := make(map[string]string)
	buf := bufio.NewReader(envfile)
	lineNumber := 0
	for {
		lineNumber++
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read line: %w", err)
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(string(line), "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line %d: %s", lineNumber, string(line))
		}
		envconfigmap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	cfg := new(Config)
	if err := envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:   cfg,
		Lookuper: envconfig.MapLookuper(envconfigmap),
	}); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// AddrPort returns the address and port in the format "addr:port".
// It supports IPv4 and IPv6 addresses.
func (c *Config) AddrPort() string {
	if c.UseIPv6Addr() {
		return fmt.Sprintf("[%s]:%d", c.Addr, c.Port)
	}

	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

// UseIPv6Addr returns true if the address is an IPv6 address
func (c *Config) UseIPv6Addr() bool {
	return net.ParseIP(c.Addr).To4() == nil && strings.Contains(c.Addr, ":")
}

// IsAuthoritative returns true if the service is in authoritative mode
func (c *Config) IsAuthoritative() bool {
	return Mode(c.Mode) == Authoritative
}

// IsAsyncAuthoritative returns true if the service is in async authoritative mode
func (c *Config) IsAsyncAuthoritative() bool {
	return Mode(c.Mode) == AsyncAuthoritative
}

// IsProxy returns true if the service is in proxy mode
func (c *Config) IsProxy() bool {
	return Mode(c.Mode) == Proxy
}

// validate validates the configuration
func (c *Config) validate() error {
	if !c.Mode.isValid() {
		return fmt.Errorf("invalid mode: %s", c.Mode)
	}
	if c.UpstreamTimeout <= 0 {
		return fmt.Errorf("invalid timeout: %d", c.UpstreamTimeout)
	}
	if c.Mode == Proxy && c.UpstreamURL == "" {
		return fmt.Errorf("upstream URL is required for proxy mode")
	}
	return nil
}

// InMemoryDatabase returns true if the service is using an in-memory database
func (c *Config) InMemoryDatabase() bool {
	return c.DataDir == ""
}

// DatabasePath returns the path to the database file
func (c *Config) DatabasePath() string {
	if c.InMemoryDatabase() {
		return ":memory:"
	}
	return filepath.Join(c.DataDir, c.DatabaseName)
}

// Logger initialize and returns a logger based on the configuration
func (c *Config) Logger() *slog.Logger {
	if c.logger != nil {
		return c.logger
	}

	var level slog.Level
	switch strings.ToLower(c.LogSeverity) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	var handler slog.Handler
	switch strings.ToLower(c.LogFormat) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:   level,
			NoColor: !isatty.IsTerminal(os.Stdout.Fd()),
		})
	}

	c.logger = slog.New(handler)

	return c.logger
}

type ConfigHelp struct {
	Default     string
	Description string
}

// Defaults returns all config keys and their default values
func (c *Config) Defaults() (map[string]ConfigHelp, error) {
	defaults := make(map[string]ConfigHelp)
	f := reflect.TypeOf(*c)
	for i := 0; i < f.NumField(); i++ {
		help := ConfigHelp{}
		tags, err := structtag.Parse(string(f.Field(i).Tag))
		if err != nil {
			return nil, err
		}
		envTag, err := tags.Get("env")
		if err != nil {
			continue
		}
		envKey := envTag.Name
		for _, opt := range envTag.Options {
			opt = strings.TrimSpace(opt)
			switch {
			case strings.HasPrefix(opt, "default="):
				help.Default = strings.TrimPrefix(opt, "default=")
			}
		}

		descriptionTag, _ := tags.Get("description")
		help.Description = descriptionTag.Value()
		defaults[envKey] = help
	}
	return defaults, nil
}

// PrintDefaults prints all config keys and their default values
func (c *Config) PrintDefaults() {
	defaults, err := c.Defaults()
	if err != nil {
		c.logger.Error("failed to get defaults", "error", err)
		os.Exit(1)
	}

	keys := make([]string, 0, len(defaults))
	for key := range defaults {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Println("# Driplimit default configuration")
	for _, key := range keys {
		if defaults[key].Description != "" {
			fmt.Printf("# %s: %s\n", key, defaults[key].Description)
		}
		fmt.Printf("%s=%s\n", key, defaults[key].Default)
	}
}
