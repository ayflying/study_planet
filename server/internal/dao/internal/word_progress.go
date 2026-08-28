// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WordProgressDao is the data access object for the table word_progress.
type WordProgressDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  WordProgressColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// WordProgressColumns defines and stores column names for the table word_progress.
type WordProgressColumns struct {
	WordId       string //
	ChildId      string //
	Known        string //
	LastReviewed string //
}

// wordProgressColumns holds the columns for the table word_progress.
var wordProgressColumns = WordProgressColumns{
	WordId:       "word_id",
	ChildId:      "child_id",
	Known:        "known",
	LastReviewed: "last_reviewed",
}

// NewWordProgressDao creates and returns a new DAO object for table data access.
func NewWordProgressDao(handlers ...gdb.ModelHandler) *WordProgressDao {
	return &WordProgressDao{
		group:    "default",
		table:    "word_progress",
		columns:  wordProgressColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *WordProgressDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *WordProgressDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *WordProgressDao) Columns() WordProgressColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *WordProgressDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *WordProgressDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *WordProgressDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
