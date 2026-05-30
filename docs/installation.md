# Installation Guide

This guide explains how to install EasyEdit on your system.

## English

### Requirements

- Go 1.21 or higher
- A terminal emulator
- Supported OS: Windows, Linux, or macOS

### Method 1: Go Install (Recommended)

```bash
go install github.com/0xA672/easyedit@latest
```

This installs the `easyedit` binary to your `$GOPATH/bin` (or `$HOME/go/bin`). Make sure this directory is in your `PATH`.

**Verify installation:**
```bash
easyedit --help
```

### Method 2: Build from Source

```bash
git clone https://github.com/0xA672/easyedit.git
cd easyedit
go build -o easyedit .
```

Move the binary to your PATH:
```bash
# Linux/macOS
sudo mv easyedit /usr/local/bin/

# Windows (PowerShell)
move easyedit.exe C:\Windows\System32\
```

### Method 3: Download Prebuilt Binary

> Coming soon — releases page will host prebuilt binaries for Windows, Linux, and macOS.

### Troubleshooting

**Issue**: `command not found: easyedit`  
**Solution**: Add Go bin directory to your PATH:
```bash
# Linux/macOS
export PATH=$PATH:$HOME/go/bin

# Windows (PowerShell)
$env:Path += ";$HOME\go\bin"
```

---

## 中文

### 系统要求

- Go 1.21 或更高版本
- 终端模拟器
- 支持的操作系统：Windows、Linux 或 macOS

### 方法一：Go 安装（推荐）

```bash
go install github.com/0xA672/easyedit@latest
```

这会将 `easyedit` 二进制文件安装到您的 `$GOPATH/bin`（或 `$HOME/go/bin`）。确保此目录在您的 `PATH` 中。

**验证安装：**
```bash
easyedit --help
```

### 方法二：从源码构建

```bash
git clone https://github.com/0xA672/easyedit.git
cd easyedit
go build -o easyedit .
```

将二进制文件移动到 PATH：
```bash
# Linux/macOS
sudo mv easyedit /usr/local/bin/

# Windows (PowerShell)
move easyedit.exe C:\Windows\System32\
```

### 方法三：下载预编译二进制

> 即将推出 — releases 页面将提供 Windows、Linux 和 macOS 的预编译二进制文件。

### 故障排除

**问题**：`command not found: easyedit`  
**解决方案**：将 Go bin 目录添加到 PATH：
```bash
# Linux/macOS
export PATH=$PATH:$HOME/go/bin

# Windows (PowerShell)
$env:Path += ";$HOME\go\bin"
```
