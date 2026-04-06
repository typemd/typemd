---
title: tmd serve
description: 啟動 Web UI 伺服器，透過瀏覽器存取 vault。
sidebar:
  order: 25
---

啟動 HTTP 伺服器，提供瀏覽器介面和 REST API 來存取你的 vault。

## 使用方式

```bash
tmd serve
tmd serve -p 8080
```

在 `http://localhost:3000`（預設 port）開啟 Web UI。介面和 TUI 相同的三欄式佈局：sidebar（type 群組 + object 列表）、body（markdown 內容）和 properties 面板。

## Flags

| Flag | 縮寫 | 預設值 | 說明 |
|------|------|--------|------|
| `--port` | `-p` | `3000` | 監聽的 port |
| `--host` | | `127.0.0.1` | 綁定的 host（使用 `0.0.0.0` 監聽所有介面） |
| `--vault` | | `.` | vault 目錄路徑（全域 flag） |

## Web UI 功能

- **三欄式佈局** — sidebar、body 和 properties 面板
- **瀏覽 Object** — 展開 type 群組，點擊查看 object
- **屬性編輯** — 雙擊屬性值即可編輯
- **內文編輯** — 點擊「Edit」切換到 markdown 編輯器（⌘+Enter 儲存，Esc 取消）
- **建立 Object** — 按 `N` 或點擊 sidebar 的「+ New」

### 快捷鍵

| 按鍵 | 動作 |
|------|------|
| `.` | 切換 focus 模式（全寬 body） |
| `p` | 切換 properties 面板 |
| `n` | 開啟建立 object 對話框 |

## REST API

伺服器在 `/api/` 下提供 REST API 供程式存取。

### Types

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/types` | 列出所有 type，包含 metadata 和 object 數量 |
| GET | `/api/types/{name}` | 取得 type schema 和屬性定義 |

### Objects

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/objects` | 列出所有 object（可用 `?type=book` 篩選） |
| GET | `/api/objects/{type}/{slug}` | 取得單一 object 的屬性和正文 |
| PUT | `/api/objects/{type}/{slug}` | 更新 object 正文（`{"body": "..."}`） |
| POST | `/api/objects` | 建立 object（`{"type": "book", "name": "...", "template": "..."}`） |

### Properties

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/properties/{type}/{slug}` | 取得 object 的顯示用屬性 |
| PUT | `/api/properties/{type}/{slug}/{key}` | 更新單一屬性（`{"value": "..."}`） |

### Templates

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/templates/{type}` | 列出指定 type 的可用 template |

## 開發

前端開發搭配 hot reload：

```bash
# Terminal 1: Go 伺服器
tmd serve

# Terminal 2: Vite dev server（將 /api 代理到 localhost:3000）
cd web/frontend && npm run dev
```

Vite dev server 在 `http://localhost:5173` 提供 hot module replacement。API 請求會自動代理到 Go 伺服器。
