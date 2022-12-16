package models

type TemplateData struct {
	Template string
	User     User
	Post     Post
	Posts    []Post
	Error    ErrorMsg
}

type ErrorMsg struct {
	Status int
	Msg    string
}
