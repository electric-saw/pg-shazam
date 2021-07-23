package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
	c := NewShazam()
	path, err := filepath.Abs("../../../conf/sample.yaml")
	assert.Nil(t, err)

	err = c.LoadFromFile(path)

	assert.Nil(t, err)
}
