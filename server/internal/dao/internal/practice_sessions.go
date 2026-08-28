// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PracticeSessionsDao is the data access object for the table practice_sessions.
type PracticeSessionsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  PracticeSessionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// PracticeSessionsColumns defines and stores column names for the table practice_sessions.
type PracticeSessionsColumns struct {
	Id         string //
	ChildId    string //
	Subject    string //
	Level      string //
	Total      string //
	Correct    string //
	MaxCombo   string //
	Bonus      string //
	Stars      string //
	Finished   string //
	CreatedAt  string //
	FinishedAt string //
}

// practiceSessionsColumns holds the columns for the table practice_sessions.
var practiceSessionsColumns = PracticeSessionsColumns{
	Id:         "id",
	ChildId:    "child_id",
	Subject:    "subject",
	Level:      "level",
	Total:      "total",
	Correct:    "correct",
	MaxCombo:   "max_combo",
	Bonus:      "bonus",
	Stars:      "stars",
	Finished:   "finished",
	CreatedAt:  "created_at",
	FinishedAt: "finished_at",
}

// NewPracticeSessionsDao creates and returns a new DAO object for table data access.
func NewPracticeSessionsDao(handlers ...gdb.ModelHandler) *PracticeSessionsDao {
	return &PracticeSessionsDao{
		group:    "default",
		table:    "practice_sessions",
		columns:  practiceSessionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PracticeSessionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PracticeSessionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PracticeSessionsDao) Columns() PracticeSessionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PracticeSessionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PracticeSessionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PracticeSessionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
