package main

import (
	"github.com/toby1991/go-zero/core/load"
	"github.com/toby1991/go-zero/core/logx"
	"github.com/toby1991/go-zero/tools/goctl/cmd"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
