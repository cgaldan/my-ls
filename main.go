package main

import (
	"ls/logic"
	"ls/utils"
)

func main() {
	path, lFlag, RFlag, aFlag, rFlag, tFlag := utils.Args()

	logic.TheMainLS(path, lFlag, RFlag, aFlag, rFlag, tFlag)
}
