package mappers

import (
	"encoding/json"
	"io"

	"deduplicator/internal/models"
)

func ReaderToRequest(body io.Reader) (*models.Request, error) {
	msg := &models.Request{}
	err := json.NewDecoder(body).Decode(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
