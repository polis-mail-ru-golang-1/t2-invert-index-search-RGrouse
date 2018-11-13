package db_model

import (
	"github.com/go-pg/pg"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/model/interfaces"
	"github.com/rs/zerolog/log"
)

type DBModel struct {
	pg *pg.DB
}

func New(pg *pg.DB) DBModel {
	cleanTables(pg)

	return DBModel{
		pg: pg,
	}
}

type Word struct {
	tableName 	struct{} 	`sql:"words"`
	Id   		int      	`sql:"id,pk"`
	Word    	string   	`sql:"word"`
}
type Source struct {
	tableName 	struct{}	`sql:"sources"`
	Id			int			`sql:"id,pk"`
	Source		string		`sql:"source"`
}
type Index struct {
	tableName 	struct{} 	`sql:"index"`
	Id 			int			`sql:"id,pk"`
	Word_id 	int			`sql:"word_id"`
	Source_id 	int 		`sql:"source_id"`
	Weight 		int 		`sql:"weight"`
}

func (m DBModel) AttachWeightedWords(src string, weightedWords map[string]int) error {
	srcEntry, err := m.putSourceEntry(Source{Source: src})
	if err != nil {
		return err
	}
	for word, weight := range weightedWords {
		word, err := m.putWordEntry(Word{Word: word})
		if err != nil {
			return err
		}
		err = m.insertIndexEntry(Index{Word_id: word.Id,
			Source_id: srcEntry.Id,
			Weight:    weight})
		if err != nil {
			return err
		}
	}
	log.Info().Msg("Добавили в индекс слова из " + src)
	return nil
}

func (m DBModel) SearchByString(str string) ([]interfaces.SearchResultEntry, error) {
	words := interfaces.WordsInString(str)
	return m.SearchByWords(words)
}

func (m DBModel) SearchByWords(words []string) ([]interfaces.SearchResultEntry, error) {
	for i, _ := range words {
		words[i] = interfaces.StemWord(words[i])
	}

	matchedWordsId := m.pg.Model().
		Column("words.id").
		Table("words").
		Where("word IN (?)", pg.In(words))

	matchedIndexes := m.pg.Model().
		With("wordsids", matchedWordsId).
		Column("index.source_id", "index.weight","index.word_id").
		Table("index", "wordsids").
		Where("index.word_id IN (wordsIDs.id)")

	qresult := []interfaces.SearchResultEntry{}
	err := m.pg.Model().
		With("indexes", matchedIndexes).
		Column("sources.source").
		ColumnExpr("SUM(indexes.weight) AS score").
		Table("indexes").
		Join("JOIN sources ON indexes.source_id = sources.id").
		Group("sources.source").
		OrderExpr("score DESC").
		Select(&qresult)

	if err!=nil{
		return qresult, err
	}
	log.Info().Msgf("Поиск по словам %v, результат %v", words, qresult)
	return qresult, nil
}

func (m DBModel) putSourceEntry(source Source) (Source, error) {
	created, err := m.pg.Model(&source).
		Where("source = ?source").
		SelectOrInsert()
	if err!=nil{
		return source, err
	}
	if(created){
		log.Info().Msg("Добавили в индекс источник "+ source.Source)
	}
	return source, nil
}

func (m DBModel) putWordEntry(word Word) (Word, error){
	created, err := m.pg.Model(&word).
		Where("word = ?word").
		SelectOrInsert(&word)
	if err!=nil{
		return word, err
	}
	if(created){
		log.Info().Msg("Добавили в индекс слово "+ word.Word)
	}
	return word, nil
}

func (m DBModel) insertIndexEntry(index Index) error {
	err := m.pg.Insert(&index)
	return err
}

func cleanTables(pg *pg.DB){
	res, err1 := pg.Exec(`DELETE FROM index;`)
	interfaces.Die(err1)
	log.Info().Msgf("Удалено %v строк из таблицы index", res.RowsAffected())

	res2, err2 := pg.Exec(`DELETE FROM words;`)
	interfaces.Die(err2)
	log.Info().Msgf("Удалено %v строк из таблицы words", res2.RowsAffected())

	res3, err3 := pg.Exec(`DELETE FROM sources;`)
	interfaces.Die(err3)
	log.Info().Msgf("Удалено %v строк из таблицы sources", res3.RowsAffected())
}

