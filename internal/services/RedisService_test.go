package services

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Set(t *testing.T) {
	miniRedis, err := miniredis.Run()
	if err != nil {
		print(err.Error())
		t.Fail()
	}

	redisService, err := NewRedisService(miniRedis.Addr(), "", 0)
	assert.ErrorIs(t, err, nil)

	err = redisService.Set("test", map[string]interface{}{"deviceID": "AZERTY"})
	assert.ErrorIs(t, err, nil)

	result, err := redisService.Get("test")
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, result["deviceID"], "AZERTY")
}

func Test_Get(t *testing.T) {
	miniRedis, err := miniredis.Run()
	assert.ErrorIs(t, err, nil)

	redisService, err := NewRedisService(miniRedis.Addr(), "", 0)
	assert.ErrorIs(t, err, nil)

	err = redisService.Set("test", map[string]interface{}{"deviceID": "AZERTY", "content": "TEST"})
	assert.ErrorIs(t, err, nil)

	result, err := redisService.Get("test")
	assert.ErrorIs(t, err, nil)

	assert.Equal(t, result["deviceID"], "AZERTY")

	assert.Equal(t, result["content"], "TEST")
}
