package goappenv

import (
	"path"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestNewAppEnv(t *testing.T) {
	fsys := fstest.MapFS{
		path.Join("resources", "config.json"): &fstest.MapFile{
			Data: []byte(`{"app":{"name":"default","version":1,"meta":"default-value"}}`),
		},
	}
	app := NewAppEnv(fsys)
	assert.Equal(t, "default", app.ActiveEnv())
	assert.Equal(t, "./resources/", app.ConfigDir())
	assert.Equal(t, "default", app.Config()["app.name"])
	assert.EqualValues(t, 1, app.Config()["app.version"])
	assert.Equal(t, "default-value", app.Config()["app.meta"])
}
