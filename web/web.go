package web

import (
	"bytes"
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
	Result []string
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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query != "" {
		q := strings.ToLower(query)

		resultmap := invertedindex.SearchByString(q)

		buf := new(bytes.Buffer)
		invertedindex.SortAndPrintResult(resultmap, buf)

		result := strings.Split(buf.String(), "\n")

		resultTmpl.Execute(w, SearchResult{ Query:q, Result:result})
	} else {
		http.Redirect(w, r, "/main", 301)
	}
}

func Start(address string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.Handle("/res/", http.StripPrefix(
		"/res/",
		http.FileServer(http.Dir("./static")),
	))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/icons/favicon.ico")
	})
	mux.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		searchTmpl.Execute(w, nil)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/main", 301)
	})

	server := http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}
