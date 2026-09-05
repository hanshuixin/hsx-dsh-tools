package dshutil

import "os/exec"

// DshInstalled 报告 dsh 命令是否可在 PATH 中被解析。
func DshInstalled() bool {
	_, err := exec.LookPath("dsh")
	return err == nil
}

// DshVersionString 返回已安装的 dsh 版本号；命令不可用时返回空字符串。
func DshVersionString() string {
	out, code := OutputCmd("dsh", "--version")
	if code != 0 {
		return ""
	}
	return out
}

// StartWeb 启动 `dsh web` 并将输出实时转发，返回其退出码。直接启动失败时
// 自动改用 npx 安装并启动最新版。
func StartWeb() int {
	if code := RunCmd("dsh", "web"); code == 0 {
		return 0
	}
	Warn("直接启动失败，尝试备选方案...")
	return RunCmd("npx", "--yes", DshPackage+"@latest", "web")
}
