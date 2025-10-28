package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Embed the built frontend assets from the web directory.
//
//go:embed web/*
var webFS embed.FS

// registerFrontendRoutes serves the embedded frontend at root and
// falls back to index.html for SPA routes while keeping /api routes intact.
func registerFrontendRoutes(r *gin.Engine) {
	static, err := fs.Sub(webFS, "web")
	if err != nil {
		// If the sub FS cannot be created (e.g., dir missing), skip registration.
		return
	}

	// Explicit root handler to avoid 301 redirects on "/" by serving content directly
	r.GET("/", func(c *gin.Context) {
		b, err := fs.ReadFile(static, "index.html")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})
	r.HEAD("/", func(c *gin.Context) {
		if _, err := fs.ReadFile(static, "index.html"); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusOK)
	})

	// SPA + static handling via NoRoute to avoid wildcard conflicts with /api
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}

		// Try to serve the requested static file first
		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		if f, openErr := static.Open(filePath); openErr == nil {
			_ = f.Close()
			c.FileFromFS(filePath, http.FS(static))
			return
		}

		// Fallback to SPA entry
		c.FileFromFS("index.html", http.FS(static))
	})
}
