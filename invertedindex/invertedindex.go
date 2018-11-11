package invertedindex

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"sort"
	"strings"
)

type IndexEntry struct {
	Source     string
	Occurrence int
}
type WordsEntry struct {
	Source       string
	CountedWords map[string]int
}

//карта вида: слово - массив из записей {источник, сколько раз слово употреблено в источнике}
var gindexmap map[string][]IndexEntry

func init() {
	NewIndexMap()
}

func NewIndexMap() {
	gindexmap = make(map[string][]IndexEntry)
}

func AttachCountedWordsFromChannel(ch chan WordsEntry, k int){
	for i := 0; i<k; i++ {
		wordsEntry := <- ch
		attachWordsOccurencesToGlobalMap(wordsEntry.Source, wordsEntry.CountedWords)
		log.Info().Msg("Добавили в индекс посчитанные слова из файла "+wordsEntry.Source)
	}
}

func attachWordsOccurencesToGlobalMap(source string, strcmap map[string]int){
	for str, c := range strcmap {
		entries, present := gindexmap[str]
		if !present {
			entries = make([]IndexEntry, 0)
		}
		entries = append(entries, IndexEntry{source, c})
		gindexmap[str] = entries
	}
}

func CountWords(words []string) *map[string]int {
	m := make(map[string]int)

	for _, word := range words {
		m[word]++
	}

	return &m
}

func SearchByString(str string) map[string]int {
	words := strings.Split(str, " ")
	return SearchByWords(words)
}

func SearchByWords(words []string) map[string]int {
	resultmap := make(map[string]int)

	for _, word := range words {
		entries, present := gindexmap[word]
		if !present {
			continue
		}

		for _, entry := range entries {
			resultmap[entry.Source]+=entry.Occurrence
		}
	}
	log.Info().Msg("Поиск по строке "+strings.Join(words, " "))
	return resultmap
}

func SortAndPrintResult(m map[string]int, w io.Writer) {
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