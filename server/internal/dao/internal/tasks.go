// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TasksDao is the data access object for the table tasks.
type TasksDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of the current DAO.
	columns TasksColumns // columns contains all the column names of Table for convenient usage.
}

// TasksColumns defines and stores column names for the table tasks.
type TasksColumns struct {
	Id          string //
	Title       string //
	Type        string //
	DueDate     string //
	Points      string //
	Status      string //
	ChildId     string //
	CreatedAt   string //
	CompletedAt string //
}

// tasksColumns holds the columns for the table tasks.
var tasksColumns = TasksColumns{
	Id:          "id",
	Title:       "title",
	Type:        "type",
	DueDate:     "due_date",
	Points:      "points",
	Status:      "status",
	ChildId:     "child_id",
	CreatedAt:   "created_at",
	CompletedAt: "completed_at",
}

// NewTasksDao creates and returns a new DAO object for table data access.
func NewTasksDao() *TasksDao {
	return &TasksDao{
		group:   "default",
		table:   "tasks",
		columns: tasksColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TasksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TasksDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TasksDao) Columns() TasksColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TasksDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TasksDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *TasksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
