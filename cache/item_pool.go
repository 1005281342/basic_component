package cache

import "sync"

var (
	itemPool = sync.Pool{
		New: func() interface{} {
			return new(item)
		},
	}

	entryPool = sync.Pool{
		New: func() interface{} {
			return new(entry)
		},
	}

	entryWithFreqPool = sync.Pool{
		New: func() interface{} {
			return new(entryWithFreq)
		},
	}

	entryWithHistoryPool = sync.Pool{
		New: func() interface{} {
			return new(entryWithHistory)
		},
	}
)
