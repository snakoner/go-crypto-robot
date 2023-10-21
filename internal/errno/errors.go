package errno

import "errors"

var (
	ErrLengthNotEqual    = errors.New("slices have different length")
	ErrStrategyConfig    = errors.New("cant parse strategy")
	ErrExchangeName      = errors.New("unknown exchange name")
	ErrBybitCouldntAuth  = errors.New("bybit authentication failed")
	ErrBybitNotConnected = errors.New("exchange not connected")
	ErrBybitReconLimit   = errors.New("all reconnection attemtps are over")
)
