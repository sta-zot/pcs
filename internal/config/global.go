package config

import (
	"sync"
)

var (
	GlobalConfig *Config
	once         sync.Once
)

func init() {
	cfg := &Config{}
	var err error
	once.Do(func() {
		cfg, err = NewConfig()
		if err != nil {
			panic(err)
		}
		GlobalConfig = cfg
	})
}
