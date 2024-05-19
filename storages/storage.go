package storages

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/Anjasfedo/go-react-fireauth/configs"
)

func UploadFile(ctx context.Context, file multipart.File) (string, error) {
	filename := uuid.New().String()

	wc := configs.StorageBucket.Object(filename).NewWriter(ctx)
	defer func() {
		if err := wc.Close(); err != nil {
			fmt.Printf("error closing writer: %v\n", err)
		}
	}()

	if _, err := io.Copy(wc, file); err != nil {
		log.Printf("Error copying file to writer: %v\n", err)
		return "", fmt.Errorf("io.Copy: %v", err)
	}

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", configs.StorageBucketName, filename)
	return imageURL, nil
}

func DeleteFile(ctx context.Context, objectName string) error {
	o := configs.StorageBucket.Object(objectName)
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", objectName, err)
	}

	log.Printf("Blob %v deleted.\n", objectName)
	return nil
}

// func DownloadFile(ctx *gin.Context, objectName string) {
// 	obj := configs.StorageBucket.Object(objectName)
// 	attrs, err := obj.Attrs(ctx)
// 	if err != nil {
// 		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Object(%q).Attrs: %v", objectName, err))
// 		return
// 	}

// 	rc, err := obj.NewReader(ctx)
// 	if err != nil {
// 		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Object(%q).NewReader: %v", objectName, err))
// 		return
// 	}
// 	defer rc.Close()

// 	// Set the appropriate headers
// 	ctx.Header("Content-Description", "File Transfer")
// 	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", objectName))
// 	ctx.Header("Content-Type", attrs.ContentType)

// 	// Serve the file directly to the client without saving it locally
// 	ctx.DataFromReader(http.StatusOK, attrs.Size, attrs.ContentType, rc, nil)

// 	log.Printf("Blob %v streamed and returned to client.\n", objectName)
// }
