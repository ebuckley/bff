package server

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

// frontend contains the react app other assets etc..
//
//go:embed dist
var frontend embed.FS

func serveReactIndex(prefix string) http.Handler {
	var viteProxy http.Handler
	if development == "true" {
		viteDevServerURL, _ := url.Parse("http://localhost:5173") // Default Vite dev server address
		viteProxy = httputil.NewSingleHostReverseProxy(viteDevServerURL)
		return viteProxy
	}
	fp, err := frontend.Open("dist/index.html")
	if err != nil {
		slog.Error("failed to open index.html", "err", err)
		panic(err)
	}
	rdr, ok := fp.(io.ReadSeeker)
	if !ok {
		slog.Error("failed to open index.html and coerce it in to a readseeker")
		panic("failed to open index.html and coerce it in to a readseeker")
	}
	stat, err := fp.Stat()
	if err != nil {
		slog.Error("failed to stat index.html", "err", err)
		panic(err)
	}
	if prefix != "" {
		slog.Debug("rewriting index.html with prefix", "prefix", prefix)
		srcContent, err := io.ReadAll(rdr)
		if err != nil {
			slog.Error("failed to read index.html", "err", err)
			panic(err)
		}
		// rewrite the index.html with the prefix
		srcRegex := regexp.MustCompile(`(src|href)="(/[^"]*)"`)
		// Replace the asset sources with prefixed versions
		modifiedContent := srcRegex.ReplaceAllStringFunc(string(srcContent), func(match string) string {
			parts := srcRegex.FindStringSubmatch(match)
			if len(parts) == 3 {
				return fmt.Sprintf(`%s="%s%s"`, parts[1], prefix, parts[2])
			}
			return match
		})
		rdr = io.ReadSeeker(strings.NewReader(modifiedContent))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// insert the prefix into the index.html
		http.ServeContent(w, r, "index.html", stat.ModTime(), rdr)
	})

}

func (s *Server) makeStaticServer() http.Handler {
	var viteProxy http.Handler
	if development == "true" {
		viteDevServerURL, _ := url.Parse("http://localhost:5173") // Default Vite dev server address
		viteProxy = httputil.NewSingleHostReverseProxy(viteDevServerURL)
		return viteProxy
	}
	staticFiles, err := fs.Sub(frontend, "dist")
	if err != nil {
		panic(err)
	}

	if s.handlerPrefix != "" {
		slog.Debug("serving static files with prefix", "prefix", s.handlerPrefix)
		// there is a prefix, we want to look up
		return http.StripPrefix(s.handlerPrefix, http.FileServer(http.FS(staticFiles)))
	} else {
		return http.FileServer(http.FS(staticFiles))
	}
}
