package exchanges

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hirokisan/bybit/v2"
	"github.com/snakoner/go-crypto-robot/internal/errno"
	"github.com/snakoner/go-crypto-robot/internal/models"
)

type BybitExchange struct {
	exchange *Exchange
	client   *bybit.Client
}

func (e *BybitExchange) Connect(apiPublic, apiSecret string) error {
	client := bybit.NewClient().WithAuth("your api key", "your api secret")
	if client == nil {
		return errno.ErrBybitCouldntAuth
	}

	e.client = client

	return nil
}

func unixTimeToDate(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func __reverseSlice[T any](s []T) {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
}

func __intervalToBybitInterval(timeframe string) bybit.Interval {
	switch timeframe {
	case "1m":
		return bybit.Interval("1")
	case "3m":
		return bybit.Interval("3")
	case "5m":
		return bybit.Interval("5")
	case "15m":
		return bybit.Interval("15")
	case "30m":
		return bybit.Interval("30")
	case "1h":
		return bybit.Interval("60")
	case "2h":
		return bybit.Interval("120")
	case "4h":
		return bybit.Interval("240")
	case "6h":
		return bybit.Interval("360")
	case "12h":
		return bybit.Interval("720")
	case "1d":
		return bybit.Interval("D")
	case "1w":
		return bybit.Interval("W")
	case "1mon":
		return bybit.Interval("M")
	}

	return bybit.Interval("Unknown")
}

// [todo] limit, max ret is 200
func (e *BybitExchange) GetKlines(name string, stable string, timeframe string) ([]*models.MarketPoint, error) {
	mPoints := []*models.MarketPoint{}

	param := &bybit.V5GetKlineParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   bybit.SymbolV5(strings.ToUpper(fmt.Sprintf("%s%s", name, stable))),
		Interval: __intervalToBybitInterval(timeframe),
	}

	resp, err := e.client.V5().Market().GetKline(*param)
	if err != nil {
		return mPoints, err
	}

	for _, k := range resp.Result.List {
		time, err := strconv.ParseInt(k.StartTime, 10, 64)
		if err != nil {
			return mPoints, err
		}

		close, err := strconv.ParseFloat(k.Close, 64)
		if err != nil {
			return mPoints, err
		}

		mPoints = append(mPoints, &models.MarketPoint{
			Time:  time / 1000,
			Price: close,
		})
	}
	// because list order from first to last
	__reverseSlice(mPoints)

	return mPoints, nil
}

func (e *BybitExchange) OpenTrade(tracker *models.TokenTracker) error {
	return nil
}

func (e *BybitExchange) CloseTrade(tracker *models.TokenTracker) error {
	return nil
}

func (e *BybitExchange) UpdateLimits(tracker *models.TokenTracker) error {
	return nil
}
