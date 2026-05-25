package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

// UploadToGCS mengunggah file gambar ke Google Cloud Storage
// Mengembalikan URL publik dari gambar yang berhasil diunggah.
func UploadToGCS(file *multipart.FileHeader) (string, error) {
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return "", fmt.Errorf("GCS_BUCKET_NAME environment variable not set")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("file.Open: %v", err)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	// Membuat format unik dengan UNIX timestamp
	filename := fmt.Sprintf("products/%d%s", time.Now().UnixNano(), ext)

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(filename)
	wc := obj.NewWriter(ctx)

	if _, err = io.Copy(wc, src); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	// Format URL Publik (pastikan bucket diset Public di GCP Console)
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
	return publicURL, nil
}
