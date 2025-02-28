package readinput

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type AutoCompleteFunc func(text string) string

type Reader struct {
	stdin *os.File
}

func setRawMode(fd int) error {
	var termios syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0); err != 0 {
		return err
	}

	// Отключаем канонический режим и эхо
	termios.Lflag &^= syscall.ICANON | syscall.ECHO

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0); err != 0 {
		return err
	}
	return nil
}

func restoreMode(fd int) error {
	var termios syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0); err != 0 {
		return err
	}

	// Восстанавливаем канонический режим и эхо
	termios.Lflag |= syscall.ICANON | syscall.ECHO

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0); err != 0 {
		return err
	}
	return nil
}

func NewReader() *Reader {
	fd := int(os.Stdin.Fd())
	if err := setRawMode(fd); err != nil {
		panic(err)
	}
	return &Reader{stdin: os.Stdin}
}

func (r *Reader) Close() {
	fd := int(r.stdin.Fd())
	if err := restoreMode(fd); err != nil {
		panic(err)
	}
}

func (r *Reader) ReadLine(autoCompleteFunc AutoCompleteFunc) (string, error) {
	var buffer []byte
	fmt.Fprint(os.Stdout, "$ ")
	os.Stdout.Sync()

	tmp := make([]byte, 1)
	for {
		_, err := r.stdin.Read(tmp)
		if err != nil {
			return "", err
		}
		b := tmp[0]

		switch b {
		case '\n', '\r':
			fmt.Println()
			return string(buffer), nil
		case '\t':
			newBuffer := autoCompleteFunc(string(buffer))
			if newBuffer != "" {
				buffer = []byte(newBuffer)
			}
		case 127: // Backspace
			if len(buffer) > 0 {
				buffer = buffer[:len(buffer)-1]
			}
		case 3: // Ctrl+C
			fmt.Println()
			return "", nil
		case 4: // Ctrl+D
			if len(buffer) == 0 {
				fmt.Println()
				return "", nil
			}
		default:
			if b >= 32 { // печатаемые символы
				buffer = append(buffer, b)
			}
		}

		// Обновляем экран после обработки символа
		fmt.Printf("\r\033[K$ %s", string(buffer))
		os.Stdout.Sync()
	}
}
