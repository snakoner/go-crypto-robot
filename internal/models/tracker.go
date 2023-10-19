package models

// Chart point
type MarketPoint struct {
	Time  int64
	Price float64
}

// Main data while trade session is on
type Statistics struct {
	CurrentStopLoss   float64
	CurrentTakeProfit float64
	DealActive        bool
	EnterPrice        float64
	EnterTime         int64
	ExitPrice         float64
	ExitTime          int64
	CurrentPrice      float64
	LastMaxPrice      float64
}

// Main trade struct
type TokenTracker struct {
	Name            string `json:"name"`
	Stablecoin      string
	MarketPoints    []*MarketPoint
	CurrentPrice    chan MarketPoint
	Exit            chan bool
	CloseConnection chan bool
	Long            bool
	Stat            Statistics
}

// Returns last time market point in list
func (tracker *TokenTracker) GetLastPoint() *MarketPoint {
	if len(tracker.MarketPoints) == 0 {
		return nil
	}

	return tracker.MarketPoints[len(tracker.MarketPoints)-1]
}

func (tracker *TokenTracker) GetLastPointTime() int64 {
	return tracker.MarketPoints[len(tracker.MarketPoints)-1].Time
}

func (tracker *TokenTracker) GetLastPointPrice() float64 {
	return tracker.MarketPoints[len(tracker.MarketPoints)-1].Price
}

// Add new market point to MP list. Size of list remains the same
func (tracker *TokenTracker) Push(marketPoint MarketPoint) {
	tracker.MarketPoints = append(tracker.MarketPoints, &marketPoint)
	tracker.MarketPoints = tracker.MarketPoints[1:]
}

// Update last time market point with new price value
func (tracker *TokenTracker) Update(marketPoint *MarketPoint) {
	tracker.GetLastPoint().Price = marketPoint.Price
}
