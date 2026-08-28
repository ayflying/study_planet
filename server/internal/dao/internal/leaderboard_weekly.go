// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LeaderboardWeeklyDao is the data access object for the table leaderboard_weekly.
type LeaderboardWeeklyDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  LeaderboardWeeklyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// LeaderboardWeeklyColumns defines and stores column names for the table leaderboard_weekly.
type LeaderboardWeeklyColumns struct {
	Id        string //
	WeekKey   string //
	ChildId   string //
	Xp        string //
	UpdatedAt string //
}

// leaderboardWeeklyColumns holds the columns for the table leaderboard_weekly.
var leaderboardWeeklyColumns = LeaderboardWeeklyColumns{
	Id:        "id",
	WeekKey:   "week_key",
	ChildId:   "child_id",
	Xp:        "xp",
	UpdatedAt: "updated_at",
}

// NewLeaderboardWeeklyDao creates and returns a new DAO object for table data access.
func NewLeaderboardWeeklyDao(handlers ...gdb.ModelHandler) *LeaderboardWeeklyDao {
	return &LeaderboardWeeklyDao{
		group:    "default",
		table:    "leaderboard_weekly",
		columns:  leaderboardWeeklyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LeaderboardWeeklyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LeaderboardWeeklyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LeaderboardWeeklyDao) Columns() LeaderboardWeeklyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LeaderboardWeeklyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LeaderboardWeeklyDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *LeaderboardWeeklyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
