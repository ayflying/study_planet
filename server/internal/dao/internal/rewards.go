// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RewardsDao is the data access object for the table rewards.
type RewardsDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of the current DAO.
	columns RewardsColumns // columns contains all the column names of Table for convenient usage.
}

// RewardsColumns defines and stores column names for the table rewards.
type RewardsColumns struct {
	Id         string //
	Name       string //
	CostPoints string //
	Status     string //
}

// rewardsColumns holds the columns for the table rewards.
var rewardsColumns = RewardsColumns{
	Id:         "id",
	Name:       "name",
	CostPoints: "cost_points",
	Status:     "status",
}

// NewRewardsDao creates and returns a new DAO object for table data access.
func NewRewardsDao() *RewardsDao {
	return &RewardsDao{
		group:   "default",
		table:   "rewards",
		columns: rewardsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RewardsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RewardsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RewardsDao) Columns() RewardsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RewardsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RewardsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *RewardsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
