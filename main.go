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
		utils.SortDirs(&paths)
		for i, path := range paths {
			fmt.Printf("%s:\n", path)
			logic.TheMainLS(path, lFlag, RFlag, aFlag, rFlag, tFlag)
			if i != len(paths)-1 {
				fmt.Println()
			}
		}
		return
	} else if len(paths) == 1 {
		dir = paths[0]
	}

	logic.TheMainLS(dir, lFlag, RFlag, aFlag, rFlag, tFlag)
}
