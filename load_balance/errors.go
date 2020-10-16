package load_balance

import "errors"

var (
	ErrParamLeastOne    = errors.New("param len 1 at least")
	ErrParamLeastTwo    = errors.New("param len 2 at least")
	ErrNodeListIsEmpty  = errors.New("ErrNodeListIsEmpty")
	ErrNodeNotAvailable = errors.New("ErrNodeNotAvailable")
)
