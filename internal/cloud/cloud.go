package cloud

import (
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"path"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrUploadFailed    = errors.New("upload failed")
)

func Upload(ctx context.Context, pathname string, fh multipart.FileHeader) (any, error) {
	bucket, err := client(ctx).DefaultBucket()
	if err != nil {
		return nil, err
	}
	fileName := path.Join(pathname, fh.Filename)

	fileType := fh.Header["Content-Type"][0]

	if !(fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/gif" || fileType == "image/webp" || fileType == "application/pdf") {
		return nil, ErrInvalidFileType
	}

	file, err := fh.Open()
	if err != nil {
		return nil, err
	}

	// Upload an object with storage.Writer.
	wc := bucket.Object(fileName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		log.Printf("storage io.Copy error; %v", err)
		return nil, ErrUploadFailed
	}
	if err := wc.Close(); err != nil {
		log.Printf("storage wc.Copy error; %v", err)
		return nil, ErrUploadFailed
	}

	return nil, nil
}
