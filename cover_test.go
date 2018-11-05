package main

import (
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/invertedindex"
	"os"
	"reflect"
	"testing"
	"time"
)

type indexEntry = invertedindex.IndexEntry
var wordsEntry1, wordsEntry2, wordsEntry3 invertedindex.WordsEntry
func init() {
	wordsEntry1 = invertedindex.WordsEntry{
		"one",
		map[string]int{
			"1": 1,
			"2": 2,
			"3": 3,
		},
	}
	wordsEntry2 = invertedindex.WordsEntry{
		"two",
		map[string]int{
			"1": 1,
			"2": 2,
			"4": 4,
			"3": 3,
		},
	}
	wordsEntry3 = invertedindex.WordsEntry{}
}

func TestNewIndex(t *testing.T){
	m1 := invertedindex.GetIndexMap()
	m1["1"]=nil
	invertedindex.NewIndexMap()
	m2 := invertedindex.GetIndexMap()

	if reflect.DeepEqual(m1, m2) {
		t.Errorf("%v is equal to unexpected %v", m1, m2)
	}
}

func TestCountWords(t *testing.T) {
	s := []string{"1", "3", "2", "2", "3", "3"}
	m := invertedindex.CountWords(s)
	actm := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	if !reflect.DeepEqual(*m, actm) {
		t.Errorf("%v is not equal to expected %v", m, actm)
	}

	emptym := invertedindex.CountWords([]string{})
	if len(*emptym) != 0 {
		t.Errorf("%v is not empty (0 len arg)", emptym)
	}

	emptym = invertedindex.CountWords(nil)
	if len(*emptym) != 0 {
		t.Errorf("%v is not empty (nil arg)", emptym)
	}
}

func TestAttachCountedWordsFromChannel(t *testing.T) {
	invertedindex.NewIndexMap()

	resultMap1 := map[string][]indexEntry{
		"1": []indexEntry{indexEntry{"one", 1}},
		"2": []indexEntry{indexEntry{"one", 2}},
		"3": []indexEntry{indexEntry{"one", 3}},
	}
	resultMap2 := map[string][]indexEntry{
		"1": []indexEntry{indexEntry{"one", 1}, indexEntry{"two", 1}},
		"2": []indexEntry{indexEntry{"one", 2}, indexEntry{"two", 2}},
		"3": []indexEntry{indexEntry{"one", 3}, indexEntry{"two", 3}},
		"4": []indexEntry{indexEntry{"two", 4}},
	}

	ch := make(chan invertedindex.WordsEntry)

	go func() {
		invertedindex.AttachCountedWordsFromChannel(ch, 3)
	}()

	ch <- wordsEntry1
	time.Sleep(100000)
	if !reflect.DeepEqual(invertedindex.GetIndexMap(), resultMap1) {
		t.Errorf("%v is not equal to expected %v (attach first map)", invertedindex.GetIndexMap(), resultMap1)
	}

	ch <- wordsEntry2
	time.Sleep(100000)
	if !reflect.DeepEqual(invertedindex.GetIndexMap(), resultMap2) {
		t.Errorf("%v is not equal to expected %v (attach second map)", invertedindex.GetIndexMap(), resultMap2)
	}

	ch <- wordsEntry3
	time.Sleep(100000)
	if !reflect.DeepEqual(invertedindex.GetIndexMap(), resultMap2) {
		t.Errorf("%v is not equal to expected %v (attach empty map)", invertedindex.GetIndexMap(), resultMap2)
	}
}
func TestSearchByWords(t *testing.T){
	invertedindex.NewIndexMap()

	ch := make(chan invertedindex.WordsEntry)
	go func() {
		invertedindex.AttachCountedWordsFromChannel(ch, 2)
	}()
	ch <- wordsEntry1
	ch <- wordsEntry2
	time.Sleep(100000)

	m := invertedindex.SearchByWords([]string{})
	if len(m)!=0{
		t.Errorf("%v is not empty (0 len arg)", m)
	}

	m = invertedindex.SearchByWords(nil)
	if len(m)!=0{
		t.Errorf("%v is not empty (nil arg)", m)
	}

	m = invertedindex.SearchByWords([]string{"1"})
	resultSearch1 := map[string]int{
		"one": 1,
		"two": 1,
	}
	if !reflect.DeepEqual(m, resultSearch1) {
		t.Errorf("%v is not equal to expected %v", m, resultSearch1)
	}

	m = invertedindex.SearchByWords([]string{"1", "4"})
	resultSearch2 := map[string]int{
		"one": 1,
		"two": 5,
	}
	if !reflect.DeepEqual(m, resultSearch2) {
		t.Errorf("%v is not equal to expected %v", m, resultSearch2)
	}
}

func TestSearchByString(t *testing.T){
	invertedindex.NewIndexMap()

	ch := make(chan invertedindex.WordsEntry)
	go func() {
		invertedindex.AttachCountedWordsFromChannel(ch, 2)
	}()
	ch <- wordsEntry1
	ch <- wordsEntry2
	time.Sleep(100000)

	m1 := invertedindex.SearchByString("")
	m2 := invertedindex.SearchByWords([]string{})
	if !reflect.DeepEqual(m1, m2){
		t.Errorf("%v is not equal to expected %v (0 len arg)", m1, m2)
	}

	m1 = invertedindex.SearchByString("1")
	m2 = invertedindex.SearchByWords([]string{"1"})
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("%v is not equal to expected %v", m1, m2)
	}

	m1 = invertedindex.SearchByString("1 4")
	m2 = invertedindex.SearchByWords([]string{"1", "4"})
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("%v is not equal to expected %v", m1, m2)
	}
}

func BenchmarkIndexing(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	searchingfolder := gopath+"/src/"+"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/search/"
	for i := 0; i < b.N; i++ {
		invertedindex.NewIndexMap()
		indexFilesInFolder(searchingfolder)
	}
}