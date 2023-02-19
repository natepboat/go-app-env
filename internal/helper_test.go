package internal

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestToEnvVarName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "app",
			expected: "APP",
		},
		{
			input:    "app.version",
			expected: "APP_VERSION",
		},
		{
			input:    "app.meta.name",
			expected: "APP_META_NAME",
		},
		{
			input:    "key.some-val.prop",
			expected: "KEY_SOME_VAL_PROP",
		},
		{
			input:    "key.prop.1.child",
			expected: "KEY_PROP_1_CHILD",
		},
	}

	for _, tc := range testCases {
		output := toEnvVarName(tc.input)
		assert.Equal(t, tc.expected, output)
	}
}

func TestLoadJsonFile(t *testing.T) {
	fileSystem := fstest.MapFS{
		"tmp/resources/config.json": &fstest.MapFile{
			Data: []byte(`{"app":{"name":"test","version":1}}`),
		},
		"tmp/resources/invalid.json": &fstest.MapFile{
			Data: []byte("not json"),
		},
		"tmp/resources/empty.json": &fstest.MapFile{
			Data: []byte(""),
		},
	}

	t.Run("file not exist", func(t *testing.T) {
		data := loadJsonFile(fileSystem, "tmp/resources/not-exist.json")
		assert.Equal(t, 0, len(data))
	})
	t.Run("empty file", func(t *testing.T) {
		data := loadJsonFile(fileSystem, "tmp/resources/empty.json")
		assert.Equal(t, 0, len(data))
	})
	t.Run("valid json", func(t *testing.T) {
		data := loadJsonFile(fileSystem, "tmp/resources/config.json")
		assert.Equal(t, 1, len(data))
	})
	t.Run("invalid json", func(t *testing.T) {
		data := loadJsonFile(fileSystem, "tmp/resources/invalid.json")
		assert.Equal(t, 0, len(data))
	})
}

func TestLoadEnvArgument(t *testing.T) {
	topic := "testkey1"
	envName := "TEST_KEY_1"
	argName := "--testkey"
	defaultValue := "test-key-val"

	t.Run("default value", func(t *testing.T) {
		output := loadEnvArgument(topic, envName, argName, defaultValue)
		assert.Equal(t, defaultValue, output)
	})

	t.Run("env value", func(t *testing.T) {
		os.Setenv(envName, "env value")

		output := loadEnvArgument(topic, envName, argName, defaultValue)
		assert.Equal(t, "env value", output)

		os.Clearenv()
	})

	t.Run("arg value", func(t *testing.T) {
		os.Args = []string{"program", argName + "=arg value"}

		output := loadEnvArgument(topic, envName, argName, defaultValue)
		assert.Equal(t, "arg value", output)

		os.Args = []string{}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("non-empty map", func(t *testing.T) {
		lv3Obj := make(map[string]interface{})
		lv3Obj["leafStr"] = "leaf-string"
		lv3Obj["leafNum"] = 1230
		lv3Obj["leafBool"] = true
		lv2Obj := make(map[string]interface{})
		lv2Obj["third"] = lv3Obj
		lv2Obj["leaf"] = "2nd-leaf"
		lv1Obj := make(map[string]interface{})
		lv1Obj["second"] = lv2Obj
		lv1Obj["leaf"] = "1st-leaf"
		mockMap := make(map[string]interface{})
		mockMap["first"] = lv1Obj
		mockMap["meta"] = "meta-string"
		mockMap["version"] = 789

		flattedMap := flatMap(mockMap)

		assert.Equal(t, "leaf-string", flattedMap["first.second.third.leafStr"])
		assert.Equal(t, 1230, flattedMap["first.second.third.leafNum"])
		assert.Equal(t, true, flattedMap["first.second.third.leafBool"])
		assert.Equal(t, "2nd-leaf", flattedMap["first.second.leaf"])
		assert.Equal(t, "1st-leaf", flattedMap["first.leaf"])
		assert.Equal(t, "meta-string", flattedMap["meta"])
		assert.Equal(t, 789, flattedMap["version"])

		_, unknownExist := flattedMap["first.second.third.unknown"]
		assert.Equal(t, false, unknownExist)

	})

	t.Run("empty map", func(t *testing.T) {
		mockMap := make(map[string]interface{})

		flattedMap := flatMap(mockMap)

		assert.Equal(t, 0, len(flattedMap))
	})
}
