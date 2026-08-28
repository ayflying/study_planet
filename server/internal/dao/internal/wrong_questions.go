// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WrongQuestionsDao is the data access object for the table wrong_questions.
type WrongQuestionsDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  WrongQuestionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// WrongQuestionsColumns defines and stores column names for the table wrong_questions.
type WrongQuestionsColumns struct {
	Id             string //
	ChildId        string //
	Subject        string //
	RefId          string //
	WrongCount     string //
	Resolved       string //
	LastWrongAt    string //
	LastReviewedAt string //
	CreatedAt      string //
}

// wrongQuestionsColumns holds the columns for the table wrong_questions.
var wrongQuestionsColumns = WrongQuestionsColumns{
	Id:             "id",
	ChildId:        "child_id",
	Subject:        "subject",
	RefId:          "ref_id",
	WrongCount:     "wrong_count",
	Resolved:       "resolved",
	LastWrongAt:    "last_wrong_at",
	LastReviewedAt: "last_reviewed_at",
	CreatedAt:      "created_at",
}

// NewWrongQuestionsDao creates and returns a new DAO object for table data access.
func NewWrongQuestionsDao(handlers ...gdb.ModelHandler) *WrongQuestionsDao {
	return &WrongQuestionsDao{
		group:    "default",
		table:    "wrong_questions",
		columns:  wrongQuestionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *WrongQuestionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *WrongQuestionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *WrongQuestionsDao) Columns() WrongQuestionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *WrongQuestionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *WrongQuestionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *WrongQuestionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
