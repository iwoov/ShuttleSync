# 构建与运行指南

本文档说明如何在不同平台上编译并运行 ShuttleSync 项目（Vue 前端 + Go 后端）。

## 基础依赖

| 平台 | 必备工具 |
| --- | --- |
| Linux / macOS / WSL | Node.js (含 npm)、Go 1.22+、Bash、rsync(可选) |
| Windows (PowerShell) | Node.js (含 npm)、Go 1.22+、[Robocopy](https://learn.microsoft.com/windows-server/administration/windows-commands/robocopy) 或其他拷贝工具 |

> 提示：Windows 推荐在 WSL 或 Git Bash 中运行 `scripts/build_and_run.sh`，这样可复用与 Linux/macOS 相同的命令。

## Linux / macOS / WSL

### 一键脚本
```bash
chmod +x scripts/build_and_run.sh
./scripts/build_and_run.sh
```
脚本将自动：
1. 安装/更新前端依赖并执行 `npm run build`；
2. 将 `frontend/dist` 拷贝进 `backend/web`；
3. 运行 `go build -o backend/bin/shuttlesync`；
4. 启动服务器并监听 `http://127.0.0.1:3456/`。

### 手动操作
```bash
# 1. 编译前端
cd frontend
npm install
npm run build

# 2. 同步前端产物
cd ..
rsync -a --delete frontend/dist/ backend/web/
# 若无 rsync，可改用：cp -R frontend/dist/. backend/web/

# 3. 编译后端
cd backend
mkdir -p bin
go build -o bin/shuttlesync ./

# 4. 运行
./bin/shuttlesync
```

## Windows (PowerShell)
```powershell
# 1. 编译前端
cd frontend
npm install
npm run build
cd ..

# 2. 同步前端产物（示例使用 robocopy）
robocopy frontend\dist backend\web /MIR

# 3. 编译后端
cd backend
mkdir bin
go build -o bin\shuttlesync.exe .\

# 4. 运行
.\bin\shuttlesync.exe
```
> 如果使用 Git Bash/WSL，可直接执行脚本 `./scripts/build_and_run.sh`。

## 交叉编译示例
使用 Go 的交叉编译能力可生成不同平台二进制（亦可参考 `.github/workflows/build.yml`）：
```bash
# Linux 上编译 Windows 版本
go env -w GOOS=windows GOARCH=amd64
cd backend
go build -o bin/shuttlesync.exe ./

# 恢复本地配置
unset GOOS GOARCH # 或 go env -u GOOS; go env -u GOARCH
```

## 常见问题
- **缺少 npm/go 命令**：确认安装并加入 PATH。
- **rsync 未安装**：按需安装，或使用 `cp -R`/`robocopy`。
- **数据库文件**：`backend/database.db` 仅用于本地开发，生产环境应使用迁移重新生成。
