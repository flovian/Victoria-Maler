package config

type Config struct {
    DBPath string
    Port   string
    Env    string
}

func Load() *Config {
    return &Config{
        DBPath: "storage/database.db",
        Port:   "8080",
        Env:    "development",
    }
}
