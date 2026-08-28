// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubjectsDao is the data access object for the table subjects.
type SubjectsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SubjectsColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SubjectsColumns defines and stores column names for the table subjects.
type SubjectsColumns struct {
	Id        string //
	Code      string //
	Name      string //
	Icon      string //
	Color     string //
	MinGrade  string //
	MaxGrade  string //
	Sort      string //
	Enabled   string //
	CreatedAt string //
}

// subjectsColumns holds the columns for the table subjects.
var subjectsColumns = SubjectsColumns{
	Id:        "id",
	Code:      "code",
	Name:      "name",
	Icon:      "icon",
	Color:     "color",
	MinGrade:  "min_grade",
	MaxGrade:  "max_grade",
	Sort:      "sort",
	Enabled:   "enabled",
	CreatedAt: "created_at",
}

// NewSubjectsDao creates and returns a new DAO object for table data access.
func NewSubjectsDao(handlers ...gdb.ModelHandler) *SubjectsDao {
	return &SubjectsDao{
		group:    "default",
		table:    "subjects",
		columns:  subjectsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SubjectsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SubjectsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SubjectsDao) Columns() SubjectsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SubjectsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SubjectsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SubjectsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
