package cloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

		name := parseFileName(attrs.Name)
		fileType := parseFileType(attrs.Name)

		// skip file in sub folder
		if strings.Count(attrs.Name, "/") > strings.Count(q.Prefix, "/")+1 && strings.HasSuffix(name, "/") {
			continue
		}
		// skip directory/file_name
		if strings.Contains(attrs.Name, "/") && !strings.HasSuffix(attrs.Name, "/") {
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

	sort.Slice(dir, sortByFolder(dir))

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

type SearchParam struct {
	InFolder string
	Keyword  string
}

func (p *SearchParam) Valid() error {
	p.Keyword = strings.ToLower(p.Keyword)
	p.InFolder = strings.TrimSpace(p.InFolder)
	return nil
}

func Search(ctx context.Context, p SearchParam) ([]*File, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err := p.Valid(); err != nil {
		return nil, err
	}

	bucket, err := client(ctx).DefaultBucket()
	if err != nil {
		return nil, err
	}
	it := bucket.Objects(ctx, &storage.Query{
		Prefix: p.InFolder,
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
		if !strings.Contains(attrs.Name, p.Keyword) {
			continue
		}
		name := parseFileName(attrs.Name)
		fileType := parseFileType(attrs.Name)

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
	sort.Slice(dir, sortByFolder(dir))
	return dir, nil
}
func ReadFile(ctx context.Context, path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	bucket, err := client(ctx).DefaultBucket()
	if err != nil {
		return nil, err
	}
	reader, err := bucket.Object(path).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	slurp, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("readFile: unable to read data %v", err)
	}
	return slurp, nil
}

func sortByFolder(dir []*File) func(i, j int) bool {
	return func(i, j int) bool {
		return dir[i].Type > dir[j].Type
	}
}

func parseFileName(name string) string {
	name = strings.TrimSuffix(name, "/")
	split := strings.Split(name, "/")
	last := split[len(split)-1]
	if strings.Contains(name, ".") {
		return last
	}
	return last + "/"
}

func parseFileType(name string) FileType {
	if strings.Contains(name, ".") {
		return FileTypeFile
	}
	return FileTypeFolder
}
