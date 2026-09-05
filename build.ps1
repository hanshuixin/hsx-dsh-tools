# 构建脚本：生成三个 exe 到 dist/ 并打包 zip 整合包。
# 用法：powershell -ExecutionPolicy Bypass -File build.ps1
$ErrorActionPreference = "Stop"
Set-Location -Path $PSScriptRoot

# 发行版本号：唯一来源是 internal\dshutil\VERSION 文件（Go 侧用 go:embed
# 嵌入同一文件），打包时此处读取，避免版本号在多处重复维护。
#
# 发行包命名规则：hsx-dsh-tools-<系统>-<架构>-<版本号>.zip
# 例如 Windows x64 v0.1.0 => hsx-dsh-tools-windows-x64-0.1.0.zip
# （文件名中的版本号与各 exe 的 --version 输出保持一致。）
$version = (Get-Content (Join-Path $PSScriptRoot "internal\dshutil\VERSION") -Raw).Trim()

$dist = Join-Path $PSScriptRoot "dist"
New-Item -ItemType Directory -Force -Path $dist | Out-Null

Write-Host "[1/3] 构建 安装DeepSeekHarness.exe ..."
go build -trimpath -ldflags="-s -w" -o (Join-Path $dist "安装DeepSeekHarness.exe") .\cmd\dsh-install
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host "[2/3] 构建 启动DeepSeekHarness.exe ..."
go build -trimpath -ldflags="-s -w" -o (Join-Path $dist "启动DeepSeekHarness.exe") .\cmd\dsh-start
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host "[3/3] 构建 升级DeepSeekHarness.exe ..."
go build -trimpath -ldflags="-s -w" -o (Join-Path $dist "升级DeepSeekHarness.exe") .\cmd\dsh-upgrade
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host "复制许可证 ..."
Copy-Item (Join-Path $PSScriptRoot "LICENSE") (Join-Path $dist "LICENSE") -Force

Write-Host "打包 zip ..."
# 清理可能残留的旧文件与崩溃转储，确保 zip 只包含三个 exe 与许可证。
Get-ChildItem $dist -Filter "DeepSeek_logo.ico" -ErrorAction SilentlyContinue | Remove-Item -Force
Get-ChildItem $dist -Filter *.stackdump -ErrorAction SilentlyContinue | Remove-Item -Force
$staleReadme = Join-Path $dist "README.md"
if (Test-Path $staleReadme) { Remove-Item $staleReadme -Force }
$staleAssets = Join-Path $dist "assets"
if (Test-Path $staleAssets) { Remove-Item $staleAssets -Recurse -Force }

$zip = Join-Path $dist "hsx-dsh-tools-windows-x64-$version.zip"
if (Test-Path $zip) { Remove-Item $zip -Force }
# 发行包只含三个 exe 与 LICENSE。
$zipItems = @(
	(Join-Path $dist "安装DeepSeekHarness.exe")
	(Join-Path $dist "启动DeepSeekHarness.exe")
	(Join-Path $dist "升级DeepSeekHarness.exe")
	(Join-Path $dist "LICENSE")
)
Compress-Archive -Path $zipItems -DestinationPath $zip

Write-Host ""
Write-Host "构建完成！产物位于 dist\ 目录"
