package ls

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	blue  = "\033[34m"
	reset = "\033[0m"
)

func Args() {
	dir := flag.String("dir", ".", "Directory to list")
	flag.Parse()

	ls(*dir)
}

func ls(dirName string) {
	dirPath, err := filepath.Abs(dirName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Printf(blue+"%s  "+reset, file.Name())
		} else {
			fmt.Printf("%s  ", file.Name())
		}
	}
	fmt.Println()
}
