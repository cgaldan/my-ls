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

	test := WidthOfEachColumn(rows, columns, names)

	for row := range rows {
		rowLength := 0
		for column := range columns {
			index := column*rows + row
			if index < len(names) {
				rowLength += test[column]
			}
			if column != columns-1 {
				rowLength += 2
			}
		}

		if rowLength > termWidth {
			return true
		}
	}
	return false
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
