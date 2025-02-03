# Safe Ollama

一个为本地部署的 Ollama 提供安全鉴权接口、用户管理和用量记录的工具。方便用户将本地模型分享给他人使用。

## 项目亮点

- **一键运行**：仅需一个可执行文件和配置文件，无需安装任何额外环境。
- **完整功能**：支持鉴权、用户管理、用量记录等功能。
- **现代化技术栈**：基于 Go 和 Gin 框架构建后端，前端使用 React + Vite + shadcn/ui。

## 技术栈

- 后端：Go & Gin Framework
- 前端：React + Vite + shadcn/ui
- 数据库：SQLite

## 安装与运行

1. 下载可执行文件

2. 配置文件：

将以下配置文件放置在与可执行文件相同的目录下，命名为`config.yml`：

```yaml
server:
  address: "localhost:8080"
  logging:
    level: "info"
    output: "stdout"
  jwtkey: "secret" # 请更换为安全的密钥

admin:
  username: "admin"
  password: "admin" # 请更换为安全的密钥

ollama:
  url: "http://localhost:11434"
  timeout: 300 # 秒

database:
  url: "safe_ollama.db"
```

3. 运行程序：

```bash
./safe-ollama
```

## 配置说明

- `server`
    - `address`：服务监听的地址，格式为 "host:port"。
    - `logging`：
        - `level`：日志级别，可选值：debug, info, warning, error, fatal.
        - `output`：日志输出位置，可选值：file_path, stdout.
    - `jwtkey`：JWT 加密密钥，建议生产环境下使用强密码。
- `admin`
    - `username`：默认管理员用户名。
    - `password`：默认管理员密码，建议在生产环境中及时更改。
- `ollama`
    - `url`：Ollama 服务地址。
    - `timeout`：请求超时时间，单位秒。
- `database`
    - `url`：SQLite 数据库文件路径。

## 构建指南

> 注意：本项目使用 go embed 将前端资源打包进可执行文件中，因此需要先构建前端。

1. 前端打包：

```bash
cd safe-ollama/safe-ollama-ui && yarn install && yarn run build
```

2. 后端编译（针对不同平台）：

```bash
# 编译为 Linux 可执行文件
GOOS=linux GOARCH=amd64 go build -o safe-ollama

# 编译为 macOS 可执行文件
GOOS=darwin GOARCH=amd64 go build -o safe-ollama

# 编译为 Windows 可执行文件
GOOS=windows GOARCH=amd64 go build -o safe-ollama.exe
```

---

## 许可证

本项目使用 MIT 许可证。请在你的任何修改或分发中包含原始许可证信息。
