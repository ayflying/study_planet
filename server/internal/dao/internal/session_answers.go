// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SessionAnswersDao is the data access object for the table session_answers.
type SessionAnswersDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  SessionAnswersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// SessionAnswersColumns defines and stores column names for the table session_answers.
type SessionAnswersColumns struct {
	Id         string //
	SessionId  string //
	RefId      string //
	Correct    string //
	Points     string //
	Combo      string //
	AnsweredAt string //
}

// sessionAnswersColumns holds the columns for the table session_answers.
var sessionAnswersColumns = SessionAnswersColumns{
	Id:         "id",
	SessionId:  "session_id",
	RefId:      "ref_id",
	Correct:    "correct",
	Points:     "points",
	Combo:      "combo",
	AnsweredAt: "answered_at",
}

// NewSessionAnswersDao creates and returns a new DAO object for table data access.
func NewSessionAnswersDao(handlers ...gdb.ModelHandler) *SessionAnswersDao {
	return &SessionAnswersDao{
		group:    "default",
		table:    "session_answers",
		columns:  sessionAnswersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SessionAnswersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SessionAnswersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SessionAnswersDao) Columns() SessionAnswersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SessionAnswersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SessionAnswersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SessionAnswersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
