package server

import (
	"embed"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// frontend contains the react app other assets etc..
//
//go:embed dist
var frontend embed.FS

func serveReactIndex() http.Handler {
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "index.html", stat.ModTime(), rdr)
	})

}

func makeStaticServer() http.Handler {
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
	return http.FileServer(http.FS(staticFiles))
}
