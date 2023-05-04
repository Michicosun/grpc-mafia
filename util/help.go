package util

import (
	"fmt"
	"os"
)

func GetEnvWithDefault(key string, default_value string) string {
	value, has_env := os.LookupEnv(key)

	if has_env {
		return value
	}

	return default_value
}

func ChatGroupName(session_id string, subgroup string) string {
	return fmt.Sprintf("%s-%s", session_id, subgroup)
}
