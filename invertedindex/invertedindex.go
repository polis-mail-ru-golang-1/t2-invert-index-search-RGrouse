package invertedindex

import (
	"github.com/rs/zerolog/log"
	"sort"
	"strings"
)

type IndexEntry struct {
	Source string
	Weight int
}

type SearchResultEntry struct {
	Source 	string
	Score 	int
}

type InvertedIndex struct {
	//карта вида: слово - массив из записей {источник, вес}
	indexmap map[string][]IndexEntry
}

func New() InvertedIndex {
	return InvertedIndex{ indexmap: make(map[string][]IndexEntry,0) }
}

func (ii *InvertedIndex) AttachWeightedWords(source string, weightedWords map[string]int){
	for str, w := range weightedWords {
		entries, present := ii.indexmap[str]
		if !present {
			entries = make([]IndexEntry, 0)
		}
		entries = append(entries, IndexEntry{source, w})
		ii.indexmap[str] = entries
	}
	log.Info().Msg("Добавили в индекс слова из файла "+source)
}

func (ii *InvertedIndex) SearchByString(str string) []SearchResultEntry {
	words := strings.Split(str, " ")
	return ii.SearchByWords(words)
}

func (ii *InvertedIndex) SearchByWords(words []string) []SearchResultEntry {
	resultmap := make(map[string]int)

	for _, word := range words {
		entries, present := ii.indexmap[word]
		if !present {
			continue
		}

		for _, entry := range entries {
			resultmap[entry.Source]+=entry.Weight
		}
	}
	log.Info().Msgf("Поиск по строке %v", words)
	return sortResult(resultmap)
}

func sortResult(m map[string]int) []SearchResultEntry {
	searchResult := []SearchResultEntry{}
	for s, w := range m {
		searchResult = append(searchResult, SearchResultEntry{s, w})
	}

	sort.Slice(searchResult, func(i, j int) bool {
		return searchResult[i].Score > searchResult[j].Score
	})

	return searchResult
}