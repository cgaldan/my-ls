package main

import (
	"ls/logic"
	"ls/utils"
)

func main() {
	paths, lFlag, RFlag, aFlag, rFlag, tFlag := utils.Args()

	if len(paths) == 0 {
		paths = []string{"."}
	} else if len(paths) == 1 {
		paths = []string{paths[0]}
	}
	logic.ProcessPaths(paths, lFlag, RFlag, aFlag, rFlag, tFlag)
}
