# 构建脚本：生成三个 exe 到 dist/ 并打包 zip 整合包。
# 用法：powershell -ExecutionPolicy Bypass -File build.ps1
$ErrorActionPreference = "Stop"
Set-Location -Path $PSScriptRoot

# 发行版本号：唯一来源是 internal\dshutil\VERSION 文件（Go 侧用 go:embed
# 嵌入同一文件），打包时此处读取，避免版本号在多处重复维护。
#
# exe 的图标 / 版本信息 / 语言由 go-winres 在构建时根据各 cmd/*/winres.json
# 模板自动生成（版本取自下方 $version），因此改版本号后无需手动重生成资源。
#
# 发行包命名规则：hsx-dsh-tools-<系统>-<架构>-<版本号>.zip
# 例如 Windows x64 v0.1.0 => hsx-dsh-tools-windows-x64-0.1.0.zip
# （文件名中的版本号与各 exe 的 --version 输出保持一致。）
$version = (Get-Content (Join-Path $PSScriptRoot "internal\dshutil\VERSION") -Raw).Trim()
# VERSIONINFO 用四段式（a.b.c.d）：VERSION=0.1.1 -> 0.1.1.0
$winVer = "$version.0"

# go-winres 用于把 winres.json 编译成 go build 所需的 .syso 资源文件。
if (-not (Get-Command go-winres -ErrorAction SilentlyContinue)) {
    throw "未找到 go-winres，请先安装：go install github.com/tc-hib/go-winres@latest，并把 GOPATH\bin 加入 PATH"
}

# Update-WindowsResources 依据 cmd\<cmdDir>\winres.json 模板重新生成该程序的
# 图标 / 版本 / 语言资源（写入 cmd\<cmdDir>\rsrc_windows_amd64.syso）。
function Update-WindowsResources {
    param([string]$cmdDir, [string]$version)

    $dir = Join-Path $PSScriptRoot "cmd\$cmdDir"
    $tpl = Join-Path $dir "winres.json"
    if (-not (Test-Path $tpl)) { throw "缺少 winres.json 模板：$tpl" }

    # 版本只出现在四段式占位符 0.0.0.0 上，替换为实际版本即可，天然幂等。
    # 必须按 UTF-8 读取：PowerShell 5.1 默认按 ANSI 解码会把中文后的引号吞掉，
    # 生成损坏的 JSON。
    $text = Get-Content $tpl -Raw -Encoding UTF8
    $text = [regex]::Replace($text, '0\.0\.0\.0', $version)

    # go-winres make 从工作目录下 winres\winres.json 读取配置。
    $winresDir = Join-Path $dir "winres"
    New-Item -ItemType Directory -Force -Path $winresDir | Out-Null
    # 用无 BOM 的 UTF-8 写回，避免 JSON 解析失败；模板本身保持占位符不被改动。
    $utf8NoBom = New-Object System.Text.UTF8Encoding($false)
    [System.IO.File]::WriteAllText((Join-Path $winresDir "winres.json"), $text, $utf8NoBom)

    Push-Location $dir
    try {
        go-winres make
        if ($LASTEXITCODE -ne 0) { throw "go-winres make 失败（cmd\$cmdDir）" }
        # 仅发布 x64，移除多余的 386 资源对象。
        Remove-Item (Join-Path $dir "rsrc_windows_386.syso") -ErrorAction SilentlyContinue
    } finally {
        Pop-Location
    }
}

$dist = Join-Path $PSScriptRoot "dist"
New-Item -ItemType Directory -Force -Path $dist | Out-Null

$apps = @(
    @{ Dir = "dsh-install"; Exe = "安装DeepSeekHarness.exe" },
    @{ Dir = "dsh-start";   Exe = "启动DeepSeekHarness.exe" },
    @{ Dir = "dsh-upgrade"; Exe = "升级DeepSeekHarness.exe" }
)

$step = 1
foreach ($app in $apps) {
    Write-Host "[$step/$($apps.Count)] 构建 $($app.Exe) ..."
    Update-WindowsResources -cmdDir $app.Dir -version $winVer
    go build -trimpath -ldflags="-s -w" -o (Join-Path $dist $app.Exe) (Join-Path ".\cmd" $app.Dir)
    if ($LASTEXITCODE -ne 0) { exit 1 }
    $step++
}

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
