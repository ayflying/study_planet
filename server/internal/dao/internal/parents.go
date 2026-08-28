// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ParentsDao is the data access object for the table parents.
type ParentsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ParentsColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ParentsColumns defines and stores column names for the table parents.
type ParentsColumns struct {
	Id          string //
	CasdoorSub  string //
	DisplayName string //
	Avatar      string //
	CreatedAt   string //
	LastLoginAt string //
}

// parentsColumns holds the columns for the table parents.
var parentsColumns = ParentsColumns{
	Id:          "id",
	CasdoorSub:  "casdoor_sub",
	DisplayName: "display_name",
	Avatar:      "avatar",
	CreatedAt:   "created_at",
	LastLoginAt: "last_login_at",
}

// NewParentsDao creates and returns a new DAO object for table data access.
func NewParentsDao(handlers ...gdb.ModelHandler) *ParentsDao {
	return &ParentsDao{
		group:    "default",
		table:    "parents",
		columns:  parentsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ParentsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ParentsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ParentsDao) Columns() ParentsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ParentsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ParentsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ParentsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
