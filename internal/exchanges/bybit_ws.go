package exchanges

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hirokisan/bybit/v2"
	"github.com/snakoner/go-crypto-robot/internal/models"
)

func (e *BybitExchange) WebSocketRun(tracker *models.TokenTracker) error {
	wsClient := bybit.NewWebsocketClient()
	svc, err := wsClient.Spot().V1().PublicV1()

	defer func() {
		fmt.Printf("Connection closed for %s", tracker.Name)
		svc.Close()
	}()

	if err != nil {
		return err
	}
	_, err = svc.SubscribeTrade(bybit.SymbolSpotBTCUSDT, func(response bybit.SpotWebsocketV1PublicV1TradeResponse) error {
		price, err := strconv.ParseFloat(response.Data[len(response.Data)-1].Price, 64)
		if err != nil {
			tracker.Exit <- true
			return err
		}

		// unix time in second
		timeSecond := int64(response.Data[len(response.Data)-1].Timestamp) / int64(1000)
		tracker.CurrentPrice <- models.MarketPoint{
			Price: price,
			Time:  timeSecond,
		}

		return nil
	})

	if err != nil {
		tracker.Exit <- true
		return err
	}

	svc.Start(context.Background())

	return nil
}
