package types

/*
Package types is used to store the context struct which
is passed while templates are executed.
*/
//Task is the struct used to identify tasks
type Task struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Created      string    `json:"created"`
	Priority     string    `json:"priority"`
	Category     string    `json:"category"`
	Referer      string    `json:"referer,omitempty"`
	Comments     []Comment `json:"comments,omitempty"`
	IsOverdue    bool      `json:"isoverdue, omitempty"`
	IsHidden     int       `json:"ishidden, omitempty`
	CompletedMsg string    `json:"completedmsg, omitempty"`
	ShowComment  bool      `json:"showComment"`
}

type Tasks []Task

//Comment is the struct used to populate comments per tasks
type Comment struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Created  string `json:"created"`
	Username string `json:"author"`
}

//Context is the struct passed to templates
type Context struct {
	Tasks      []Task
	Navigation string
	Search     string
	Message    string
	CSRFToken  string
	Categories []Category
	Referer    string
}

//CategoryCount is the struct used to populate the sidebar
//which contains the category name and the count of the tasks
//in each category
type CategoryCount struct {
	Name  string
	Count int
}

//Status is the JSON struct to be returned
type Status struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

//Category is the structure of the category table
type Category struct {
	ID        int    `json:"categoryId"`
	Name      string `json:"categoryName"`
	TaskCount int    `json:"taskCount"`
}

//Categories will show
type Categories []Category
