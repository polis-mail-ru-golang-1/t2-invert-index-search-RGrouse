package invertedindex

import (
	"strings"
	"sync"
)

type safeIndexMap struct {
	sync.RWMutex
	//карта вида: слово - массив из записей {источник, сколько раз слово употреблено в источнике}
	indexmap map[string][]entry
}

type entry struct {
	source string
	occurrence  int
}

var gSafeIndexMap *safeIndexMap

func init() {
	newSafeIndexMap()
}

func newSafeIndexMap() {
	gSafeIndexMap = &safeIndexMap {
		indexmap : make(map[string][]entry),
	}
}

func AttachWordsOccurencesToGlobalMap(source string, strcmap map[string]int){
	gSafeIndexMap.Lock()
	defer gSafeIndexMap.Unlock()

	for str, c := range strcmap {
		entries, present := gSafeIndexMap.indexmap[str]
		if !present {
			entries = make([]entry, 0)
		}
		entries = append(entries, entry{source, c})
		gSafeIndexMap.indexmap[str] = entries
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

	gSafeIndexMap.RLock()
	defer gSafeIndexMap.RUnlock()

	for _, word := range words {
		entries, present := gSafeIndexMap.indexmap[word]
		if !present {
			continue
		}

		for _, entry := range entries {
			resultmap[entry.source]+=entry.occurrence
		}
	}
	return resultmap
}