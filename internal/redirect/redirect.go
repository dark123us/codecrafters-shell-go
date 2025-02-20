package redirect

import (
	"os"
)

func RedirectFile(file string, data []byte) {
	os.WriteFile(file, data, 0644)

}
