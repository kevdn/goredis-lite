package core

import "goredis-lite/internal/data_structure"

var (
	dictStore *data_structure.Dict
	zsetStore map[string]*data_structure.SortedSet
	setStore  map[string]*data_structure.SimpleSet
	cmsStore  map[string]*data_structure.CMS
	bloomStore map[string]*data_structure.Bloom
)

func init() {
	dictStore = data_structure.CreateDict()
	zsetStore = make(map[string]*data_structure.SortedSet)
	setStore = make(map[string]*data_structure.SimpleSet)
	cmsStore = make(map[string]*data_structure.CMS)
	bloomStore = make(map[string]*data_structure.Bloom)
}
