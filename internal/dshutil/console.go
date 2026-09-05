// 包 dshutil 为 DeepSeek Harness 的 Windows 三件套（安装 / 启动 / 升级）
// 提供共享的基础能力：控制台初始化与彩色输出、外部命令执行、注册表 PATH
// 处理、Node.js 安装、npm 与 dsh 操作、桌面快捷方式创建。
package dshutil

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

// procSetConsoleTitleW 封装 kernel32.SetConsoleTitleW（x/sys 未直接导出）。
var procSetConsoleTitleW = windows.NewLazySystemDLL("kernel32.dll").NewProc("SetConsoleTitleW")

// cpUTF8 是能正确显示 UTF-8 输出的 Windows 控制台代码页。
const cpUTF8 = 65001

// 仅当控制台支持时才输出的 ANSI 转义序列。
const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiRed    = "\x1b[31m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
)

// ruleLine 是 Banner 与 Rule 使用的分隔线。
const ruleLine = "============================================================"

// Stdout 与 Stderr 是控制台助手使用的输出流，便于在测试中替换。
var (
	Stdout = os.Stdout
	Stderr = os.Stderr
)

// colorEnabled 记录是否允许输出 ANSI 颜色码。
var colorEnabled bool

// SetupConsole 将控制台切换为 UTF-8 以正确显示中文，并在支持时启用 ANSI
// 颜色。stdout 被重定向或未连接真实控制台时调用是安全的。
func SetupConsole() {
	h, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil || h == windows.InvalidHandle {
		return
	}
	var mode uint32
	if windows.GetConsoleMode(h, &mode) != nil {
		return
	}
	_ = windows.SetConsoleOutputCP(cpUTF8)
	_ = windows.SetConsoleCP(cpUTF8)
	if windows.SetConsoleMode(h, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING) == nil {
		colorEnabled = true
	}
}

// SetTitle 设置控制台窗口标题。
func SetTitle(title string) {
	p, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return
	}
	_, _, _ = procSetConsoleTitleW.Call(uintptr(unsafe.Pointer(p)))
}

// ShowVersion 在命令行出现 --version 或 -v 时打印程序名与版本并返回 true。
func ShowVersion(program, version string) bool {
	for _, a := range os.Args[1:] {
		if a == "--version" || a == "-v" {
			fmt.Fprintf(Stdout, "%s v%s\n", program, version)
			return true
		}
	}
	return false
}

func ansi(code string) string {
	if colorEnabled {
		return code
	}
	return ""
}

// Info 输出一行普通信息。
func Info(format string, a ...any) {
	fmt.Fprintf(Stdout, format+"\n", a...)
}

// Success 输出绿色 [OK] 行。
func Success(format string, a ...any) {
	fmt.Fprintf(Stdout, "%s[OK]%s %s\n", ansi(ansiGreen), ansi(ansiReset), fmt.Sprintf(format, a...))
}

// Warn 输出黄色 [警告] 行。
func Warn(format string, a ...any) {
	fmt.Fprintf(Stdout, "%s[警告]%s %s\n", ansi(ansiYellow), ansi(ansiReset), fmt.Sprintf(format, a...))
}

// Fail 输出红色 [失败] 行。
func Fail(format string, a ...any) {
	fmt.Fprintf(Stdout, "%s[失败]%s %s\n", ansi(ansiRed), ansi(ansiReset), fmt.Sprintf(format, a...))
}

// Step 输出带编号的进度行，如 "[1/3] 检测 Node.js 环境..."。
func Step(n, total int, format string, a ...any) {
	fmt.Fprintf(Stdout, "%s[%d/%d]%s %s\n", ansi(ansiCyan), n, total, ansi(ansiReset), fmt.Sprintf(format, a...))
}

// Banner 输出带标题的方框。
func Banner(title string) {
	Info(ruleLine)
	Info("%s  %s%s", ansi(ansiBold), title, ansi(ansiReset))
	Info(ruleLine)
}

// Rule 输出一条分隔线。
func Rule() {
	Info(ruleLine)
}

// PressAnyKey 等待按下任意键，等价于 cmd 的 pause。stdin 非交互式控制台时
// 也能优雅降级。
func PressAnyKey() {
	fmt.Fprint(Stdout, "请按任意键继续...")
	_ = readKey()
	fmt.Fprintln(Stdout)
}

// readKey 在可能的情况下读取一次按键（无需回车）。
func readKey() bool {
	h, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil || h == windows.InvalidHandle {
		return false
	}
	var mode uint32
	if windows.GetConsoleMode(h, &mode) != nil {
		// 非交互式控制台：改为阻塞读取单个字节。
		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)
		return err == nil
	}
	raw := mode &^ (windows.ENABLE_LINE_INPUT | windows.ENABLE_ECHO_INPUT)
	_ = windows.SetConsoleMode(h, raw)
	defer windows.SetConsoleMode(h, mode)
	buf := make([]byte, 1)
	_, err = os.Stdin.Read(buf)
	return err == nil
}
