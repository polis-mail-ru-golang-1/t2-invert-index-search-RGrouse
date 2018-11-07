package main

import (
	"bufio"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/invertedindex"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/web"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

func main() {
	searchingfolder := os.Getenv("SDIR") //"./search"
	interfaceAddr := os.Getenv("ADDR")	//"127.0.0.1:8080"

	indexFilesInFolder(searchingfolder)

	check(web.Start(interfaceAddr))
}

func check(err error) {
	if err != nil {
		panic(err)
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