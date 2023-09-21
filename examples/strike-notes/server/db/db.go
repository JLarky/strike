package db

import (
	_ "embed"
	"encoding/json"
	"math"
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
	if q != "" {
		time.Sleep(time.Duration(math.Min(20, math.Pow(2, float64(len(q))))) * 100 * time.Millisecond)
	}
	return out, nil
}

//go:embed notes.json
var jsonFile []byte

func readNotes() ([]Note, error) {
	var notes []Note
	err := json.Unmarshal(jsonFile, &notes)
	if err != nil {
		return []Note{}, err
	}
	return notes, nil
}
