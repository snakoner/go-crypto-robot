package core

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/snakoner/go-crypto-robot/internal/algorithm"
	"github.com/snakoner/go-crypto-robot/internal/errno"
	"github.com/snakoner/go-crypto-robot/internal/exchanges"
	"github.com/snakoner/go-crypto-robot/internal/models"
)

type Core struct {
	Config        *Config
	Logger        *logrus.Logger
	Strategy      *algorithm.Strategy
	Exchange      exchanges.ExchangeI
	TokenTrackers []*models.TokenTracker
}

// Create trading core and setup, call from main
func New(config *Config) (*Core, error) {
	core := &Core{
		Config: config,
		Logger: logrus.New(),
	}

	// strategy setup
	strategy := algorithm.NewStrategy(config.Algos)
	if strategy == nil {
		return core, errno.ErrStrategyConfig
	}

	core.Strategy = strategy

	core.Logger.Debug(fmt.Sprintf("Algos: %s", core.Strategy.String()))

	// new exchange
	switch config.Exchange {
	case "bybit":
		core.Exchange = exchanges.NewBybit(config.PublicKey, config.PrivateKey)
	default:
		return core, errno.ErrExchangeName
	}

	// exchange connect
	if err := core.Exchange.Connect(); err != nil {
		core.Logger.Error(err)
	}

	core.Logger.Debug(fmt.Sprintf("%s conn success", config.Exchange))
	core.Logger.Debug("Core run success")

	return core, nil
}

// Start trading core, call from main
func (core *Core) Start(ctx context.Context) error {
	for _, name := range core.Config.Coins {
		var tokenTracker *models.TokenTracker
		mp, err := core.Exchange.GetKlines(name, core.Config.Stablecoin, core.Config.Timeframe)
		if err != nil {
			return err
		}

		tokenTracker = &models.TokenTracker{
			MarketPoints: mp,
			Name:         name,
			Stablecoin:   core.Config.Stablecoin,
			CurrentPrice: make(chan models.MarketPoint),
			Exit:         make(chan bool),
		}

		core.TokenTrackers = append(core.TokenTrackers, tokenTracker)

		core.Logger.Debug(fmt.Sprintf("Size of initial data for %s : %d",
			name+core.Config.Stablecoin,
			len(tokenTracker.MarketPoints)))

		go core.Exchange.WebSocketRun(tokenTracker)
		go core.trackersStart(ctx)
	}

	<-ctx.Done()

	return nil
}
