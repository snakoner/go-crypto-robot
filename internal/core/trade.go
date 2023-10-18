package core

import (
	"context"
	"fmt"
	"time"

	"github.com/snakoner/go-crypto-robot/internal/models"
)

func timeframeToSeconds(timeframe string) int64 {
	switch timeframe {
	case "1m":
		return 60
	case "3m":
		return 180
	case "5m":
		return 300
	case "15m":
		return 900
	case "30m":
		return 1800
	case "1h":
		return 3600
	case "2h":
		return 7200
	case "4h":
		return 14400
	case "12h":
		return 43200
	case "1d":
		return 86400
	}

	return 0
}

// Round time
func roundedValue(prevTime, currentTime int64, timeframe string) (int64, bool) {
	timeFrameSeconds := timeframeToSeconds(timeframe)
	// round to minute
	currentTime = currentTime - currentTime%60
	if currentTime-prevTime >= timeFrameSeconds {
		return currentTime, true
	}

	return currentTime, false
}

// If next time > timeframe then push market point, else update price value
func (core *Core) updateMarketPoints(tracker *models.TokenTracker, currentMarketPoint *models.MarketPoint) {
	prevMarketPoint := tracker.GetLastPoint()
	roundedTime, needPush := roundedValue(int64(prevMarketPoint.Time), int64(currentMarketPoint.Time), core.Config.Timeframe)
	currentMarketPoint.Time = roundedTime

	if needPush {
		tracker.Push(*currentMarketPoint)
	} else {
		tracker.Update(currentMarketPoint)
	}
}

// Wrapper for open order
func (core *Core) openTrade(tracker *models.TokenTracker, mp *models.MarketPoint) error {
	tracker.Stat.DealActive = true
	tracker.Stat.EnterPrice = mp.Price
	tracker.Stat.EnterTime = mp.Time
	tracker.Stat.CurrentPrice = mp.Price
	tracker.Stat.CurrentStopLoss = mp.Price * (1. - core.Config.MaxStopLoss/100.)
	tracker.Stat.CurrentTakeProfit = mp.Price * (1. + core.Config.MaxTakeProfit/100.)

	// [todo] bybit.close_deal
	if err := core.Exchange.OpenTrade(tracker); err != nil {
		return err
	}

	// do initial stop/take
	if err := core.Exchange.UpdateLimits(tracker); err != nil {
		return err
	}

	return nil
}

// Wrapper for close order
func (core *Core) closeTrade(tracker *models.TokenTracker, mp *models.MarketPoint) error {
	tracker.Stat.ExitPrice = mp.Price
	tracker.Stat.ExitTime = time.Now().Unix()
	tracker.Stat.DealActive = false

	// [todo] bybit.close_deal
	if err := core.Exchange.CloseTrade(tracker); err != nil {
		return err
	}

	return nil
}

// Change stoploss / takeprofit or close the order if hit
func (core *Core) evaluateDeal(tracker *models.TokenTracker, mp *models.MarketPoint) error {
	if !core.Config.UseTrailing {
		return nil
	}

	lastPrice := tracker.Stat.CurrentPrice
	currPrice := mp.Price

	if currPrice <= tracker.Stat.CurrentStopLoss || currPrice >= tracker.Stat.CurrentTakeProfit {
		// close deal
		core.closeTrade(tracker, mp)
	}

	diff := currPrice - lastPrice

	if diff > 0 {
		newStopLoss := tracker.Stat.CurrentStopLoss + diff
		newTakeProfit := tracker.Stat.CurrentTakeProfit + diff
		if newStopLoss > tracker.Stat.CurrentStopLoss {
			// set stop/take
			tracker.Stat.CurrentStopLoss = newStopLoss
			tracker.Stat.CurrentTakeProfit = newTakeProfit
			// [todo] bybit.set_new_stop_take
			if err := core.Exchange.UpdateLimits(tracker); err != nil {
				if err = core.Exchange.CloseTrade(tracker); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Run as goroutine. Fetches the value from market point channel and evaluate stratefy
func (core *Core) trackersStart(ctx context.Context) error {
	var mp models.MarketPoint

	tracker := core.TokenTrackers[0]
	time := tracker.GetLastPoint().Time
	price := tracker.GetLastPoint().Price

	fmt.Println(mp, time, price)

	for {
		select {
		case mp = <-tracker.CurrentPrice:
			fmt.Println(mp)
			startTrade := false
			if !tracker.Stat.DealActive {
				core.updateMarketPoints(tracker, &mp)
				startTrade = core.Strategy.Calculate(tracker)
			} else {
				// recalculare stop loss and sell if hit
				if err := core.evaluateDeal(tracker, &mp); err != nil {
					return err
				}
			}

			// start trade over here
			if startTrade {
				if err := core.openTrade(tracker, &mp); err != nil {
					core.Logger.Error(fmt.Sprintf("cant start trading %s at %s", tracker.Name, core.Config.Exchange))
					return err
				}
			}
			// code
		case <-tracker.Exit:
			core.Logger.Info("close websocket connection")
			return nil
		case <-ctx.Done():
			core.Logger.Info("app finished by user")
			tracker.Exit <- true
			return nil
		}
	}
}
