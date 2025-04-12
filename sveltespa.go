package sveltespa

import (
	"embed"
	"errors"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func EmbeddedRouter(embeddedFrontend embed.FS, root, notFoundFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Root path
		if path == "/" {
			fullPath := filepath.Join(root, "index.html")
			file, err := embeddedFrontend.ReadFile(fullPath)
			if !errors.Is(err, fs.ErrNotExist) && err != nil {
				http.NotFound(w, r)
				return
			}
			if errors.Is(err, fs.ErrNotExist) {
				fullPath = filepath.Join(root, notFoundFile)
				file, err = embeddedFrontend.ReadFile(fullPath)
				if err != nil {
					http.NotFound(w, r)
					return
				}
			}
			w.Header().Set("Content-Type", http.DetectContentType(file))
			w.Write(file)
			return
		}

		fileExtension := filepath.Ext(path)
		filePath := filepath.Join(root, path)
		file, err := embeddedFrontend.ReadFile(filePath)
		if err == nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(fileExtension))
			w.Write(file)
			return
		}

		// Try with ".html"
		filePath = filepath.Join(root, path+".html")
		file, err = embeddedFrontend.ReadFile(filePath)
		if err == nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(fileExtension))
			w.Write(file)
			return
		}

		// Fallback to notFoundFile
		filePath = filepath.Join(root, notFoundFile)
		file, err = embeddedFrontend.ReadFile(filePath)
		if err == nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(fileExtension))
			w.Write(file)
			return
		}
		http.NotFound(w, r)
	}
}

func Router(root, notFoundFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clean the request path
		path := r.URL.Path

		// Check if the request is for the root "/"
		if path == "/" {
			http.ServeFile(w, r, filepath.Join(root, "index.html"))
			return
		}

		// Redirect if path ends with a slash
		if len(path) > 1 && strings.HasSuffix(path, "/") {
			http.Redirect(w, r, strings.TrimSuffix(path, "/"), http.StatusMovedPermanently)
			return
		}

		// Try to serve the requested file as-is (for images and other files with extensions)
		filePath := filepath.Join(root, path)

		if fileInfo, err := os.Stat(filePath); err == nil && !fileInfo.IsDir() {
			http.ServeFile(w, r, filePath)
			return
		}

		// If not found, try to serve the requested file with ".html" extension  this is due to svelte prerender stuff
		filePath = filepath.Join(root, path+".html")
		if fileInfo, err := os.Stat(filePath); err == nil && !fileInfo.IsDir() {
			http.ServeFile(w, r, filePath)
			return
		}

		// If still not found, serve "not-found.html"
		filePath = filepath.Join(root, notFoundFile)
		if fileInfo, err := os.Stat(filePath); err == nil && !fileInfo.IsDir() {
			http.ServeFile(w, r, filePath)
			return
		}
		http.NotFound(w, r)
	}
}
