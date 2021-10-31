package services

import (
	"time"

	"deduplicator/internal/services/hasher"
	"deduplicator/internal/services/mappers"
)

type DeduplicatorService struct {
	redisService *RedisService
}

//go:generate mockgen -destination=$PROJECT_PATH/internals/tests/mocks/mock_IDeduplicatorService.go -package=mocks . IDeduplicatorService
type IDeduplicatorService interface {
	CheckHashExists(hash string) (bool, error)
	InsertNewHash(hash string, attributes []string) error
	UpdateExistingHash(hash string) error
	CreateOrUpdate(deduplicatorAttributes []string) (bool, error)
}

func NewDeduplicatorService(redisService *RedisService) *DeduplicatorService {
	return &DeduplicatorService{redisService: redisService}
}

func (deduplicatorService *DeduplicatorService) CheckHashExists(hash string) (bool, error) {
	data, err := deduplicatorService.redisService.Get(hash)
	if err != nil {
		return false, err
	}
	msg := mappers.MapToMessage(data)
	if msg.Count == 0 {
		return false, nil
	}

	return true, nil
}

func (deduplicatorService *DeduplicatorService) InsertNewHash(hash string, attributes []string) error {
	newMsg := mappers.RawDataToMessage(time.Now().Format("01-02-2006 15:04:05.000000 UTC"), 1, attributes)
	err := deduplicatorService.redisService.Set(hash, mappers.MessageToMap(newMsg))
	if err != nil {
		return err
	}
	return nil
}

func (deduplicatorService *DeduplicatorService) UpdateExistingHash(hash string) error {
	msgMapped, err := deduplicatorService.redisService.Get(hash)
	if err != nil {
		return err
	}
	msg := mappers.MapToMessage(msgMapped)
	msg.Count++
	err = deduplicatorService.redisService.Set(hash, mappers.MessageToMap(msg))
	if err != nil {
		return err
	}
	return nil
}

func (deduplicatorService *DeduplicatorService) CreateOrUpdate(deduplicatorAttributes []string) (bool, error) {
	hash := hasher.GenerateHash(deduplicatorAttributes...)
	exists, err := deduplicatorService.CheckHashExists(hash)
	if err != nil {
		return false, err
	}
	if exists {
		err = deduplicatorService.UpdateExistingHash(hash)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	err = deduplicatorService.InsertNewHash(hash, deduplicatorAttributes)
	if err != nil {
		return false, err
	}

	return true, nil
}
