package main

import (
	"fmt"
	"ls/logic"
	"ls/utils"
)

func main() {
	var dir string
	paths, lFlag, RFlag, aFlag, rFlag, tFlag := utils.Args()

	if len(paths) > 1 {
		for _, path := range paths {
			fmt.Printf("%s:\n", path)
			logic.TheMainLS(path, lFlag, RFlag, aFlag, rFlag, tFlag)
			fmt.Println()
		}
		return
	}

	logic.TheMainLS(dir, lFlag, RFlag, aFlag, rFlag, tFlag)
}
