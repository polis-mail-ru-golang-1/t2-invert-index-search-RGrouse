package web

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"time"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/invertedindex"
)

type SearchResult struct {
	Query string
	Result template.HTML
}

var searchTmpl, resultTmpl *template.Template

func init(){
	searchTmpl = template.Must(template.ParseFiles("./templates/search.html"))
	resultTmpl = template.Must(template.ParseFiles("./templates/result.html"))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query != "" {
		q := strings.ToLower(query)

		resultmap := invertedindex.SearchByString(q)

		buf := new(bytes.Buffer)
		invertedindex.SortAndPrintResult(resultmap, buf)

		result := strings.Replace(buf.String(), "\n", "<br>", -1)

		resultTmpl.Execute(w, SearchResult{ Query:q, Result:template.HTML(result)})
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
