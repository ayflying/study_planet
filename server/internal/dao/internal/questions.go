// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// QuestionsDao is the data access object for the table questions.
type QuestionsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  QuestionsColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// QuestionsColumns defines and stores column names for the table questions.
type QuestionsColumns struct {
	Id          string //
	Subject     string //
	Grade       string //
	Topic       string //
	Qtype       string //
	Passage     string //
	Question    string //
	Options     string //
	Answer      string //
	Explanation string //
	Difficulty  string //
	Source      string //
	ContentHash string //
	Enabled     string //
	CreatedAt   string //
}

// questionsColumns holds the columns for the table questions.
var questionsColumns = QuestionsColumns{
	Id:          "id",
	Subject:     "subject",
	Grade:       "grade",
	Topic:       "topic",
	Qtype:       "qtype",
	Passage:     "passage",
	Question:    "question",
	Options:     "options",
	Answer:      "answer",
	Explanation: "explanation",
	Difficulty:  "difficulty",
	Source:      "source",
	ContentHash: "content_hash",
	Enabled:     "enabled",
	CreatedAt:   "created_at",
}

// NewQuestionsDao creates and returns a new DAO object for table data access.
func NewQuestionsDao(handlers ...gdb.ModelHandler) *QuestionsDao {
	return &QuestionsDao{
		group:    "default",
		table:    "questions",
		columns:  questionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *QuestionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *QuestionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *QuestionsDao) Columns() QuestionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *QuestionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *QuestionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *QuestionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
