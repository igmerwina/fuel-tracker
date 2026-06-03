package main

import "os"

func getenvRaw(key string) string {
	return os.Getenv(key)
}
