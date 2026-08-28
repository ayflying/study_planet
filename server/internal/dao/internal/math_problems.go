// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MathProblemsDao is the data access object for the table math_problems.
type MathProblemsDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of the current DAO.
	columns MathProblemsColumns // columns contains all the column names of Table for convenient usage.
}

// MathProblemsColumns defines and stores column names for the table math_problems.
type MathProblemsColumns struct {
	Id          string //
	Level       string //
	Type        string //
	Question    string //
	Options     string //
	Answer      string //
	Explanation string //
}

// mathProblemsColumns holds the columns for the table math_problems.
var mathProblemsColumns = MathProblemsColumns{
	Id:          "id",
	Level:       "level",
	Type:        "type",
	Question:    "question",
	Options:     "options",
	Answer:      "answer",
	Explanation: "explanation",
}

// NewMathProblemsDao creates and returns a new DAO object for table data access.
func NewMathProblemsDao() *MathProblemsDao {
	return &MathProblemsDao{
		group:   "default",
		table:   "math_problems",
		columns: mathProblemsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MathProblemsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MathProblemsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MathProblemsDao) Columns() MathProblemsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MathProblemsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MathProblemsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *MathProblemsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
