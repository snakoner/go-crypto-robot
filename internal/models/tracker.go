package models

type MarketPoint struct {
	Time  int64
	Price float64
}

type Statistics struct {
	CurrentStopLoss   float64
	CurrentTakeProfit float64
	DealActive        bool
	EnterPrice        float64
	EnterTime         int64
	ExitPrice         float64
	ExitTime          int64
	CurrentPrice      float64
}

type TokenTracker struct {
	Name         string `json:"name"`
	Stablecoin   string
	MarketPoints []*MarketPoint
	CurrentPrice chan MarketPoint
	Exit         chan bool
	Long         bool
	Stat         Statistics
}

func (tracker *TokenTracker) GetLastPoint() *MarketPoint {
	if len(tracker.MarketPoints) == 0 {
		return nil
	}

	return tracker.MarketPoints[len(tracker.MarketPoints)-1]
}

func (tracker *TokenTracker) Push(marketPoint MarketPoint) {
	tracker.MarketPoints = append(tracker.MarketPoints, &marketPoint)
	tracker.MarketPoints = tracker.MarketPoints[1:]
}

func (tracker *TokenTracker) Update(marketPoint *MarketPoint) {
	tracker.GetLastPoint().Price = marketPoint.Price
}
