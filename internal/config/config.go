package config

var (
	Protocol          = "tcp"
	Port              = ":3000"
	MaxConnection     = 20000
	MaxKeyNumber  int = 10
)

var (
	EvictionRatio         = 0.1
	EvictionPolicy string = "allkeys-lru"
)

var (
	EpoolMaxSize       = 5
	EpoolLruSampleSize = 5
)
