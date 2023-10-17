package errno

import "errors"

var (
	ErrLengthNotEqual   = errors.New("slices have different length")
	ErrStrategyConfig   = errors.New("cant parse strategy")
	ErrBybitCouldntAuth = errors.New("bybit authentication failed")
)
