package exchanges

import (
	"github.com/snakoner/go-crypto-robot/internal/models"
)

type ExchangeI interface {
	Connect() error
	GetWalletBalance() error
	GetKlines(string, string, string) ([]*models.MarketPoint, error)
	WebSocketRun(tracker *models.TokenTracker) error
	OpenTrade(tracker *models.TokenTracker) error
	CloseTrade(tracker *models.TokenTracker) error
	UpdateLimits(tracker *models.TokenTracker) error
}
