package main

import (
	"github.com/bilalcaliskan/varnish-cache-invalidator/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
}
