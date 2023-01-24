package models

import "html/template"

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
	ImagesPath   []template.URL
	Categories   []string
}
