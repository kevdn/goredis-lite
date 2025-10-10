package core

import "goredis-lite/internal/data_structure"

var (
	dictStore *data_structure.Dict
	zsetStore map[string]*data_structure.SortedSet
	setStore  map[string]*data_structure.SimpleSet
)

func init() {
	dictStore = data_structure.CreateDict()
	zsetStore = make(map[string]*data_structure.SortedSet)
	setStore = make(map[string]*data_structure.SimpleSet)
}
