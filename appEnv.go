package goappenv

import (
	"io/fs"

	"github.com/natepboat/go-app-env/internal"
)

type AppEnv struct {
	activeEnv string
	configDir string
	config    map[string]interface{}
}

func (a *AppEnv) ActiveEnv() string {
	return a.activeEnv
}

func (a *AppEnv) ConfigDir() string {
	return a.configDir
}

func (a *AppEnv) Config() map[string]interface{} {
	return a.config
}

func NewAppEnv(fsys fs.FS) *AppEnv {
	processor := internal.NewProcessor(fsys)
	activeEnv := processor.GetActiveEnv()
	configDir := processor.GetConfigDir()

	return &AppEnv{
		activeEnv: activeEnv,
		configDir: configDir,
		config:    processor.LoadConfig(activeEnv, configDir),
	}
}
