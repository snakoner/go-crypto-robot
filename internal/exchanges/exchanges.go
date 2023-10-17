package exchanges

import (
	"github.com/snakoner/go-crypto-robot/internal/models"
)

type BybitExchangeI interface {
	Connect(string, string) error
	GetKlines(string, string, string) ([]*models.MarketPoint, error)
	WebSocketRun(tracker *models.TokenTracker) error
	OpenTrade(tracker *models.TokenTracker) error
	CloseTrade(tracker *models.TokenTracker) error
	UpdateLimits(tracker *models.TokenTracker) error
}

// type BinanceExchangeI interface {
// 	Get()
// }

type ExchangeI interface {
	Bybit() BybitExchangeI
	// Binance() BinanceExchangeI
}

type Exchange struct {
	bybit *BybitExchange
}

func (ex *Exchange) Bybit() BybitExchangeI {
	if ex.bybit != nil {
		return ex.bybit
	}

	ex.bybit = &BybitExchange{
		exchange: ex,
	}

	return ex.bybit
}
