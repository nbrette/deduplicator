package mappers

import (
	"strconv"
	"strings"

	"deduplicator/internal/models"
)

const (
	msgDate          = "date"
	count            = "count"
	hashAttributeKey = "hash_"
)

func MessageToMap(msg models.Message) map[string]interface{} {
	mappedMessage := make(map[string]interface{})
	mappedMessage[msgDate] = msg.Date
	mappedMessage[count] = msg.Count
	for key, value := range msg.Attributes {
		mappedMessage[key] = value
	}
	return mappedMessage
}

func MapToMessage(data map[string]string) models.Message {
	var msg models.Message
	deduplicatorCount, _ := strconv.Atoi(data[count])
	deduplicatorAttribute := make(map[string]string)
	for key, value := range data {
		if strings.Contains(key, hashAttributeKey) {
			deduplicatorAttribute[key] = value
		}
	}
	msg.Attributes = deduplicatorAttribute
	msg.Count = deduplicatorCount
	msg.Date = data[msgDate]
	return msg
}

func RawDataToMessage(date string, count int, attributes []string) models.Message {
	deduplicatorAttributes := make(map[string]string)
	for i := 0; i < len(attributes); i++ {
		key := hashAttributeKey + strconv.Itoa(i+1)
		deduplicatorAttributes[key] = attributes[i]
	}
	return models.Message{Date: date, Count: count, Attributes: deduplicatorAttributes}
}
