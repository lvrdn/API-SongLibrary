package song

import (
	"encoding/json"
	"time"
)

// Song model info
// @Description song information
type Song struct {
	ID          int    `json:"id"`
	Name        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongHandler struct {
	Storage
	ExternalAPI string
}

type Storage interface {
	Add(name, group, releaseDate, text, link string) (int, error)
	Delete(id int) error
	Update(id int, releaseDate, text, link string) error
	Get(id int) (string, error)
	GetAll(limit, offset int, songNameFragment, groupNameFragment, year, textFragment, linkExist string) ([]*Song, error)
	GetErrorAlreadyExist() error
	GetErrorNoUpdate() error
	GetErrorBadID() error
	GetErrorBadDate() error
	GetErrorNoRows() error
}

// @Description response format
type Response map[string]interface{}

func GetAnswerWithError(message, path string) ([]byte, error) {
	answer := Response{
		"error": Response{
			"timestamp": time.Now(),
			"message":   message,
			"path":      path,
		},
	}

	answerData, err := json.Marshal(answer)
	if err != nil {
		return nil, err
	}

	return answerData, nil
}
