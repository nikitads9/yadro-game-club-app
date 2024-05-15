package main

import (
	"flag"
	"log"
	"os"

	"github.com/nikitads9/yadro-game-club-app/internal/process"
)

var path string

func init() {
	flag.StringVar(&path, "path", "C:\\Users\\swnik\\Desktop\\projects\\yadro-game-club-app\\testdata\\test2.txt", "путь к файлу с исходными данными по итогам дня")
}

func main() {
	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("could not open file with path %s. error: %v", path, err)
	}
	defer file.Close() //nolint:errcheck

	process.ReadLogs(file)
}
