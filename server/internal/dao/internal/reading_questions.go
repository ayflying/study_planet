// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReadingQuestionsDao is the data access object for the table reading_questions.
type ReadingQuestionsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  ReadingQuestionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// ReadingQuestionsColumns defines and stores column names for the table reading_questions.
type ReadingQuestionsColumns struct {
	Id        string //
	ReadingId string //
	Question  string //
	OptionA   string //
	OptionB   string //
	OptionC   string //
	OptionD   string //
	Answer    string //
}

// readingQuestionsColumns holds the columns for the table reading_questions.
var readingQuestionsColumns = ReadingQuestionsColumns{
	Id:        "id",
	ReadingId: "reading_id",
	Question:  "question",
	OptionA:   "option_a",
	OptionB:   "option_b",
	OptionC:   "option_c",
	OptionD:   "option_d",
	Answer:    "answer",
}

// NewReadingQuestionsDao creates and returns a new DAO object for table data access.
func NewReadingQuestionsDao(handlers ...gdb.ModelHandler) *ReadingQuestionsDao {
	return &ReadingQuestionsDao{
		group:    "default",
		table:    "reading_questions",
		columns:  readingQuestionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReadingQuestionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReadingQuestionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReadingQuestionsDao) Columns() ReadingQuestionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReadingQuestionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReadingQuestionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReadingQuestionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
