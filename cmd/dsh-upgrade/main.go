// 升级DeepSeekHarness.exe：一键将 @deepseek-ai/dsh 升级到 npm 最新版。
package main

import (
	"fmt"
	"os"
	"strings"

	"hsx-dsh-tools/internal/dshutil"
)

// title 是窗口标题与欢迎横幅共用的名称。
const title = "DeepSeek Harness 升级程序"

func main() {
	if dshutil.ShowVersion(strings.TrimSuffix(dshutil.ExeUpgrade, ".exe"), dshutil.Version) {
		return
	}
	dshutil.SetupConsole()
	dshutil.SetTitle(title)

	dshutil.Banner(title)
	dshutil.Info("将升级到 npm 仓库最新版")
	dshutil.Info("")
	dshutil.Info("  [注意] 升级全程需要联网")
	dshutil.Rule()
	dshutil.Info("")

	if err := run(); err != nil {
		dshutil.Fail("%s", err)
		dshutil.PressAnyKey()
		os.Exit(1)
	}

	dshutil.Info("")
	dshutil.Banner("升级完成！")
	dshutil.PressAnyKey()
}

// run 切换镜像源并升级 dsh，失败时返回错误。
func run() error {
	dshutil.Info("正在升级 %s 到最新版...", dshutil.DshPackage)
	dshutil.Info("")
	if code := dshutil.SetNpmRegistry(); code != 0 {
		dshutil.Warn("设置 npm 镜像源失败，将使用默认源")
	}
	if code := dshutil.UpgradeDsh(); code != 0 {
		return fmt.Errorf("升级失败（退出码 %d），请检查网络连接后重试", code)
	}
	dshutil.Info("")
	if dshutil.DshInstalled() {
		dshutil.Success("升级成功，当前版本：%s", dshutil.DshVersionString())
	} else {
		dshutil.Warn("升级后 dsh 命令不可用，请重新运行 %s", dshutil.ExeInstall)
	}
	return nil
}
