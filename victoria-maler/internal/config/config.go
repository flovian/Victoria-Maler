package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	DBPath     string
	Port       string
	Env        string
	JWTSecret  string
	UploadDir  string
	Bitcoin    BitcoinConfig
}

type BitcoinConfig struct {
	Enabled  bool
	RPCURL   string
	RPCUser  string
	RPCPass  string
	Network  string
	MaxRetry int
}

func Load() *Config {
	loadDotEnv()

	cfg := &Config{
		DBPath:    getEnv("DB_PATH", "storage/database.db"),
		Port:      getEnv("PORT", "8080"),
		Env:       getEnv("ENV", "development"),
		JWTSecret: getEnv("JWT_SECRET", "victoria-maler-dev-secret-change-me"),
		UploadDir: getEnv("UPLOAD_DIR", "web/static/uploads"),
	}

	cfg.Bitcoin = BitcoinConfig{
		Enabled:  getEnv("BITCOIN_ENABLED", "false") == "true",
		RPCURL:   getEnv("BITCOIN_RPC_URL", "http://127.0.0.1:18443"),
		RPCUser:  getEnv("BITCOIN_RPC_USER", "bitcoin"),
		RPCPass:  getEnv("BITCOIN_RPC_PASS", "bitcoin"),
		Network:  getEnv("BITCOIN_NETWORK", "regtest"),
		MaxRetry: 2,
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadDotEnv() {
	envPath := ".env"
	if _, err := os.Stat(envPath); err != nil {
		return
	}

	file, err := os.Open(envPath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}
