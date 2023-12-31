package exchanges

import (
	"github.com/snakoner/go-crypto-robot/internal/models"
)

type Exchange interface {
	Connect() error
	GetBalance() (*models.Assets, error)
	GetKlines(string, string, string) ([]*models.MarketPoint, error)
	WsRun(tracker *models.TokenTracker) error
	OpenTrade(tracker *models.TokenTracker) error
	CloseTrade(tracker *models.TokenTracker) error
	UpdateLimits(tracker *models.TokenTracker) error
}
