package core

import (
	"context"
	"fmt"
	"time"

	"github.com/snakoner/go-crypto-robot/internal/models"
)

var (
	timeframeToSeconds = map[string]int64{
		"1m":  60,
		"3m":  180,
		"5m":  300,
		"15m": 900,
		"30m": 1800,
		"1h":  3600,
		"2h":  7200,
		"4h":  14400,
		"12h": 43200,
		"1d":  86400,
	}
)

// Round time
func roundedValue(prevTime, currentTime int64, timeframe string) (int64, bool) {
	timeFrameSeconds := timeframeToSeconds[timeframe]
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
	tracker.Stat.LastMaxPrice = mp.Price
	tracker.Stat.CurrentStopLoss = mp.Price * (1. - core.Config.MaxStopLoss/100.)
	tracker.Stat.CurrentTakeProfit = mp.Price * (1. + core.Config.MaxTakeProfit/100.)

	core.Logger.Info(fmt.Sprintf("open deal. price: %v", mp.Price))
	core.Logger.Info(fmt.Sprintf("set inital limits. sl: %v tp: %v", tracker.Stat.CurrentStopLoss, tracker.Stat.CurrentTakeProfit))

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

	core.Logger.Info(fmt.Sprintf("close deal. price: %v", mp.Price))

	// [todo] bybit.close_deal
	if err := core.Exchange.CloseTrade(tracker); err != nil {
		return err
	}

	return nil
}

// Change stoploss / takeprofit or close the order if hit
func (core *Core) evaluateDeal(tracker *models.TokenTracker, mp *models.MarketPoint) error {
	lastPrice := tracker.Stat.LastMaxPrice
	currPrice := mp.Price

	if currPrice > lastPrice {
		tracker.Stat.LastMaxPrice = currPrice
	}

	if currPrice <= tracker.Stat.CurrentStopLoss || currPrice >= tracker.Stat.CurrentTakeProfit {
		// close deal
		core.closeTrade(tracker, mp)
	}

	if !core.Config.UseTrailing {
		return nil
	}

	diff := currPrice - lastPrice

	core.Logger.Info(fmt.Sprintf("price: %v", currPrice))

	if diff > 0 {
		newStopLoss := tracker.Stat.CurrentStopLoss + diff
		newTakeProfit := tracker.Stat.CurrentTakeProfit + diff

		if newStopLoss > tracker.Stat.CurrentStopLoss {
			// set stop/take
			tracker.Stat.CurrentStopLoss = newStopLoss
			tracker.Stat.CurrentTakeProfit = newTakeProfit
			// [todo] bybit.set_new_stop_take
			core.Logger.Info(fmt.Sprintf("update limits. sl: %v tp: %v", newStopLoss, newTakeProfit))
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
func (core *Core) trackersStart(tracker *models.TokenTracker, ctx context.Context) error {
	var mp models.MarketPoint
	i := 0

	for {
		select {
		case mp = <-tracker.CurrentPrice:
			fmt.Printf(" [%d] %v\r", i, mp)
			i++
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
				core.Logger.Info("start trading")
				if err := core.openTrade(tracker, &mp); err != nil {
					core.Logger.Error(fmt.Sprintf("cant start trading %s at %s", tracker.Name, core.Config.Exchange))
					return err
				}
			}
			// code
		case <-tracker.Exit:
			core.Logger.Info("close websocket connection")
			core.Logger.Info("try to reconnect")
			go core.Exchange.WsRun(tracker)
		case <-ctx.Done():
			core.Logger.Info("app finished by user")
			tracker.CloseConnection <- true
			return nil
		}
	}
}
