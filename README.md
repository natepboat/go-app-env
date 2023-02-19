# GO APP ENV

Library to read JSON file and transform into flatten configuration map with ability to overwrite config by using specific environment file or environment variable.

## Active environment

Active environment can configure from environment variable name `GO_APP_ACTIVE_ENV` or command line argument `--env`.

In case of found both environment variable and command line argument, AppEnv will use value from **environment variable**

If both not provide, default active environment as `default`

## Config Directory

Config directory can configure from environment variable name `GO_APP_CONFIG_DIR` or command line argument `--configDir`. 

In case of found both environment variable and command line argument, AppEnv will use value from **environment variable**

If both not provide, default active environment as `./resources/`

---
## Configuration Mapping

By default, App env will read JSON file name `config.json` under config directory. Will return empty map of configuration in case of file not found.

### Flatten Configuration

App env will read JSON and transform into 1 level map by joining key with dot (`.`).
For example, following JSON

```json
{
    "app":{
        "version": 1,
        "name": "demo",
        "meta": {
            "something": "value"
        }
    },
    "active": true
}
```

will transform to key-value map

```
"app.version" = 1
"app.name" = "demo"
"app.meta.something" = "value"
"active" = true
```

### active environment overwrite

In case of active environment is not `default`, App env will trying to read `config-{activeEnv}.json` under config directory. In case of file not found, will fallback to use configuration from `config.json`.

If `config-{activeEnv}.json` found, will merge configuration with `config.json` by take value from active env configuration as priority.

``` json
Ex.
// config.json
{
    "app":{
        "version": 1,
        "name": "demo",
        "meta": {
            "something": "value"
        }
    },
    "active": true
}

// config-{activeEnv}.json
{
    "app":{
        "version": 2
    },
    "active": true
    "provider": "internal"
}

// merge result
{
    "app":{
        "version": 2,
        "name": "demo",
        "meta": {
            "something": "value"
        }
    },
    "active": true
    "provider": "internal"
}
```

### Environment variable overwrite

App env will try to read flatten configuration key, by replace dash (`-`) and dot (`.`) with underscore (`_`) and change all character to upper case. If no variable found will use active env config or default.

```
Ex.
app.version >> APP_VERSION
app.api.get-user >> APP_API_GET_USER
```

---
## Usage

```go
appEnv := NewAppEnv(os.DirFS("{desired directory}"))
```

## Example

```go

// in file main.go
...
app := appenv.NewAppEnv(os.DirFS("."))
fmt.Printf("env: %s, configDir:%s \n", app.ActiveEnv(), app.ConfigDir())
...

/*
run:
    go run main.go
result:
    env: default, configDir: ./resources/

run:
    go run main.go --env=test
result:
    env: test, configDir: ./resources/
    
run:
    go run main.go --env=test --configDir=./cfg
result:
    env: test, configDir: ./cfg
    
run:
    GO_APP_ACTIVE_ENV=local GO_APP_CONFIG_DIR=./tmp go run main.go --env=test --configDir=./cfg
result:
    env: local, configDir: ./tmp
*/
```