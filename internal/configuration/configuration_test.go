package configuration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	_, err := GetConfiguration()
	assert.Error(t, err)

	*listID = "ls001"
	cfg, err := GetConfiguration()
	assert.NoError(t, err)
	assert.Equal(t, []string{"ls001"}, cfg.ListID)
}
