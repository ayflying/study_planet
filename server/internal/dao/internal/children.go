// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChildrenDao is the data access object for the table children.
type ChildrenDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ChildrenColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ChildrenColumns defines and stores column names for the table children.
type ChildrenColumns struct {
	Id        string //
	Name      string //
	Username  string //
	Avatar    string //
	Grade     string //
	CreatedAt string //
	Xp        string //
}

// childrenColumns holds the columns for the table children.
var childrenColumns = ChildrenColumns{
	Id:        "id",
	Name:      "name",
	Username:  "username",
	Avatar:    "avatar",
	Grade:     "grade",
	CreatedAt: "created_at",
	Xp:        "xp",
}

// NewChildrenDao creates and returns a new DAO object for table data access.
func NewChildrenDao(handlers ...gdb.ModelHandler) *ChildrenDao {
	return &ChildrenDao{
		group:    "default",
		table:    "children",
		columns:  childrenColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ChildrenDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ChildrenDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ChildrenDao) Columns() ChildrenColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ChildrenDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ChildrenDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ChildrenDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
