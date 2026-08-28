// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WordsDao is the data access object for the table words.
type WordsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  WordsColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// WordsColumns defines and stores column names for the table words.
type WordsColumns struct {
	Id        string //
	Level     string //
	Word      string //
	Meaning   string //
	Phonetic  string //
	Example   string //
	CreatedAt string //
}

// wordsColumns holds the columns for the table words.
var wordsColumns = WordsColumns{
	Id:        "id",
	Level:     "level",
	Word:      "word",
	Meaning:   "meaning",
	Phonetic:  "phonetic",
	Example:   "example",
	CreatedAt: "created_at",
}

// NewWordsDao creates and returns a new DAO object for table data access.
func NewWordsDao(handlers ...gdb.ModelHandler) *WordsDao {
	return &WordsDao{
		group:    "default",
		table:    "words",
		columns:  wordsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *WordsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *WordsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *WordsDao) Columns() WordsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *WordsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *WordsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *WordsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
