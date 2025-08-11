package configclient

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ToxicToast/Azkaban-Go/libs/shared/helper"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	App struct {
		Name    string
		Env     string
		Version string
	}
	Auth struct {
		Provider       string
		JwksUrl        string `koanf:"AUTHENTIK_JWKS_URL"`
		Audience       string
		Issuer         string `koanf:"AUTHENTIK_ISSUER"`
		RequiredGroups []string
		RoutePolicies  []struct {
			Prefix string
			Groups []string
		}
	}
	Cache struct {
		redis struct {
			Addr         string `koanf:"REDIS_ADDR"`
			DB           int    `koanf:"REDIS_DB"`
			TLS          bool
			DialTimeout  time.Duration
			ReadTimeout  time.Duration
			WriteTimeout time.Duration
			DefaultTTL   time.Duration
			KeyPrefix    string `koanf:"REDIS_PREFIX"`
			Password     string `koanf:"REDIS_PASSWORD"`
		}
	}
	Downstreams struct {
		DiscoveryMode string
		Services      map[string]helper.ServiceCfg
		Registry      struct {
			Target   string
			CacheTTL time.Duration
		}
		Policies map[string]struct {
			Target     string
			MinVersion string
		}
	}
	Events struct {
		Kafka struct {
			Brokers  []string `koanf:"KAFKA_BROKERS"`
			ClientID string   `koanf:"KAFKA_CLIENT_ID"`
			Producer struct {
				Acks        string
				Compression string
				Linger      time.Duration
				MaxInflight int
			}
			Consumer struct {
				GroupId         string
				SessionTimeout  time.Duration
				AutoOffsetReset string
			}
			Topics []string `koanf:"KAFKA_TOPICS"`
		}
	}
	Features struct {
		SSE struct {
			Enabled           bool
			HeartBeatInterval time.Duration
		}
		HttpToGRPC struct {
			StrictValidation bool `koanf:"HTTP_TO_GRPC_STRICT_VALIDATION"`
		}
	}
	Health struct {
		Services       []string
		LivenessPath   string
		ReadinessPath  string
		CheckIntervals struct {
			GrpcTargets time.Duration
			Redis       time.Duration
			Kafka       time.Duration
		}
	}
	Observability struct {
		Log struct {
			Level    string `koanf:"LOG_LEVEL"`
			Json     bool   `koanf:"LOG_JSON"`
			Sampling struct {
				Enabled    bool `koanf:"LOG_SAMPLING_ENABLED"`
				Initial    int
				ThereAfter int
			}
		}
		Tracing struct {
			Enabled     bool   `koanf:"TRACING_ENABLED"`
			Exporter    string `koanf:"TRACING_EXPORTER"`
			Insecure    bool   `koanf:"TRACING_INSECURE"`
			ServiceName string `koanf:"TRACING_SERVICE_NAME"`
		}
		Metrics struct {
			Enabled bool `koanf:"METRICS_ENABLED"`
			Prom    struct {
				Path string `koanf:"PROMETHEUS_PATH"`
				Port int    `koanf:"PROMETHEUS_PORT"`
			}
		}
		Profiling struct {
			Pprof struct {
				Enabled bool   `koanf:"PPROF_ENABLED"`
				Path    string `koanf:"PPROF_PATH"`
			}
		}
	}
	Server struct {
		Http struct {
			Port         string
			ReadTimeout  time.Duration
			WriteTimeout time.Duration
			IdleTimeout  time.Duration
		}
		Cors struct {
			Enabled          bool
			Origins          []string
			Methods          []string
			Headers          []string
			AllowCredentials bool
			MaxAge           time.Duration
		}
		RateLimit struct {
			Enabled bool
		}
	}
	Routes []helper.Routes
}

type Client struct {
	k                   *koanf.Koanf
	root                string
	baseConfigPaths     []string
	overrideConfigPaths []string
}

func NewClient(root string, overrideConfigPaths []string) *Client {
	return &Client{
		k:    koanf.New("."),
		root: root,
		baseConfigPaths: []string{
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "app.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "auth.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "cache.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "downstreams.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "events.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "features.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "health.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "observability.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "server.yaml"),
			filepath.Join("..", "..", "libs", "shared", "configclient", "base", "routes.yaml"),
		},
		overrideConfigPaths: overrideConfigPaths,
	}
}

func (c *Client) Load() *Config {
	for _, bcp := range c.baseConfigPaths {
		if _, err := os.Stat(bcp); err == nil {
			if err := c.k.Load(file.Provider(bcp), yaml.Parser()); err != nil {
				log.Fatalf("error loading config file %s: %v", bcp, err)
			}
		}
	}
	for _, ocp := range c.overrideConfigPaths {
		if _, err := os.Stat(ocp); err == nil {
			if err := c.k.Load(file.Provider(ocp), yaml.Parser()); err != nil {
				log.Fatalf("error loading config file %s: %v", ocp, err)
			}
		}
	}
	// ENV-Overrides (SERVER_HTTP_PORT etc.)
	if err := c.k.Load(env.Provider("", ".", func(s string) string {
		fmt.Printf("Loading env var: %s\n", s)
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	}), nil); err != nil {
		log.Fatalf("error loading env vars: %v", err)
	}

	var cfg Config
	if err := c.k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return &cfg
}

func LoadEnvFile(root string) {
	_ = godotenv.Load(filepath.Join(root, ".env"))
}

func GetenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
