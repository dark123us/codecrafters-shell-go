package redirect

import (
	"fmt"
	"os"
)

func RedirectFile(filename string, data []byte) {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи в файл: %v\n", err)
		return
	}
}

func RedirectFileAppend(filename string, data []byte) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи в файл: %v\n", err)
		return
	}
}
