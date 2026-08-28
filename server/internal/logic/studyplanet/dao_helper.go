package studyplanet

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"studyplanet/internal/dao"
)

// ctxOf 取 GF 上下文（handler 内部延迟获取，避免每处都传）。
func ctxOf() context.Context { return gctx.New() }

// daoChildren / daoTasks / daoRewards / daoPointsLog 常用表对象快捷入口。
var (
	daoChildren  = dao.Children
	daoTasks     = dao.Tasks
	daoRewards   = dao.Rewards
	daoPointsLog = dao.PointsLog
)

var _ = g.DB // 保留 g.DB 引用，供后续扩展
