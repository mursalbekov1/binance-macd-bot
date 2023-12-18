package binance

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func saveLaunchDataToFile(chatID int64, command string) error {
	file, err := os.OpenFile(launchDataFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := strconv.FormatInt(chatID, 10) + " " + command + "\n"
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func removeActiveSession(chatID int64) error {
	lines, err := readLines(launchDataFile)
	if err != nil {
		return err
	}

	// Ищем индекс строки, соответствующей chatID
	indexToRemove := -1
	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 {
			if id, err := strconv.ParseInt(parts[0], 10, 64); err == nil && id == chatID {
				indexToRemove = i
				break
			}
		}
	}

	// Если нашли строку, удаляем ее из среза
	if indexToRemove != -1 {
		lines = append(lines[:indexToRemove], lines[indexToRemove+1:]...)
	}

	// Записываем обновленные данные обратно в файл
	err = writeLines(lines, launchDataFile)
	if err != nil {
		return err
	}

	return nil
}

// readLines читает строки из файла в срез строк
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// writeLines записывает срез строк в файл
func writeLines(lines []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	return writer.Flush()
}
