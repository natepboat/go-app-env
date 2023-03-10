package goappenv

import (
	"context"
	"io/fs"

	"github.com/natepboat/go-app-env/contextKey"
	"github.com/natepboat/go-app-env/internal"
)

type IAppEnv interface {
	ActiveEnv() string
	ConfigDir() string
	Config() map[string]interface{}
	Context() context.Context
}

type AppEnv struct {
	activeEnv string
	configDir string
	config    map[string]interface{}
	context   context.Context
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

func (a *AppEnv) Context() context.Context {
	return a.context
}

func NewAppEnv(fsys fs.FS, ctx context.Context) *AppEnv {
	processor := internal.NewProcessor(fsys)
	activeEnv := processor.GetActiveEnv()
	configDir := processor.GetConfigDir()

	return &AppEnv{
		activeEnv: activeEnv,
		configDir: configDir,
		config:    processor.LoadConfig(activeEnv, configDir),
		context:   context.WithValue(ctx, contextKey.AppCtxKey{}, make(map[string]string, 0)),
	}
}

func ConfigOrDefault(appenv *AppEnv, key string, defaultValue interface{}) interface{} {
	if appenv != nil {
		val, exist := appenv.Config()[key]

		if exist && val != nil {
			return val
		}
	}

	return defaultValue
}
