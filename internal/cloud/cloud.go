package cloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrUploadFailed    = errors.New("upload failed")
)

type FileType int

const (
	FileTypeFile FileType = iota
	FileTypeFolder
)

type File struct {
	Name         string    `json:"name"`
	Type         FileType  `json:"type"`
	CreatedAt    time.Time `json:"createdAt"`
	Prefix       string    `json:"prefix"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"contentType"`
	Owner        string    `json:"owner"`
	OriginalName string    `json:"originalName"`
	// DB

}
type ListParam struct {
	Prefix string `json:"prefix"`
	Limit  int    `json:"limit"`
}

func (p *ListParam) Valid() error {
	p.Prefix = strings.ToLower(strings.TrimSpace(p.Prefix))
	if p.Prefix != "" {
		//p.Prefix = "/"
		p.Prefix += "/"
	}
	return nil
}

func List(ctx context.Context, q ListParam) ([]*File, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err := q.Valid(); err != nil {
		return nil, err
	}

	bucket, err := client(ctx).DefaultBucket()
	if err != nil {
		return nil, err
	}
	it := bucket.Objects(ctx, &storage.Query{
		Prefix: q.Prefix,
	})

	dir := make([]*File, 0)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects: %v", bucket, err)
		}
		if attrs.Name == q.Prefix {
			continue
		}

		name := strings.Replace(attrs.Name, q.Prefix, "", 1)
		fileType := FileTypeFile

		if strings.Contains(name, "/") && !strings.Contains(name, ".") {
			fileType = FileTypeFolder
		}
		// skip file in sub folder
		if strings.Count(attrs.Name, "/") > strings.Count(q.Prefix, "/")+1 && strings.HasSuffix(name, "/") {
			continue
		}
		// skip directory/file_name
		if strings.Contains(name, "/") && !strings.HasSuffix(name, "/") {
			continue
		}

		dir = append(dir, &File{
			OriginalName: attrs.Name,
			Name:         name,
			Type:         fileType,
			CreatedAt:    attrs.Created,
			Prefix:       attrs.Prefix,
			Size:         attrs.Size,
			ContentType:  attrs.ContentType,
			Owner:        attrs.Owner,
		})
	}

	sort.Slice(dir, func(i, j int) bool {
		return dir[i].Type > dir[j].Type
	})

	return dir, nil
}

func CreateFolder(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	name = strings.TrimSuffix(strings.TrimPrefix(name, "/"), "/") + "/"

	// TODO: validate name
	bucket, err := client(ctx).DefaultBucket()
	if err != nil {
		return err
	}
	log.Println("create", name)
	wc := bucket.Object(name).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
	wc.ContentType = "application/x-www-form-urlencoded;charset=UTF-8"
	src := bytes.NewBufferString("")

	if _, err = io.Copy(wc, src); err != nil {
		log.Printf("storage io.Copy error; %v", err)
		return ErrUploadFailed
	}
	if err := wc.Close(); err != nil {
		log.Printf("storage wc.Copy error; %v", err)
		return ErrUploadFailed
	}
	return nil
}

func Upload(ctx context.Context, pathname string, fh *multipart.FileHeader) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

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
	wc := bucket.Object(fileName).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
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
