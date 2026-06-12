# MyQQBot

基于 Go + Wails v3 的桌面 QQ 机器人，连接 OneBot v11 协议适配器（NapCatQQ / Lagrange.OneBot），集成 OpenAI 兼容 API 与 Function Calling 搜索工具。

## 功能

- 🖥️ 现代化 Dark Mode (OLED) 桌面 GUI
- 🤖 OpenAI 兼容 API 对话（支持中转 / OneAPI）
- 🔍 Function Calling 自动触发 `search_web` 搜索实时信息
- 💬 QQ 群聊 / 私聊自动回复
- ⚙️ GUI 内配置 LLM / OneBot / 搜索参数
- 📊 Dashboard、Chat、Settings、Logs 多页面

## 技术栈

- 后端：Go 1.25 + Wails v3
- 前端：React 18 + TypeScript + Vite + Tailwind CSS
- 协议：OneBot v11 (WebSocket)
- LLM：OpenAI 兼容 API
- 搜索：DuckDuckGo（默认免费）/ Tavily

## 前置依赖

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Wails v3 CLI](https://v3.wails.io/)
- [NapCatQQ](https://github.com/NapNeko/NapCatQQ) 或 [Lagrange.OneBot](https://github.com/LagrangeDev/Lagrange.OneBot)
- 有效的 OpenAI 兼容 API Key

## 快速开始

```bash
git clone https://github.com/dragonbaba/MyQQBot.git
cd MyQQBot
wails3 dev
```

首次启动后：

1. 在 Settings 页面填写 LLM Base URL、API Key、Model。
2. 填写 OneBot WebSocket URL（如 `ws://127.0.0.1:3001`）。
3. 启动 NapCatQQ / Lagrange.OneBot 并确认 OneBot HTTP API 在 `127.0.0.1:3000` 可用。
4. 返回 Settings → 启动机器人。

## 配置

应用启动时会读取工作目录下的 `config.yaml`。首次运行时文件不存在会自动生成默认值。示例模板见 `config.example.yaml`。

> ⚠️ `config.yaml` 包含 API Key，已被 `.gitignore` 排除，请勿提交到仓库。

## 截图

（待补充）

## License

MIT
