package models

type Comment struct {
	ID           int
	UserID       int
	PostID       int
	LikeCount    int
	DislikeCount int
	Vote         int
	Content      string
	Author       string
}
