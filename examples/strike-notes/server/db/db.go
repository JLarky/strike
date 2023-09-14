package db

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

type Note struct {
	Id        uint64    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SearchNotes(q string) ([]Note, error) {
	var out []Note
	notes, err := readNotes()
	if err != nil {
		return out, err
	}
	search := strings.ToLower(q)
	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Title), search) || strings.Contains(strings.ToLower(note.Body), search) {
			out = append(out, note)
		}
	}
	time.Sleep(time.Duration(math.Min(20, math.Pow(2, float64(len(q))))) * 100 * time.Millisecond)
	return out, nil
}

func readNotes() ([]Note, error) {
	jsonFile, err := os.Open("server/notes/notes.json")
	if err != nil {
		return []Note{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var notes []Note
	err = json.Unmarshal(byteValue, &notes)
	if err != nil {
		return []Note{}, err
	}
	return notes, nil
}
