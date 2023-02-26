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

func TestEnvOrDefault(t *testing.T) {
	t.Run("nil appEnv", func(t *testing.T) {
		assert.Equal(t, "default value", ConfigOrDefault(nil, "config.key", "default value"))
	})

	t.Run("key not exist", func(t *testing.T) {
		fsys := fstest.MapFS{
			path.Join("resources", "config.json"): &fstest.MapFile{
				Data: []byte(`{"config":{"key-x":"key-x-val"}}`),
			},
		}
		app := NewAppEnv(fsys)

		assert.Equal(t, "default value", ConfigOrDefault(app, "config.key", "default value"))
	})

	t.Run("key exist", func(t *testing.T) {
		fsys := fstest.MapFS{
			path.Join("resources", "config.json"): &fstest.MapFile{
				Data: []byte(`{"config":{"key":"key-val"}}`),
			},
		}
		app := NewAppEnv(fsys)

		assert.Equal(t, "key-val", ConfigOrDefault(app, "config.key", "default value"))
	})
}
