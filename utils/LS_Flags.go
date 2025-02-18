package utils

import (
	"os"
	"strings"
)

func Args() (path string, lFlag, RFlag, aFlag, rFlag, tFlag bool) {

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "a") {
				aFlag = true
			}
			if strings.Contains(arg, "R") {
				RFlag = true
			}
			if strings.Contains(arg, "t") {
				tFlag = true
			}
			if strings.Contains(arg, "l") {
				lFlag = true
			}
			if strings.Contains(arg, "r") {
				rFlag = true
			}
		} else {
			path = arg
		}
	}

	return
}
