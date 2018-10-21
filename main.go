package main

import (
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/invertedindex"
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	searchingfolder := os.Args[1] 	//"./search"
	filenames := os.Args[2:]		//[]string{"ex1.txt", "ex2.txt", "ex3.txt", "ex4.txt"}

	for _, filename := range filenames {
		words, err := wordsInFile(searchingfolder+"/"+filename)
		check(err)
		invertedindex.AttachWordsListToGlobalMap(filename, words)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Поисковая фраза: ")
	for scanner.Scan() {
		str := scanner.Text()

		resultmap := invertedindex.SearchByString(str)
		if len(resultmap)>0 {
			sortAndPrintResultMap(resultmap)
		}

		fmt.Print("\nПоисковая фраза: ")
	}
	check(scanner.Err())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func sortAndPrintResultMap(m map[string]int) {
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
			fmt.Printf("- %s; совпадений - %d\n", s, k)
		}
	}
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
		words = append(words, scanner.Text())
	}

	return words, nil
}