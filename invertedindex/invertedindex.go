package invertedindex

import (
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

func GetIndexMap() map[string][]IndexEntry {
	return gindexmap
}

func NewIndexMap() {
	gindexmap = make(map[string][]IndexEntry)
}

func AttachCountedWordsFromChannel(ch chan WordsEntry, k int){
	for i := 0; i<k; i++ {
		wordsEntry := <- ch
		attachWordsOccurencesToGlobalMap(wordsEntry.Source, wordsEntry.CountedWords)
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
	return resultmap
}