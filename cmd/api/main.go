package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moonrhythm/parapet"

	"github.com/xkamail/epx5-bot/fsctx"
	"github.com/xkamail/epx5-bot/internal/cloud"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalln(err)
	}
}

func readEnv(key, defaultValue string) string {
	read := os.Getenv(key)
	if read == "" {
		return defaultValue
	}
	return read
}

func run(ctx context.Context) error {
	projectID := readEnv("FIREBASE_PROJECT_ID", "")
	port := readEnv("PORT", "8080")
	defaultBucket := readEnv("FIREBASE_BUCKET_NAME", "epx5-2bfff.appspot.com")

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID, StorageBucket: defaultBucket})

	if err != nil {
		return err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	storageClient, err := app.Storage(ctx)
	if err != nil {
		return err
	}
	srv := parapet.NewBackend()
	srv.Handler = handlers()
	// inject client context
	srv.Use(parapet.MiddlewareFunc(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = fsctx.NewContext(ctx, client)
			ctx = cloud.NewContext(ctx, storageClient)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}))
	srv.Use(parapet.MiddlewareFunc(middleware.StripSlashes))
	srv.Use(parapet.MiddlewareFunc(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					middleware.PrintPrettyStack(rvr)
					// TODO: sent json response
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}))
	srv.Use(parapet.MiddlewareFunc(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if len(w.Header().Get("Content-Type")) == 0 {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}))
	srv.Addr = net.JoinHostPort("", port)
	log.Println("ListenAndServe on ", srv.Addr)
	return srv.ListenAndServe()
}

func handlers() http.Handler {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello world"))
		})
		r.Route("/file", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				var p cloud.ListParam
				p.Prefix = r.URL.Query().Get("prefix")
				result, err := cloud.List(r.Context(), p)
				if err != nil {
					JSONErr(w, err)
					return
				}
				JSON(w, result)
			})
			r.Post("/folder", func(w http.ResponseWriter, r *http.Request) {
				var p struct {
					Name string `json:"name"`
				}
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					JSONErr(w, err)
					return
				}
				if err := cloud.CreateFolder(r.Context(), p.Name); err != nil {
					JSONErr(w, err)
					return
				}
				JSON(w, nil)
			})
			r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				_, fh, err := r.FormFile("file")
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if _, err := cloud.Upload(ctx, "test test", fh); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println("cloud upload error: ", err)
					return
				}
			})
		})

	})
	return r
}

func JSONErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	JSON(w, map[string]any{
		"message": err.Error(),
	})
}

func JSON(w http.ResponseWriter, i any) {
	json.NewEncoder(w).Encode(i)
}
