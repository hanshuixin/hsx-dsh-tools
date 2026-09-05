// 启动DeepSeekHarness.exe：一键启动 DeepSeek Harness 的 Web 界面。
package main

import (
	"os"
	"strings"

	"hsx-dsh-tools/internal/dshutil"
)

// title 是窗口标题与欢迎横幅共用的名称。
const title = "DeepSeek Harness 启动器"

func main() {
	if dshutil.ShowVersion(strings.TrimSuffix(dshutil.ExeStart, ".exe"), dshutil.Version) {
		return
	}
	dshutil.SetupConsole()
	dshutil.SetTitle(title)

	dshutil.Banner(title)
	dshutil.Info("")
	dshutil.Info("  服务地址：%s", dshutil.WebURL)
	dshutil.Info("  关闭本窗口 = 停止服务")
	dshutil.Rule()
	dshutil.Info("")

	if !dshutil.DshInstalled() {
		dshutil.Warn("未检测到 %s", dshutil.DshPackage)
		dshutil.Info("")
		dshutil.Info("请先双击运行 %s 进行安装", dshutil.ExeInstall)
		dshutil.Info("")
		dshutil.PressAnyKey()
		os.Exit(1)
	}

	dshutil.Info("当前版本：%s", dshutil.DshVersionString())
	dshutil.Info("")
	dshutil.Info("正在启动 DeepSeek Harness Web 界面...")
	dshutil.Info("")

	dshutil.StartWeb()

	dshutil.Info("")
	dshutil.Banner("服务已停止")
	dshutil.PressAnyKey()
}
