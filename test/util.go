package test

import (
	"fmt"
	"math/rand"
	"os"
)

var (
	apiPort      = getEnvOrDefault("API_PORT", "8080")
	portsService = fmt.Sprintf("http://localhost:%s/api/v1/ports", apiPort)
	testFile     = getEnvOrDefault("TEST_FILE", "./assets/ports.json")
)

func getEnvOrDefault(env, def string) string {
	e := os.Getenv(env)
	if len(e) == 0 {
		return def
	}

	return e
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))] // nolint: gosec
	}
	return string(b)
}
