package dshutil

// SetNpmRegistry 将用户级 npm 镜像源切换为国内镜像，返回命令退出码。
func SetNpmRegistry() int {
	return RunCmd("npm", "config", "set", "registry", NpmRegistry)
}

// InstallDshLatest 全局安装 @deepseek-ai/dsh 命令行工具，返回命令退出码。
func InstallDshLatest() int {
	return RunCmd("npm", "install", "-g", DshPackage+"@latest")
}

// UpgradeDsh 将已全局安装的 @deepseek-ai/dsh 升级到最新版，返回命令退出码。
func UpgradeDsh() int {
	return RunCmd("npm", "update", "-g", DshPackage)
}
