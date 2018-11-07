// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period   time.Duration `config:"period"`
  Classes  []ClassConfig
}

type ClassConfig struct {
	Class       string    `config:"class"`
	Fields      []string  `config:"fields"`
	WhereClause string    `config:"whereclause"`
	ObjectTitle string    `config:"objecttitlecolumn"`
}

var DefaultConfig = Config{
  Period: 1*time.Second,
  Classes:  make([]ClassConfig,0),
}
