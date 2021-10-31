package mappers

import (
	"testing"

	"deduplicator/internal/models"

	"github.com/stretchr/testify/assert"
)

func Test_MessageToMap(t *testing.T) {
	attributes := make(map[string]string)
	attributes["hash_1"] = "FFFAAA"
	attributes["hash_2"] = "abcde"
	msg := models.Message{Date: "25-02-2021 15:04:05.000000 UTC", Count: 2, Attributes: attributes}
	mappedMsg := MessageToMap(msg)
	assert.Equal(t, msg.Count, mappedMsg["count"])
	assert.Equal(t, msg.Date, mappedMsg["date"])
	assert.Equal(t, msg.Attributes["hash_1"], mappedMsg["hash_1"])
	assert.Equal(t, msg.Attributes["hash_2"], mappedMsg["hash_2"])
}

func Test_MapToMessage(t *testing.T) {
	mappedMsg := map[string]string{
		"hash_1": "FFFAAA",
		"date":   "25-02-2021 15:04:05.000000 UTC",
		"count":  "2",
	}

	msg := MapToMessage(mappedMsg)

	assert.Equal(t, msg.Count, 2)
	assert.Equal(t, msg.Attributes["hash_1"], "FFFAAA")
	assert.Equal(t, msg.Date, "25-02-2021 15:04:05.000000 UTC")
}

func Test_RawDataToMessage(t *testing.T) {
	msg := RawDataToMessage("25-02-2021 15:04:05.000000 UTC", 1, []string{"90DCCBBD587706546543", "TEST"})
	assert.Equal(t, msg.Count, 1)
	assert.Equal(t, msg.Date, "25-02-2021 15:04:05.000000 UTC")
	assert.Equal(t, msg.Attributes["hash_1"], "90DCCBBD587706546543")
	assert.Equal(t, msg.Attributes["hash_2"], "TEST")
}
