package models

type Post struct {
	ID           int
	AuthorID     int
	LikeCount    int
	DislikeCount int
	CommentCount int
	Vote         int
	Author       string
	Title        string
	Content      string
	ImgPath      string
	Categories   []string
}
