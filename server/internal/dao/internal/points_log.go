// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointsLogDao is the data access object for the table points_log.
type PointsLogDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of the current DAO.
	columns PointsLogColumns // columns contains all the column names of Table for convenient usage.
}

// PointsLogColumns defines and stores column names for the table points_log.
type PointsLogColumns struct {
	Id        string //
	ChildId   string //
	Delta     string //
	Reason    string //
	CreatedAt string //
}

// pointsLogColumns holds the columns for the table points_log.
var pointsLogColumns = PointsLogColumns{
	Id:        "id",
	ChildId:   "child_id",
	Delta:     "delta",
	Reason:    "reason",
	CreatedAt: "created_at",
}

// NewPointsLogDao creates and returns a new DAO object for table data access.
func NewPointsLogDao() *PointsLogDao {
	return &PointsLogDao{
		group:   "default",
		table:   "points_log",
		columns: pointsLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointsLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointsLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointsLogDao) Columns() PointsLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointsLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointsLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *PointsLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
