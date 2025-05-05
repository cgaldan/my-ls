package logic

import (
	"fmt"
	"ls/data"
	"strconv"
	"strings"
	"time"
)

func formatLongEntry(file data.MyLSFiles, lenNLink int, maxOwner, maxGroup, maxSize, maxMajor, maxMinor int) string {
	permission := getPermission(file)

	modTime := formatTime(file.ModTime)

	size := fmt.Sprintf("%*d", maxSize, file.Size)

	fileName := formatFileNames(file.Name)

	if file.IsBlockDevice || file.IsCharDevice {
		size = fmt.Sprintf("%*d, %*d",
			maxMajor, file.MajorNumber,
			maxMinor, file.MinorNumber,
		)
	}

	fileColor := file.GetColor()
	if file.IsLink {
		targetColor := reset
		if file.TargetFile != nil {
			targetColor = file.TargetFile.GetColor()
		}
		fileName = fmt.Sprintf("%s%s%s -> %s%s%s", fileColor, fileName, reset, targetColor, file.LinkTarget, reset)
	}

	// fmt.Println(maxSize)
	return fmt.Sprintf("%10s %*d %-*s %-*s %s %12s %s%s%s",
		permission,
		lenNLink, file.NLink,
		maxOwner, file.OwnerName,
		maxGroup, file.GroupName,
		size,
		modTime,
		fileColor,
		fileName,
		reset,
	)
}

func getPermission(file data.MyLSFiles) string {
	permission := file.Mode.String()

	if file.IsLink {
		permission = strings.Replace(permission, "L", "l", 1)
	}
	if file.IsBlockDevice {
		permission = strings.Replace(permission, "D", "b", 1)
	}
	if file.IsCharDevice {
		permission = strings.Replace(permission, "D", "", 1)
	}
	if file.IsSetuid {
		permission = strings.Replace(permission, "u", "-", 1)
	}
	if file.IsSetgid {
		permission = strings.Replace(permission, "g", "-", 1)
	}

	return permission
}

func updateMaxNlink(maxNlink *int, file data.MyLSFiles) {
	strNLink := strconv.Itoa(int(file.NLink))
	if *maxNlink < len(strNLink) {
		*maxNlink = len(strNLink)
	}
}

func calculateMaxWidth(files []data.MyLSFiles) (maxOwner, maxGroup, maxSize, maxMajor, maxMinor int) {
	maxMajor, maxMinor, maxRegular := 0, 0, 0

	for _, file := range files {
		if len(file.OwnerName) > maxOwner {
			maxOwner = len(file.OwnerName)
		}
		if len(file.GroupName) > maxGroup {
			maxGroup = len(file.GroupName)
		}

		if file.IsBlockDevice || file.IsCharDevice {
			majorLen := len(fmt.Sprint(file.MajorNumber))
			minorLen := len(fmt.Sprint(file.MinorNumber))

			if majorLen > maxMajor {
				maxMajor = majorLen
			}
			if minorLen > maxMinor {
				maxMinor = minorLen
			}

		} else {
			sizeLen := len(fmt.Sprintf("%d", file.Size))
			if sizeLen > maxRegular {
				maxRegular = sizeLen
			}
		}
	}

	deviceField := maxMajor + maxMinor + 2
	if deviceField > maxRegular && deviceField > 2 {
		maxSize = deviceField
	} else {
		maxSize = maxRegular
	}

	return maxOwner, maxGroup, maxSize, maxMajor, maxMinor
}

func formatTime(modTime time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0)

	if now.Sub(modTime) > 0 && modTime.After(sixMonthsAgo) {
		return modTime.Format("Jan _2 15:04")
	} else {
		return modTime.Format("Jan _2  2006")
	}
}
