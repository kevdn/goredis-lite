package config

var (
	Protocol          = "tcp"
	Port              = ":3000"
	MaxConnection     = 20000
	MaxKeyNumber  int = 1000000
)

var (
	EvictionRatio         = 0.1
	EvictionPolicy string = "allkeys-lru"
)

var (
	EpoolMaxSize       = 16
	EpoolLruSampleSize = 5
)

var ListenerNumber int = 2
