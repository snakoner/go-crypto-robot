package bybit

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hirokisan/bybit/v2"
	"github.com/snakoner/go-crypto-robot/internal/errno"
	"github.com/snakoner/go-crypto-robot/internal/models"
)

var (
	intervalToBybitInterval = map[string]bybit.Interval{
		"1m":   bybit.Interval("1"),
		"3m":   bybit.Interval("3"),
		"5m":   bybit.Interval("5"),
		"15m":  bybit.Interval("15"),
		"30m":  bybit.Interval("30"),
		"1h":   bybit.Interval("60"),
		"2h":   bybit.Interval("120"),
		"4h":   bybit.Interval("240"),
		"6h":   bybit.Interval("360"),
		"12h":  bybit.Interval("720"),
		"1d":   bybit.Interval("D"),
		"1w":   bybit.Interval("W"),
		"1mon": bybit.Interval("M"),
	}
)

// Main Bybit exchange structure
type BybitExchange struct {
	privateKey string
	publicKey  string
	client     *bybit.Client
}

func NewBybit(publicKey, privateKey string) *BybitExchange {
	bybit := &BybitExchange{
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	fmt.Println(privateKey, publicKey)

	return bybit
}

// Connect to bybit spot, save connection to BybitExchange
func (e *BybitExchange) Connect() error {
	client := bybit.NewClient().WithAuth(e.publicKey, e.privateKey)
	if client == nil {
		return errno.ErrBybitCouldntAuth
	}

	e.client = client

	return nil
}

func unixTimeToDate(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func reverseSlice[T any](s []T) {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
}

// [todo] limit, max ret is 200
func (e *BybitExchange) GetKlines(name string, stable string, timeframe string) ([]*models.MarketPoint, error) {
	mPoints := []*models.MarketPoint{}

	param := &bybit.V5GetKlineParam{
		Category: bybit.CategoryV5Spot,
		Symbol:   bybit.SymbolV5(strings.ToUpper(fmt.Sprintf("%s%s", name, stable))),
		Interval: intervalToBybitInterval[timeframe],
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
	reverseSlice(mPoints)

	return mPoints, nil
}

func (e *BybitExchange) isConnected() bool {
	return e.client != nil
}

func (e *BybitExchange) GetBalance() (*models.Assets, error) {
	assets := &models.Assets{}
	if !e.isConnected() {
		return assets, errno.ErrBybitNotConnected
	}

	client := bybit.NewClient().WithAuth(e.publicKey, e.privateKey)
	resp, err := client.V5().Account().GetWalletBalance(bybit.AccountTypeUnified, nil)
	if err != nil {
		fmt.Println("err", err)
		return assets, err
	}

	for _, val := range resp.Result.List {
		assets.TotalUsdValue, err = strconv.ParseFloat(val.TotalEquity, 64)
		if err != nil {
			return assets, err
		}

		for _, c := range val.Coin {
			amount, err := strconv.ParseFloat(c.Equity, 64)
			if err != nil {
				return assets, err
			}

			usdValue, err := strconv.ParseFloat(c.UsdValue, 64)
			if err != nil {
				return assets, err
			}

			assets.Tokens = append(assets.Tokens, models.TokenAsset{
				Name:     string(c.Coin),
				Amount:   amount,
				UsdValue: usdValue,
			})
		}
	}

	return assets, nil
}

// Bybit: open market order
func (e *BybitExchange) OpenTrade(tracker *models.TokenTracker) error {
	return nil
}

// Bybit: close market order
func (e *BybitExchange) CloseTrade(tracker *models.TokenTracker) error {
	return nil
}

// Bybit: change stop loss and take profit
func (e *BybitExchange) UpdateLimits(tracker *models.TokenTracker) error {
	return nil
}
