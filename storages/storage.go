package storages

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

    log.Println("Starting file upload to Firebase Storage")

    if _, err := io.Copy(wc, file); err != nil {
        log.Printf("Error copying file to writer: %v\n", err)
        return "", fmt.Errorf("io.Copy: %v", err)
    }

    log.Println("File uploaded successfully")

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

func DownloadFile(ctx *gin.Context, objectName string, destFilePath string) {
	rc, err := configs.StorageBucket.Object(objectName).NewReader(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Object(%q).NewReader: %v", objectName, err))
		return
	}
	defer rc.Close()

	f, err := os.Create(destFilePath)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("os.Create: %v", err))
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, rc); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("io.Copy: %v", err))
		return
	}

	// Set the appropriate headers
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+objectName)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")

	// Send the file as response
	ctx.File(destFilePath)

	log.Printf("Blob %v downloaded and returned to client.\n", objectName)
}
