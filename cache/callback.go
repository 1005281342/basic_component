package cache

// 当元素从缓存中移除时进行回调
type EvictCallback func(key interface{}, value interface{})
