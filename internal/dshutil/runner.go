package dshutil

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// RunCmd 执行 name 与 args，将标准输出 / 标准错误实时转发到控制台，并返回
// 进程退出码；命令完全无法启动时返回 -1。Go >= 1.20 会自动通过 cmd /c
// 运行 .cmd / .bat 垫片，因此 npm、dsh 无需额外包装。
func RunCmd(name string, args ...string) int {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	fmt.Fprintf(Stderr, "[内部错误] 无法启动命令 %s：%v\n", name, err)
	return -1
}

// OutputCmd 运行 name 与 args，捕获标准输出，并连同退出码一起返回（输出已
// 去除首尾空白）。
func OutputCmd(name string, args ...string) (string, int) {
	return OutputCmdTimeout(0, name, args...)
}

// OutputCmdTimeout 与 OutputCmd 类似，但可用 timeout 限制执行时长
// （0 表示不限制）。超时的命令返回 -1 退出码。
func OutputCmdTimeout(timeout time.Duration, name string, args ...string) (string, int) {
	// 统一走 CommandContext：timeout 为 0 时等价于无超时的普通命令。
	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out)), 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return strings.TrimSpace(string(out)), ee.ExitCode()
	}
	return "", -1
}
