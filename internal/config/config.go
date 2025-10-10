package config

var (
	Protocol              = "tcp"
	Port                  = ":3000"
	MaxConnection         = 20000
	MaxKeyNumber   int    = 10
	EvictionRatio         = 0.1
	EvictionPolicy string = "allkeys-random"
)
