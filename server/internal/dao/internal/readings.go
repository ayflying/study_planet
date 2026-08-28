// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReadingsDao is the data access object for the table readings.
type ReadingsDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of the current DAO.
	columns ReadingsColumns // columns contains all the column names of Table for convenient usage.
}

// ReadingsColumns defines and stores column names for the table readings.
type ReadingsColumns struct {
	Id      string //
	Title   string //
	Content string //
	Level   string //
}

// readingsColumns holds the columns for the table readings.
var readingsColumns = ReadingsColumns{
	Id:      "id",
	Title:   "title",
	Content: "content",
	Level:   "level",
}

// NewReadingsDao creates and returns a new DAO object for table data access.
func NewReadingsDao() *ReadingsDao {
	return &ReadingsDao{
		group:   "default",
		table:   "readings",
		columns: readingsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReadingsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReadingsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReadingsDao) Columns() ReadingsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReadingsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReadingsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *ReadingsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
