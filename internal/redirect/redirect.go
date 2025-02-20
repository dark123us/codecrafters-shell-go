package redirect

import (
	"fmt"
	"os"
)

func RedirectFile(file string, data []byte) {
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи в файл: %v\n", err)
		return
	}
}
