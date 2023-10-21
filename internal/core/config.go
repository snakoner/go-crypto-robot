package core

// Initial core config
type Config struct {
	IsSpot                 bool     `toml:"is_spot"`
	LogFilename            string   `toml:"log_filename"`
	LogLevel               string   `toml:"log_level"`
	Stablecoin             string   `toml:"stablecoin"`
	Coins                  []string `toml:"coins"`
	Algos                  []string `toml:"algos"`
	Timeframe              string   `toml:"timeframe"`
	Exchange               string   `toml:"exchange"`
	PrivateKey             string   `toml:"private_key"`
	PublicKey              string   `toml:"public_key"`
	MaxStopLoss            float64  `toml:"max_stop_loss"`
	MaxTakeProfit          float64  `toml:"max_take_profit"`
	UseTrailing            bool     `toml:"use_trailing"`
	WsReconnectionAttempts int      `toml:"ws_reconnection_attempts"`
}

// New ...
func NewConfig() (*Config, error) {
	return &Config{}, nil
}
