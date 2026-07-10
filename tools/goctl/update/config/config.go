package config

import "github.com/toby1991/go-zero/core/logx"

// Config defines a service configure for goctl update
type Config struct {
	logx.LogConf
	ListenOn string
	FileDir  string
	FilePath string
}
