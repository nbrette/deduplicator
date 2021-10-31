package services

import (
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func Test_CheckHashExists(t *testing.T) {
	miniRedis, err := miniredis.Run()
	assert.ErrorIs(t, err, nil)

	redisService, err := NewRedisService(miniRedis.Addr(), "", 0)
	assert.ErrorIs(t, err, nil)

	deduplicatorService := NewDeduplicatorService(redisService)
	deduplicatorAttributes := []string{"FFFAAA", "abcde"}
	err = deduplicatorService.InsertNewHash("azerty123456", deduplicatorAttributes)
	assert.ErrorIs(t, err, nil)

	exists, err := deduplicatorService.CheckHashExists("azerty123456")
	assert.ErrorIs(t, err, nil)

	assert.True(t, exists)
}

func Test_InsertNewHash(t *testing.T) {
	miniRedis, err := miniredis.Run()
	assert.ErrorIs(t, err, nil)

	redisService, err := NewRedisService(miniRedis.Addr(), "", 0)
	assert.ErrorIs(t, err, nil)

	deduplicatorService := NewDeduplicatorService(redisService)
	deduplicatorAttributes := []string{"FFFAAA", "abcde"}
	err = deduplicatorService.InsertNewHash("azerty123456", deduplicatorAttributes)
	assert.ErrorIs(t, err, nil)
}

func Test_UpdateExistingHash(t *testing.T) {
	miniRedis, err := miniredis.Run()
	assert.ErrorIs(t, err, nil)

	redisService, err := NewRedisService(miniRedis.Addr(), "", 0)
	assert.ErrorIs(t, err, nil)

	deduplicatorService := NewDeduplicatorService(redisService)
	deduplicatorAttributes := []string{"FFFAAA", "abcde"}
	err = deduplicatorService.InsertNewHash("azerty123456", deduplicatorAttributes)
	assert.ErrorIs(t, err, nil)

	err = deduplicatorService.UpdateExistingHash("azerty123456")
	assert.ErrorIs(t, err, nil)

	data, err := redisService.Get("azerty123456")
	assert.ErrorIs(t, err, nil)

	count, err := strconv.Atoi(data["count"])
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, 2, count)
}

func Test_CreateOrUpdate(t *testing.T) {
	miniRedis, err := miniredis.Run()
	assert.ErrorIs(t, err, nil)

	redisService, err := NewRedisService(miniRedis.Addr(), "", 0)
	assert.ErrorIs(t, err, nil)

	deduplicatorService := NewDeduplicatorService(redisService)
	deduplicatorAttributes := []string{"FFFAAA", "abcde"}
	err = deduplicatorService.InsertNewHash("azerty123456", deduplicatorAttributes)
	assert.ErrorIs(t, err, nil)

	mustForward, err := deduplicatorService.CreateOrUpdate(deduplicatorAttributes)
	assert.ErrorIs(t, err, nil)
	assert.True(t, mustForward)

	mustForward, err = deduplicatorService.CreateOrUpdate(deduplicatorAttributes)
	assert.ErrorIs(t, err, nil)
	assert.False(t, mustForward)
}
