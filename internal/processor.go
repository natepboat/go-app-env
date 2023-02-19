package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
)

const activeEnvName = "GO_APP_ACTIVE_ENV"
const configDirEnvName = "GO_APP_CONFIG_DIR"
const defaultEnvironment = "default"
const defaultConfigDir = "./resources/"

type processor struct {
	fsys fs.FS
}

func NewProcessor(fsys fs.FS) processor {
	return processor{
		fsys: fsys,
	}
}

func (p processor) GetActiveEnv() string {
	topic := "active environment"
	activeEnv := loadEnvArgument(topic, activeEnvName, "--env", defaultEnvironment)
	log.Printf("::: Application %s = %s :::\n", topic, activeEnv)
	return activeEnv
}

func (p processor) GetConfigDir() string {
	topic := "config directory"
	configDir := loadEnvArgument(topic, configDirEnvName, "--configDir", defaultConfigDir)
	log.Printf("::: Application %s = %s :::\n", topic, configDir)
	return configDir
}

func (p processor) LoadConfig(activeEnv, configDir string) map[string]interface{} {
	defaultConfig := loadJsonFile(p.fsys, path.Join(configDir, "config.json"))
	defaultConfig = flatMap(defaultConfig)

	if !strings.EqualFold(activeEnv, defaultEnvironment) {
		activeEnvConfig := loadJsonFile(p.fsys, path.Join(configDir, fmt.Sprintf("config-%s.json", activeEnv)))
		activeEnvConfig = flatMap(activeEnvConfig)
		return merge(defaultConfig, activeEnvConfig)
	}

	return defaultConfig
}

func merge(defaultConfig, activeEnvConfig map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for k, v := range activeEnvConfig {
		out[k] = v
	}

	for k, v := range defaultConfig {
		_, exist := out[k]
		if !exist {
			out[k] = v
		}
	}

	for k := range out {
		envVar := os.Getenv(toEnvVarName(k))
		if len(envVar) > 0 {
			out[k] = envVar
		}
	}

	return out
}
