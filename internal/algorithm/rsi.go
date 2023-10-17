package algorithm

import "github.com/snakoner/go-crypto-robot/internal/models"

const (
	rsiLen        = 14
	rsiLongLimit  = 30.0
	rsiShortLimit = 70.0
)

func Rsi(mp []*models.MarketPoint) (bool, bool) {
	u := make([]float64, len(mp)-1)
	d := make([]float64, len(mp)-1)
	r := make([]float64, len(mp)-1)

	for i := 0; i < len(mp)-1; i++ {
		if mp[i+1].Price >= mp[i].Price {
			u[i] = mp[i+1].Price - mp[i].Price
		} else {
			d[i] = mp[i].Price - mp[i+1].Price
		}
	}

	emaU := Ema(u, rsiLen)
	emaD := Ema(d, rsiLen)

	for i := 0; i < len(mp)-1; i++ {
		r[i] = 100.0
		if emaD[i] != 0.0 {
			r[i] = 100.0 - 100.0/(1+emaU[i]/emaD[i])
		}
	}

	lastValue := r[len(r)-1]

	return lastValue < rsiLongLimit, lastValue > rsiShortLimit
}
