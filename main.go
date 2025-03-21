package main

import (
	"fmt"
	"ls/logic"
	"ls/utils"
	"os"
)

func main() {
	paths, lFlag, RFlag, aFlag, rFlag, tFlag := utils.Args()
	if len(paths) > 1 {
		utils.SortDirs(&paths)
		for i, path := range paths {
			fileInfo, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("myls: cannot access '%s': No such file or direcory\n", path)
				} else {
					fmt.Println("Error:", err)
				}
				return
			}
			if fileInfo.IsDir() {
				if !RFlag {
					fmt.Println(path + ":")
				}

				logic.TheMainLS(path, lFlag, RFlag, aFlag, rFlag, tFlag)
			} else {
				logic.TheMainLS(path, lFlag, RFlag, aFlag, rFlag, tFlag)
			}
			if i != len(paths)-1 {
				fmt.Println()
			}
		}
		return
	}

	if len(paths) == 1 {
		logic.TheMainLS(paths[0], lFlag, RFlag, aFlag, rFlag, tFlag)
		return
	}

	logic.TheMainLS("", lFlag, RFlag, aFlag, rFlag, tFlag)
}
