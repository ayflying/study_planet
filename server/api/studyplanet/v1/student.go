package v1

import "github.com/gogf/gf/v2/frame/g"

// Student 学生信息（students 列表与创建/修改返回体）。
type Student struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Grade     int    `json:"grade"`
	CreatedAt string `json:"created_at"`
}

// ListStudentsReq 学生列表（学生模式可访问，公开接口）。
type ListStudentsReq struct {
	g.Meta `path:"/students" method:"get" tags:"Student" summary:"学生列表"`
}
type ListStudentsRes []Student

// CreateStudentReq 新建学生账号（家长鉴权）。
type CreateStudentReq struct {
	g.Meta   `path:"/parent/students" method:"post" tags:"Student" summary:"新建学生"`
	Name     string `json:"name" v:"required"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Grade    int    `json:"grade"`
}
type CreateStudentRes Student

// UpdateStudentReq 修改学生信息（家长鉴权），指针字段为可选。
type UpdateStudentReq struct {
	g.Meta   `path:"/parent/students/:id" method:"put" tags:"Student" summary:"修改学生"`
	ID       int     `in:"path" json:"-"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Avatar   *string `json:"avatar"`
	Grade    *int    `json:"grade"`
}
type UpdateStudentRes Student

// DeleteStudentReq 删除学生及其学习数据（家长鉴权），至少保留一个学生。
type DeleteStudentReq struct {
	g.Meta `path:"/parent/students/:id" method:"delete" tags:"Student" summary:"删除学生"`
	ID     int `in:"path" json:"-"`
}
type DeleteStudentRes struct {
	OK bool `json:"ok"`
}
