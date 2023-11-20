package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type NotesDb struct {
	Notes []*Note `json:"notes"`
}

type Note struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	DoneTs  int64  `json:"done_ts"`
}

func main() {
	if len(os.Args) < 2 {
		log.Println("You need at least 2 arguments")
		return
	}

	notesDb := loadNotes()
	dirty := false
	switch os.Args[1] {
	case "list":
		notesDb.print()
	case "done":
		notesDb.complete(os.Args[2])
		dirty = true
	case "add":
		notesDb.add(os.Args[2:])
		dirty = true
	}
	if dirty {
		saveNotes(notesDb)
	}
}

func loadNotes() *NotesDb {
	notesDbPath := resolveDbPath()
	notesDb := NotesDb{make([]*Note, 0)}
	bytes, err := os.ReadFile(notesDbPath)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(bytes, &notesDb)
	if err != nil {
		log.Panic(err)
	}
	return &notesDb
}

func saveNotes(notesDb *NotesDb) {
	notesDbPath := resolveDbPath()
	file, err := os.Create(notesDbPath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	jsonStr, err := json.Marshal(notesDb)
	if err != nil {
		log.Panic(err)
	}
	_, err = file.Write(jsonStr)
	if err != nil {
		log.Panic(err)
	}
}

func resolveDbPath() string {
	str := os.Getenv("SIMPLYNOTES_DBPATH")
	if len(str) == 0 {
		str = filepath.Join(os.Getenv("HOME"), ".simplynotes.json")
	}
	return str
}

func (ndb *NotesDb) print() {
	prefix := ""
	for i, v := range ndb.Notes {
		if v.Done {
			prefix = fmt.Sprintf("[DONE %s] ", time.Unix(v.DoneTs, 0).Format("2006-01-02 15:04:05"))
		} else {
			prefix = ""
		}
		fmt.Printf("%d: %s%s\n", i, prefix, v.Content)
	}
}

func (ndb *NotesDb) complete(idxStr string) {
	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		log.Panic(err)
	}
	ndb.Notes[idx].Done = true
	ndb.Notes[idx].DoneTs = time.Now().Unix()
	ndb.print()
}

func (ndb *NotesDb) add(args []string) {
	str := strings.Join(args, " ")
	note := Note{str, false, 0}
	ndb.Notes = append(ndb.Notes, &note)
}


