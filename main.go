package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rakyll/statik/fs"

	"github.com/tobiaskohlbau/dustrobo/api"
	_ "github.com/tobiaskohlbau/dustrobo/statik"
)

func main() {
	ui := flag.Bool("ui", false, "enable ui")
	flag.Parse()

	r := chi.NewRouter()
	r.Mount("/api", api.NewHandler())

	if *ui {
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}

		r.Mount("/", http.RedirectHandler("/ui/", http.StatusMovedPermanently))
		FileServer(r, "/ui", statikFS)
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fileSystem := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".js") || strings.Contains(r.URL.Path, ".css") || strings.Contains(r.URL.Path, ".html") {
			fileSystem.ServeHTTP(w, r)
			return
		}

		data, err := fs.ReadFile(root, "/index.html")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Write(data)
	}))
}
