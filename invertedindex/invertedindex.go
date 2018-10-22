package invertedindex

import (
	"strings"
)

type entry struct {
	source string
	occurrence  int
}

//карта вида: слово - массив из записей {источник, сколько раз слово употреблено в источнике}
var globalindexmap map[string][]entry

func init() {
	globalindexmap = make(map[string][]entry)
}

func AttachWordsOccurencesToGlobalMap(source string, strcmap map[string]int){
	for str, c := range strcmap {
		entries, present := globalindexmap[str]
		if !present {
			entries = make([]entry, 0)
		}
		entries = append(entries, entry{source, c})
		globalindexmap[str] = entries
	}
}

func AttachWordsListToGlobalMap(source string, words []string){
	m := make(map[string]int)

	for _, word := range words {
		m[word]++
	}

	AttachWordsOccurencesToGlobalMap(source, m)
}

func SearchByString(str string) map[string]int {
	words := strings.Split(str, " ")
	return SearchByWords(words)
}

func SearchByWords(words []string) map[string]int {
	resultmap := make(map[string]int)

	for _, word := range words {
		entries, present := globalindexmap[word]
		if !present {
			continue
		}

		for _, entry := range entries {
			resultmap[entry.source]+=entry.occurrence
		}
	}
	return resultmap
}