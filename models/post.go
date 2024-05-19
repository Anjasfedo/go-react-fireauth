package models

import (
	"context"
	"errors"
	"log"
	"mime/multipart"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Anjasfedo/go-react-fireauth/configs"
	"github.com/Anjasfedo/go-react-fireauth/storages"
)

var ErrorDocumentNotFound = errors.New("document not found")

type PostResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"imageURL"`
}

type PostRequest struct {
	Title    string `json:"title" form:"title" firestore:"title" binding:"required,min=5"`
	Content  string `json:"content" form:"content" firestore:"content"`
	ImageURL string `json:"imageURL" form:"imageURL" firestore:"imageURL"`
}

func (h *PostResponse) GetAll(ctx context.Context) ([]PostResponse, error) {
	var posts []PostResponse

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
		imageURL, ok3 := data["imageURL"].(string)
		if !ok1 || !ok2 || !ok3 {
			continue
		}

		post := PostResponse{
			ID:       doc.Ref.ID,
			Title:    title,
			Content:  content,
			ImageURL: imageURL,
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (h *PostResponse) GetByID(ctx context.Context, ID string) (*PostResponse, error) {

	dsnap, err := configs.FirestoreClient.Collection("posts").Doc(ID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrorDocumentNotFound
		}

		log.Printf("Error retrieving document with ID %s: %v", ID, err)
		return nil, err
	}

	post := &PostResponse{
		ID: ID,
	}

	if err := dsnap.DataTo(&post); err != nil {
		return nil, err
	}

	return post, nil
}

func (h *PostResponse) Add(ctx context.Context, data PostRequest, file multipart.File) (*string, error) {
	var imageURL string
	var err error
	
	if file != nil {
		imageURL, err = storages.UploadFile(ctx, file)
		if err != nil {
			log.Printf("Error uploading file: %s\n", err)
			return nil, err
		}
	}

	data.ImageURL = imageURL

	ref := configs.FirestoreClient.Collection("posts").NewDoc()

	_, err = ref.Set(ctx, data)
	if err != nil {
		log.Printf("An error has occurred: %s\n", err)
		return nil, err
	}

	return &ref.ID, err
}

func (h *PostResponse) UpdateByID(ctx context.Context, ID string, data PostRequest, file multipart.File) (*PostResponse, error) {
	_, err := h.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	updates := []firestore.Update{
		{Path: "title", Value: data.Title},
		{Path: "content", Value: data.Content},
		{Path: "imageURL", Value: data.ImageURL},
	}

	_, err = configs.FirestoreClient.Collection("posts").Doc(ID).Update(ctx, updates)
	if err != nil {
		log.Printf("An error has occured: %s\n", err)
		return nil, err
	}

	updatedPost := &PostResponse{
		ID:       ID,
		Title:    data.Title,
		Content:  data.Content,
		ImageURL: data.ImageURL,
	}

	return updatedPost, nil
}

func (h *PostResponse) DeleteById(ctx context.Context, ID string) error {
	_, err := h.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	_, err = configs.FirestoreClient.Collection("posts").Doc(ID).Delete(ctx)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		return err
	}

	return nil
}
