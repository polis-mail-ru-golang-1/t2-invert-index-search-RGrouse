package web

import (
	"fmt"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/invertedindex"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type SearchResult struct {
	Query string
	Result []invertedindex.SearchResultEntry
}

type Web struct {
	Address string
	Index invertedindex.InvertedIndex
}

var searchTmpl, resultTmpl *template.Template

func init(){
	var allFiles []string
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}

	templates, err := template.ParseFiles(allFiles...)

	searchTmpl = templates.Lookup("search")
	resultTmpl = templates.Lookup("result")
}

func (web *Web) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query != "" {
		q := strings.ToLower(query)

		result := web.Index.SearchByString(q)

		resultTmpl.Execute(w, SearchResult{ Query:q, Result:result})
	} else {
		http.Redirect(w, r, "/main", 301)
	}
}

func (web *Web) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/search", web.searchHandler)
	mux.Handle("/res/", http.StripPrefix(
		"/res/",
		http.FileServer(http.Dir("./static")),
	))
	mux.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		searchTmpl.Execute(w, nil)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/main", 301)
	})

	server := http.Server{
		Addr:         web.Address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
