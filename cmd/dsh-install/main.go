// 安装DeepSeekHarness.exe：一键安装 DeepSeek Harness——自动安装 Node.js、
// 全局安装 @deepseek-ai/dsh 并创建桌面快捷方式。
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"hsx-dsh-tools/internal/dshutil"
)

// title 是窗口标题与欢迎横幅共用的名称。
const title = "DeepSeek Harness 一键安装"

func main() {
	if dshutil.ShowVersion(strings.TrimSuffix(dshutil.ExeInstall, ".exe"), dshutil.Version) {
		return
	}
	dshutil.SetupConsole()
	dshutil.SetTitle(title)

	// Node.js 官方 MSI 是每机器安装，写入 C:\Program Files，必须以管理员权限
	// 运行。未提权时先请求一次 UAC 授权，以管理员身份重启整个流程，否则
	// msiexec /qn 会在无 UAC 的情况下直接失败（退出码 1603）。
	if !dshutil.Elevated() {
		dshutil.Info("本程序需要管理员权限来安装 Node.js，正在请求授权...")
		if err := dshutil.RelaunchElevated(); err != nil {
			dshutil.Fail("%s", err)
			dshutil.PressAnyKey()
			os.Exit(1)
		}
		return // 提权后的新实例已在独立窗口运行，当前实例退出。
	}

	if err := run(); err != nil {
		dshutil.Fail("%s", err)
		dshutil.PressAnyKey()
		os.Exit(1)
	}
	// 打开当前文件夹，方便用户查看。
	openInExplorer(exeDirectory())
	dshutil.PressAnyKey()
}

// run 依次完成 Node.js 检测安装、dsh 安装与桌面快捷方式创建。
func run() error {
	printHeader()
	if err := installNode(); err != nil {
		return err
	}
	if err := installDsh(); err != nil {
		return err
	}
	createShortcut()
	printSummary()
	return nil
}

func printHeader() {
	dshutil.Banner(title)
	dshutil.Info("")
	dshutil.Info("  [功能简介]")
	dshutil.Info("  1. 检测并安装 Node.js（如缺失自动下载）")
	dshutil.Info("  2. 一键安装 %s", dshutil.DshPackage)
	dshutil.Info("  3. 创建桌面快捷方式")
	dshutil.Info("")
	dshutil.Info("  [注意事项]")
	dshutil.Info("  安装全程需要联网，请确保网络畅通")
	dshutil.Info("  首次运行会请求管理员权限（UAC 提示），请点击\"是\"")
	dshutil.Rule()
	dshutil.Info("")
}

// installNode 确保系统中有可用的 node，失败时返回错误。
func installNode() error {
	dshutil.Step(1, 3, "检测 Node.js 环境")
	dshutil.Info("")
	if dshutil.NodeInstalled() {
		dshutil.Success("Node.js 已安装，版本：%s", dshutil.NodeVersionString())
		return nil
	}
	dshutil.Warn("未检测到 Node.js，开始自动安装...")
	if err := dshutil.InstallNode(); err != nil {
		dshutil.Info("请检查网络连接后重新运行，或访问 %s 手动下载安装", dshutil.NodeMirrorBase)
		return err
	}
	dshutil.Success("Node.js 安装完成，版本：%s", dshutil.NodeVersionString())
	return nil
}

// installDsh 全局安装 @deepseek-ai/dsh，失败时返回错误。
func installDsh() error {
	dshutil.Step(2, 3, "检测 / 安装 %s", dshutil.DshPackage)
	dshutil.Info("")
	dshutil.Info("正在切换到国内镜像源...")
	if code := dshutil.SetNpmRegistry(); code != 0 {
		dshutil.Warn("设置 npm 镜像源失败，将使用默认源")
	}
	dshutil.Info("正在安装最新版（首次安装可能需要 1-2 分钟）...")
	if code := dshutil.InstallDshLatest(); code != 0 {
		return fmt.Errorf("安装 %s 失败（退出码 %d）", dshutil.DshPackage, code)
	}
	dshutil.Success("%s 安装完成", dshutil.DshPackage)
	if dshutil.DshInstalled() {
		dshutil.Success("dsh 命令可用，版本：%s", dshutil.DshVersionString())
	} else {
		dshutil.Warn("dsh 命令未加入 PATH，请重新打开终端或重启电脑")
	}
	dshutil.Info("")
	return nil
}

// createShortcut 创建桌面快捷方式，失败仅提示不中断流程。
func createShortcut() {
	dshutil.Step(3, 3, "创建桌面快捷方式")
	dshutil.Info("")
	startExe := filepath.Join(exeDirectory(), dshutil.ExeStart)
	if err := dshutil.CreateDesktopShortcut(startExe); err != nil {
		dshutil.Warn("%s", err)
		dshutil.Info("可手动双击 %s 使用", startExe)
		return
	}
	dshutil.Success("桌面快捷方式创建成功")
	dshutil.Info("")
}

func printSummary() {
	dshutil.Banner("安装完成！")
	dshutil.Info("")
	dshutil.Info("  安装成功，现在可以使用 DeepSeek Harness")
	dshutil.Info("  推荐双击桌面上的「%s」快捷方式", dshutil.ShortcutName)
	dshutil.Info("  如需升级，请双击 %s", dshutil.ExeUpgrade)
	dshutil.Rule()
	dshutil.Info("")
}

// openInExplorer 在资源管理器中打开 dir，并立即释放进程句柄（explorer
// 独立运行，与本工具无关）。
func openInExplorer(dir string) {
	cmd := exec.Command("explorer", dir)
	if err := cmd.Start(); err == nil {
		_ = cmd.Process.Release()
	}
}

// exeDirectory 返回当前可执行文件所在目录。
func exeDirectory() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}
