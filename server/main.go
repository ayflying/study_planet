package main

import (
	_ "studyplanet/internal/logic"

	"studyplanet/internal/cmd"

	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
