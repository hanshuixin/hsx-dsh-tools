package dshutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// procShellExecuteW 封装 shell32.ShellExecuteW（x/sys 未直接导出）。
var procShellExecuteW = windows.NewLazySystemDLL("shell32.dll").NewProc("ShellExecuteW")

// Elevated 报告当前进程是否已以管理员权限（高完整性令牌）运行。
func Elevated() bool {
	token, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return false
	}
	defer token.Close()
	return token.IsElevated()
}

// RelaunchElevated 以管理员身份重新启动当前可执行文件（会触发一次 UAC 提示）。
// 成功启动后返回 nil，调用方应立即退出当前实例；用户拒绝授权或启动失败时
// 返回可读错误。
func RelaunchElevated() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取当前程序路径：%w", err)
	}
	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString(exe)
	args, _ := windows.UTF16PtrFromString(strings.Join(os.Args[1:], " "))
	dir, _ := windows.UTF16PtrFromString(filepath.Dir(exe))

	// 返回值 >32 表示成功启动；否则是 SE_ERR_* 之类的错误码。UAC 被拒绝时
	// 系统将 last error 置为 ERROR_CANCELLED(1223)。
	r, _, callErr := procShellExecuteW.Call(
		0, // 无父窗口句柄
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(dir)),
		uintptr(windows.SW_SHOWNORMAL),
	)
	if r > 32 {
		return nil
	}
	var errno syscall.Errno
	if errors.As(callErr, &errno) && errno == windows.ERROR_CANCELLED {
		return errors.New("已取消管理员授权，操作需要管理员权限")
	}
	return fmt.Errorf("无法以管理员身份重新启动（错误码 %d），请右键选择“以管理员身份运行”", r)
}
