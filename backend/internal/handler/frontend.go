package handler

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:static
var staticFiles embed.FS

// RegisterFrontend serves the embedded frontend files
func RegisterFrontend(r *gin.Engine) {
	// Get the sub-filesystem for the static directory
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	// Serve static files
	fileServer := http.FileServer(http.FS(staticFS))

	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path

		// If it's an API request, let it fall through to 404
		if strings.HasPrefix(p, "/api") {
			return
		}

		// Clean the path to avoid directory traversal or weirdness
		cleanPath := strings.TrimPrefix(p, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		// Check if the file exists and is not a directory
		f, err := staticFS.Open(cleanPath)
		if err == nil {
			stat, err := f.Stat()
			_ = f.Close()
			if err == nil && !stat.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		// Otherwise, serve index.html for SPA routing
		// We read it manually to avoid redirects and ensure it's served for all SPA routes
		indexContent, err := fs.ReadFile(staticFS, "index.html")
		if err != nil {
			c.String(http.StatusNotFound, "index.html not found")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", indexContent)
	})
}
