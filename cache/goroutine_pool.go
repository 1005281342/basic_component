package cache

import "github.com/panjf2000/ants/v2"

type goroutinePool struct {
	*ants.Pool
}

func newGoroutinePool(size int, optionList ...ants.Option) goroutinePool {
	var pool, err = ants.NewPool(size, optionList...)
	if err != nil {
		panic(err)
	}
	return goroutinePool{
		Pool: pool,
	}
}
