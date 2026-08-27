package main

import (
	_ "studyplanet/examples/gf-init-example/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"studyplanet/examples/gf-init-example/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
