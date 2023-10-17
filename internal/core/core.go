package core

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snakoner/go-crypto-robot/internal/errno"
	"github.com/snakoner/go-crypto-robot/internal/exchanges"
	"github.com/snakoner/go-crypto-robot/internal/models"
)

type Core struct {
	Config        *Config
	Logger        *logrus.Logger
	Strategy      *models.Strategy
	Exch          *exchanges.Exchange
	TokenTrackers []*models.TokenTracker
}

func New(config *Config) (*Core, error) {
	core := &Core{
		Config: config,
		Logger: logrus.New(),
	}

	fmt.Println(config.Algos)
	// strategy setup
	strategy := models.NewStrategy(config.Algos)
	if strategy == nil {
		return core, errno.ErrStrategyConfig
	}
	core.Strategy = strategy

	core.Logger.Info(fmt.Sprintf("Algos: %s", core.Strategy.String()))

	// exchange setup
	core.Exch = new(exchanges.Exchange)
	if err := core.Exch.Bybit().Connect(
		core.Config.BybitPublicKey,
		core.Config.BybitPrivateKey); err != nil {
		core.Logger.Error(err)
	}

	core.Logger.Info(fmt.Sprintf("%s conn success", config.Exchange))

	core.Logger.Info("Core run success")

	return core, nil
}

func (core *Core) Start() error {
	for _, name := range core.Config.Coins {
		var tokenTracker *models.TokenTracker
		mp, err := core.Exch.Bybit().GetKlines(name, core.Config.Stablecoin, core.Config.Timeframe)
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

		core.Logger.Info(fmt.Sprintf("Size of initial data for %s : %d",
			name+core.Config.Stablecoin,
			len(tokenTracker.MarketPoints)))

		go core.Exch.Bybit().WebSocketRun(tokenTracker)
		go core.TrackersStart()
	}

	time.Sleep(100 * time.Second)

	return nil
}
