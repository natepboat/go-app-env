package internal

import (
	"os"
	"path"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestGetActiveEnv(t *testing.T) {
	p := NewProcessor(fstest.MapFS{})

	t.Run("default env", func(t *testing.T) {
		assert.Equal(t, "default", p.GetActiveEnv())
	})

	t.Run("non-default env", func(t *testing.T) {
		os.Args = []string{"program", "--env=test"}
		assert.Equal(t, "test", p.GetActiveEnv())
		os.Args = []string{}
	})
}

func TestGetConfigDir(t *testing.T) {
	p := NewProcessor(fstest.MapFS{})

	t.Run("default dir", func(t *testing.T) {
		assert.Equal(t, "./resources/", p.GetConfigDir())
	})

	t.Run("non-default dir", func(t *testing.T) {
		os.Args = []string{"program", "--configDir=/tmp/configs/"}
		assert.Equal(t, "/tmp/configs/", p.GetConfigDir())
		os.Args = []string{}
	})
}

func TestLoadConfig(t *testing.T) {
	fsys := fstest.MapFS{
		path.Join("resources", "config.json"): &fstest.MapFile{
			Data: []byte(`{"app":{"name":"default","version":1,"meta":"default-value"}}`),
		},
		path.Join("resources", "config-test.json"): &fstest.MapFile{
			Data: []byte(`{"app":{"name":"test","version":2,"test":"test-value"}}`),
		},
	}
	p := NewProcessor(fsys)

	t.Run("default config", func(t *testing.T) {
		data := p.LoadConfig("default", "resources")
		assert.Equal(t, "default", data["app.name"])
		assert.EqualValues(t, 1, data["app.version"])
		assert.Equal(t, "default-value", data["app.meta"])
	})

	t.Run("active env config", func(t *testing.T) {
		data := p.LoadConfig("test", "resources")
		assert.Equal(t, "test", data["app.name"])
		assert.EqualValues(t, 2, data["app.version"])
		assert.Equal(t, "default-value", data["app.meta"])
		assert.Equal(t, "test-value", data["app.test"])
	})

	t.Run("active env but without config exist", func(t *testing.T) {
		data := p.LoadConfig("local", "resources")
		assert.Equal(t, "default", data["app.name"])
		assert.EqualValues(t, 1, data["app.version"])
		assert.Equal(t, "default-value", data["app.meta"])
	})

	t.Run("env var config", func(t *testing.T) {
		os.Setenv("APP_VERSION", "999")

		data := p.LoadConfig("test", "resources")
		assert.Equal(t, "test", data["app.name"])
		assert.EqualValues(t, "999", data["app.version"])
		assert.Equal(t, "default-value", data["app.meta"])
		assert.Equal(t, "test-value", data["app.test"])

		os.Clearenv()
	})
}
