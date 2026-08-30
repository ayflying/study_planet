package v1

import "github.com/gogf/gf/v2/frame/g"

// Task 每日任务。
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	DueDate     string `json:"due_date"`
	Points      int    `json:"points"`
	Status      string `json:"status"` // pending | done | overdue(计算)
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at"`
}

// ListTasksReq 学生任务列表（按 student_id，公开接口）。
type ListTasksReq struct {
	g.Meta    `path:"/tasks" method:"get" tags:"Task" summary:"任务列表"`
	StudentID int    `json:"student_id" in:"query"`
	Status    string `json:"status" in:"query"`
}
type ListTasksRes []Task

// AddTaskReq 发布任务（家长鉴权）。
type AddTaskReq struct {
	g.Meta    `path:"/parent/tasks" method:"post" tags:"Task" summary:"发布任务"`
	Title     string `json:"title" v:"required"`
	Type      string `json:"type"`
	DueDate   string `json:"due_date"`
	Points    int    `json:"points"`
	StudentID int    `json:"student_id"`
}
type AddTaskRes struct {
	OK bool `json:"ok"`
}

// DeleteTaskReq 删除任务（家长鉴权）。
type DeleteTaskReq struct {
	g.Meta `path:"/parent/tasks/:id" method:"delete" tags:"Task" summary:"删除任务"`
	ID     int `in:"path" json:"-"`
}
type DeleteTaskRes struct {
	OK bool `json:"ok"`
}

// CompleteTaskReq 完成任务（学生操作，公开接口）。
type CompleteTaskReq struct {
	g.Meta    `path:"/tasks/:id/complete" method:"post" tags:"Task" summary:"完成任务"`
	ID        int `in:"path" json:"-"`
	StudentID int `json:"student_id" in:"query"`
}
type CompleteTaskRes struct {
	OK      bool `json:"ok"`
	Already bool `json:"already,omitempty"`
}

// StudentAddTaskReq 学生自建任务（公开接口，无需家长鉴权）。
type StudentAddTaskReq struct {
	g.Meta    `path:"/tasks" method:"post" tags:"Task" summary:"学生自建任务"`
	Title     string `json:"title" v:"required"`
	Type      string `json:"type"`
	DueDate   string `json:"due_date"`
	Points    int    `json:"points"`
	StudentID int    `json:"student_id"`
}
type StudentAddTaskRes struct {
	OK bool `json:"ok"`
}

// StudentDeleteTaskReq 学生删除自建任务（公开接口）。
type StudentDeleteTaskReq struct {
	g.Meta    `path:"/tasks/:id" method:"delete" tags:"Task" summary:"学生删除任务"`
	ID        int `in:"path" json:"-"`
	StudentID int `json:"student_id" in:"query"`
}
type StudentDeleteTaskRes struct {
	OK bool `json:"ok"`
}
