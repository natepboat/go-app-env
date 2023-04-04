package internal

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func toEnvVarName(key string) string {
	envVarName := strings.ReplaceAll(strings.ReplaceAll(key, ".", "_"), "-", "_")
	return strings.ToUpper(envVarName)
}

func loadJsonFile(fsys fs.FS, fileName string) map[string]interface{} {
	fileContent, err := fs.ReadFile(fsys, fileName)
	if err != nil {
		log.Println("Failed to load configuration file", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		log.Println("Failed to parse configuration", err)
	}

	return data
}

func loadEnvArgument(topic, envName, argName, defaultValue string) string {
	envVar := os.Getenv(envName)
	if len(envVar) > 0 {
		log.Printf("Use %s from environment variable %s \n", topic, envName)
		return envVar
	} else if len(os.Args) > 1 {
		arguments := os.Args[1:]
		for _, arg := range arguments {
			if strings.HasPrefix(arg, argName) {
				argValue := strings.Split(arg, "=")

				log.Printf("Use %s from command line argument \n", topic)
				return argValue[1]
			}
		}
	}

	return defaultValue
}

func flatMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			flatChild := flatMap(child)
			for flatternedKey, flatternedValue := range flatChild {
				out[k+"."+flatternedKey] = flatternedValue
			}
		default:
			kind := reflect.TypeOf(child).Kind()
			if kind == reflect.Slice || kind == reflect.Array {
				flatMapCollection(child, k, out)
			} else {
				out[k] = child
			}
		}
	}
	return out
}

func flatMapCollection(object interface{}, key string, out map[string]interface{}) {
	arrVAlue := reflect.ValueOf(object)
	for i := 0; i < arrVAlue.Len(); i++ {
		indexString := strconv.Itoa(i)

		switch item := arrVAlue.Index(i).Interface().(type) {
		case map[string]interface{}:
			flatChild := flatMap(item)
			for flatternedKey, flatternedValue := range flatChild {
				out[key+"."+indexString+"."+flatternedKey] = flatternedValue
			}
		default:
			out[key+"."+indexString] = item
		}
	}
}
