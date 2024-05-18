package models

import (
	"context"

	"google.golang.org/api/iterator"

	"github.com/Anjasfedo/go-react-fireauth/configs"
)

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *Post) GetAll(ctx context.Context) ([]Post, error) {
	var posts []Post

	iter := configs.FirestoreClient.Collection("posts").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		data := doc.Data()
		title, ok1 := data["title"].(string)
		content, ok2 := data["content"].(string)
		if !ok1 || !ok2 {
			continue // Skip this item if type assertion fails
		}

		post := Post{
			Title:   title,
			Content: content,
		}
		posts = append(posts, post)
	}

	return posts, nil
}
