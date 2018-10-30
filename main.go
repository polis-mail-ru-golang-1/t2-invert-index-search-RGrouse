package main

import (
	"bufio"
	"fmt"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/invertedindex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query != "" {
		q := strings.ToLower(query)
		resultmap := invertedindex.SearchByString(q)
		if len(resultmap)>0 {
			sortAndPrintResultMap(resultmap, w)
		}
	}
}

func main() {
	http.HandleFunc("/search", handler)

	searchingfolder := os.Args[1] 	//"./search"
	interfaceAddr := os.Args[2]		//"127.0.0.1:8080"

	indexFilesInFolder(searchingfolder)

	check(http.ListenAndServe(interfaceAddr, nil))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func sortAndPrintResultMap(m map[string]int, w io.Writer) {
	n := map[int][]string{}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, k := range a {
		for _, s := range n[k] {
			fmt.Fprintf(w, "- %s; совпадений - %d\n", s, k)
		}
	}
}

func indexFilesInFolder(searchingfolder string) {
	filesInfos, err := ioutil.ReadDir(searchingfolder)
	check(err)

	ch := make(chan invertedindex.WordsEntry)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		invertedindex.AttachCountedWordsFromChannel(ch, len(filesInfos))	//слушаем канал и добавляем в общий индекс посчитанные слова
		wg.Done()
	}()

	for _, fileInfo := range filesInfos {
		wg.Add(1)

		go func(f os.FileInfo) {
			defer wg.Done()

			if (!f.IsDir()) {
				words, err := wordsInFile(searchingfolder + "/" + f.Name()) //разбиваем файл по словам
				check(err)
				countedWords := invertedindex.CountWords(words)      //считаем, сколько раз слово появилось в файле
				ch <- invertedindex.WordsEntry{f.Name(), *countedWords} //пишем в канал источник и карту посчитанных слов
			}
		}(fileInfo)
	}
	wg.Wait()	//ждем пока все файлы проиндексируются
}

func wordsInFile(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	words := make([]string, 0)
	for scanner.Scan() {
		words = append(words, strings.ToLower(scanner.Text()))
	}
	return words, nil
}