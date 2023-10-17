package core

type Config struct {
	IsSpot        bool     `toml:"is_spot"`
	Stablecoin    string   `toml:"stablecoin"`
	Coins         []string `toml:"coins"`
	Algos         []string `toml:"algos"`
	Timeframe     string   `toml:"timeframe"`
	Exchange      string   `toml:"exchange"`
	PrivateKey    string   `toml:"private_key"`
	PublicKey     string   `toml:"public_key"`
	MaxStopLoss   float64  `toml:"max_stop_loss"`
	MaxTakeProfit float64  `toml:"max_take_profit"`
	UseTrailing   bool     `toml:"use_trailing"`
}

func NewConfig() (*Config, error) {
	return &Config{}, nil
}
