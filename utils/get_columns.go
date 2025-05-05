package utils

func GetColumns(numFiles, width int, filesName []string) int {
	columns := numFiles

	for IsRowBiggerThanTermWidth(columns, width, filesName) {
		columns--
	}

	return columns
}

func IsRowBiggerThanTermWidth(columns, termWidth int, names []string) bool {
	rows := 1

	if columns > 0 {
		rows = (len(names) + columns - 1) / columns
	}

	widthMap := WidthOfEachColumn(rows, columns, names)

	rowLength := 0
	mapLen := len(widthMap)
	for i := range mapLen {
		rowLength += widthMap[i]
		if i != mapLen-1 {
			rowLength += 2
		}
	}

	return rowLength >= termWidth
}

func WidthOfEachColumn(rows, columns int, allFileNames []string) map[int]int {
	widthOfColumns := make(map[int]int)

	for row := range rows {
		for column := range columns {
			index := column*rows + row
			if index < len(allFileNames) {
				if len(allFileNames[index]) > widthOfColumns[column] {
					widthOfColumns[column] = len(allFileNames[index])
				}
			}
		}
	}

	return widthOfColumns
}

func SortDirs(s *[]string) {
	for i := 0; i < len(*s)-1; i++ {
		for j := 0; j < len(*s)-i-1; j++ {
			if (*s)[j] > (*s)[j+1] {
				(*s)[j], (*s)[j+1] = (*s)[j+1], (*s)[j]
			}
		}
	}
}
