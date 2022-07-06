//go:build windows

package winter

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _RECT struct {
	Left, Top, Right, Bottom int32
}

type Winter struct {
	mod                     *windows.LazyDLL
	procGetForegroundWindow *windows.LazyProc
	procGetWindowRect       *windows.LazyProc
}

func NewWinter() (*Winter, error) {
	w := &Winter{}
	w.mod = windows.NewLazyDLL("user32.dll")
	w.procGetForegroundWindow = w.mod.NewProc("GetForegroundWindow")
	w.procGetWindowRect = w.mod.NewProc("GetWindowRect")

	return w, nil
}

func (w *Winter) GetActiveWindowDimensions() (x, y, width, height int, err error) {
	hwnd, _, _ := w.procGetForegroundWindow.Call()

	var rect _RECT

	if ret, _, _ := syscall.Syscall(w.procGetWindowRect.Addr(), 2, hwnd, uintptr(unsafe.Pointer(&rect)), 0); ret != 0 {
		return int(rect.Left), int(rect.Top), int(rect.Right - rect.Left), int(rect.Bottom - rect.Top), nil
	}
	return 0, 0, 0, 0, fmt.Errorf("no foreground window")
}
