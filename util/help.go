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

func CreateExchangeName(session_id string, role string) string {
	return fmt.Sprintf("%s-exchange-%s", session_id, role)
}

func CreateQueueName(session_id string, name string) string {
	return fmt.Sprintf("%s-user-%s", session_id, name)
}
