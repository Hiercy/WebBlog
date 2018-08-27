package models

import (
	"fmt"
)

// Post comment
type Post struct {
	ID      string
	Title   string
	Content string
}

// NewPost comment
func NewPost(id, title, content string) *Post {
	if title != "" && content != "" {
		return &Post{id, title, content}
	} else {
		fmt.Println("error")
	}
	return nil
}
