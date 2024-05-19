package models

import (
	"context"
	"errors"
	"log"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Anjasfedo/go-react-fireauth/configs"
)

var ErrorDocumentNotFound = errors.New("document not found")

type Post struct {
	ID      string `json:"id"`
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
			ID:      doc.Ref.ID,
			Title:   title,
			Content: content,
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (h *Post) GetByID(ctx context.Context, ID string) (*Post, error) {

	dsnap, err := configs.FirestoreClient.Collection("posts").Doc(ID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrorDocumentNotFound
		}

		log.Printf("Error retrieving document with ID %s: %v", ID, err)
		return nil, err
	}
	var post Post

	if err := dsnap.DataTo(&post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (h *Post) Add(ctx context.Context, data Post) error {
	_, _, err := configs.FirestoreClient.Collection("posts").Add(ctx, data)
	if err != nil {
		log.Printf("An Error has occurred: %s\n", err)
	}

	return nil
}
