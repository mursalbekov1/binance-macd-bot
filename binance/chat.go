package binance

import (
	"binance_tg/logging"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func СheckIDInFile(chatID int64) bool {
	mu.Lock()
	defer mu.Unlock()

	lines, err := ReadLines(launchDataFile)
	if err != nil {
		return false
	}

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 3 {
			if id, err := strconv.ParseInt(parts[1], 10, 64); err == nil && id == chatID {
				return true
			}
		}
	}

	return false
}

func saveLaunchDataToFile(chatID int64, command string) error {
	mu.Lock()
	defer mu.Unlock()

	lines, err := ReadLines(launchDataFile)
	if err != nil {
		return err
	}

	found := false
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 3 {
			if id, err := strconv.ParseInt(parts[1], 10, 64); err == nil && id == chatID {
				found = true
				break
			}
		}
	}

	if !found {
		data := logging.CurrentDatetime() + " " + strconv.FormatInt(chatID, 10) + " " + command
		lines = append(lines, data)

		err := writeLines(lines, launchDataFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeActiveSession(chatID int64) error {
	mu.Lock()
	defer mu.Unlock()

	// Чтение текущего содержимого файла
	lines, err := ReadLines(launchDataFile)
	if err != nil {
		return err
	}

	indexToRemove := -1
	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 3 {
			if id, err := strconv.ParseInt(parts[1], 10, 64); err == nil && id == chatID {
				indexToRemove = i
				break
			}
		}
	}

	if indexToRemove != -1 {
		lines = append(lines[:indexToRemove], lines[indexToRemove+1:]...)

		// Запись обновленных данных в файл
		err := writeLines(lines, launchDataFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadLines(filename string) ([]string, error) {
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
