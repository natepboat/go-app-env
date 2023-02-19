package goappenv

import (
	"io/fs"

	"github.com/natepboat/go-app-env/internal"
)

type appEnv struct {
	activeEnv string
	configDir string
	config    map[string]interface{}
}

func (a appEnv) ActiveEnv() string {
	return a.activeEnv
}

func (a appEnv) ConfigDir() string {
	return a.configDir
}

func (a appEnv) Config() map[string]interface{} {
	return a.config
}

func NewAppEnv(fsys fs.FS) appEnv {
	processor := internal.NewProcessor(fsys)
	activeEnv := processor.GetActiveEnv()
	configDir := processor.GetConfigDir()

	return appEnv{
		activeEnv: activeEnv,
		configDir: configDir,
		config:    processor.LoadConfig(activeEnv, configDir),
	}
}
