package logic

import (
	"fmt"
	"ls/data"
	"os"
	"strconv"
	"time"
	"unicode"
)

func FormatLongEntry(file data.MyLSFiles, lenNLink int, maxOwner, maxGroup, maxSize, maxMajor, maxMinor int) string {
	permission := GetPermission(file)

	modTime := FormatTime(file.ModTime)

	size := fmt.Sprintf("%*d", maxSize, file.Size)

	fileName := FormatFileNames(file.Name)

	if file.IsBlockDevice || file.IsCharDevice {
		size = fmt.Sprintf("%*d, %*d",
			maxMajor, file.MajorNumber,
			maxMinor, file.MinorNumber,
		)
	}

	fileColor := file.GetColor()
	if file.IsLink {
		targetColor := Reset
		if file.TargetFile != nil {
			targetColor = file.TargetFile.GetColor()
		}
		fileName = fmt.Sprintf("%s%s%s -> %s%s%s", fileColor, fileName, Reset, targetColor, file.LinkTarget, Reset)
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
		Reset,
	)
}

func GetPermission(file data.MyLSFiles) string {
	mode := file.Mode
	perm := make([]byte, 10) // 1 type + 9 permissions

	// File type
	switch {
	case file.IsLink:
		perm[0] = 'l'
	case file.IsDir:
		perm[0] = 'd'
	case file.IsBlockDevice:
		perm[0] = 'b'
	case file.IsCharDevice:
		perm[0] = 'c'
	default:
		perm[0] = '-'
	}

	// Permissions
	perm[1] = boolToChar(mode&0400 != 0, 'r') // User read
	perm[2] = boolToChar(mode&0200 != 0, 'w') // User write
	perm[3] = specialChar(mode&0100 != 0, mode&os.ModeSetuid != 0, 'x', 's')

	perm[4] = boolToChar(mode&0040 != 0, 'r') // Group read
	perm[5] = boolToChar(mode&0020 != 0, 'w') // Group write
	perm[6] = specialChar(mode&0010 != 0, mode&os.ModeSetgid != 0, 'x', 's')

	perm[7] = boolToChar(mode&0004 != 0, 'r') // Others read
	perm[8] = boolToChar(mode&0002 != 0, 'w') // Others write
	perm[9] = specialChar(mode&0001 != 0, mode&os.ModeSticky != 0, 'x', 't')

	return string(perm)
}

// Helper functions
func boolToChar(has bool, char byte) byte {
	if has {
		return char
	}
	return '-'
}

func specialChar(hasPerm, hasSpecial bool, normal, special byte) byte {
	if hasSpecial {
		if hasPerm {
			return special
		}
		return byte(unicode.ToUpper(rune(special))) // S or T
	}
	return boolToChar(hasPerm, normal)
}

// func GetPermission(file data.MyLSFiles) string {
// 	permission := file.Mode.String()

// 	if file.IsLink {
// 		permission = strings.Replace(permission, "L", "l", 1)
// 	}
// 	if file.IsBlockDevice {
// 		permission = strings.Replace(permission, "D", "b", 1)
// 	}
// 	if file.IsCharDevice {
// 		permission = strings.Replace(permission, "D", "", 1)
// 	}
// 	if file.IsSetuid {
// 		permission = strings.Replace(permission, "u", "-", 1)
// 	}
// 	if file.IsSetgid {
// 		permission = strings.Replace(permission, "g", "-", 1)
// 	}

// 	return permission
// }

func UpdateMaxNlink(maxNlink *int, file data.MyLSFiles) {
	strNLink := strconv.Itoa(int(file.NLink))
	if *maxNlink < len(strNLink) {
		*maxNlink = len(strNLink)
	}
}

func CalculateMaxWidth(files []data.MyLSFiles) (maxOwner, maxGroup, maxSize, maxMajor, maxMinor int) {
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

func FormatTime(modTime time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0)

	if now.Sub(modTime) > 0 && modTime.After(sixMonthsAgo) {
		return modTime.Format("Jan _2 15:04")
	} else {
		return modTime.Format("Jan _2  2006")
	}
}
