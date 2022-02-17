package ginx

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

const INDEX = "index.html"

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		fmt.Printf("Exists prefix: %v filepath: %v true \n", prefix, filepath)

		name := path.Join(l.root, p)
		fmt.Printf("Exists name: %v\n", name)

		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}

	fmt.Printf("Exists prefix: %v filepath: %v false \n", prefix, filepath)
	return false
}

// ServeRoot returns a middleware handler that serves static files in the given directory.
func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

// Serve returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	if urlPrefix != "" {
		fileServer = http.StripPrefix(urlPrefix, fileServer)
	}
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, urlPrefix) {
			fmt.Printf("HAS PREFIX urlPrefix: %v c.Request.URL.Path: %v \n", urlPrefix, c.Request.URL.Path)

			if fs.Exists(urlPrefix, c.Request.URL.Path) {
				fileServer.ServeHTTP(c.Writer, c.Request)
				c.Abort()
			}
		}
	}
}

type staticFileSystem struct {
	fileSystem     fs.FS
	httpFileSystem http.FileSystem
	root           string
	indexes        bool
}

func StaticFS(fileSystem fs.FS, root string, indexes bool) *staticFileSystem {
	httpFS := http.FS(fileSystem)

	return &staticFileSystem{
		fileSystem:     fileSystem,
		root:           root,
		indexes:        indexes,
		httpFileSystem: httpFS,
	}
}

func (l *staticFileSystem) Open(name string) (http.File, error) {
	fmt.Printf("Open name: %v root: %v\n", name, l.root)
	return l.httpFileSystem.Open(name)
	//
	//if p := strings.TrimPrefix(name, l.root); len(p) < len(name) {
	//
	//    fmt.Printf("Open name: %v\n", p)
	//
	//    return l.httpFileSystem.Open(p)
	//}
	//
	//return nil, errors.Errorf("unable to find %v", name)
}

func (l *staticFileSystem) Exists(prefix string, filepath string) bool {
	fmt.Printf("Exists prefix: %v filepath: %v \n", prefix, filepath)
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		fmt.Printf("Exists prefix: %v filepath: %v true \n", prefix, filepath)

		fmt.Printf("Exists name: %v\n", p)

		if p == "" || p == "/" {
			p = "index.html"

		}
		fmt.Printf("Exists name: %v\n", p)

		file, _ := l.fileSystem.Open(p)
		return file != nil
	}

	fmt.Printf("Exists prefix: %v filepath: %v false \n", prefix, filepath)
	return false
}
