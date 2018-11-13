package interfaces

import "github.com/rs/zerolog/log"

type InvertedIndexModel interface {
	AttachWeightedWords(source string, weightedWords map[string]int) error
	SearchByString(str string) ([]SearchResultEntry, error)
	SearchByWords(words []string) ([]SearchResultEntry, error)
}

type SearchResultEntry struct {
	Source 	string
	Score 	int
}

func Check(err error){
	if err!=nil {
		log.Error().Msgf("%v",err)
	}
}
func Die(err error)  {
	if err!=nil {
		log.Fatal().Msgf("%v",err)
		panic(err)
	}
}

func CountWords(words []string) map[string]int {
	m := make(map[string]int)

	for _, word := range words {
		m[word]++
	}

	return m
}