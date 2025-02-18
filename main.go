package main

import (
	"ls/myls"
	"ls/utils"
)

func main() {
	path, lFlag, RFlag, aFlag, rFlag, tFlag := utils.Args()

	myls.TheMainLS(path, lFlag, RFlag, aFlag, rFlag, tFlag)

}
