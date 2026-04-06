---
title: tmd config
description: 從 CLI 管理 vault 設定。
sidebar:
  order: 23
---

從命令列管理 `.typemd/config.yaml`。提供 `get`、`set`、`list` 子命令，讓你不需直接編輯檔案就能讀寫 vault 設定。

## 子命令

### `tmd config set`

設定一個 config 值。如果 `.typemd/config.yaml` 不存在，會自動建立。

```bash
tmd config set cli.default_type idea
tmd config set ai.enabled true
```

只接受已知的 key。設定未知的 key 會回傳錯誤，並列出有效的 keys。

### `tmd config get`

取得一個 config 值。將值輸出到 stdout。如果 key 未設定，不輸出任何內容。

```bash
tmd config get cli.default_type
# → idea
```

### `tmd config list`

列出所有已設定（非空）的 config 值，格式為 `key: value`。如果沒有 config 或所有值都是空的，不輸出任何內容。

```bash
tmd config list
# cli.default_type: idea
# ai.enabled: true
```

使用 `--all` 顯示所有已知的 key，包括未設定的 key 及其預設值。

## 已知的 Keys

### CLI

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `cli.default_type` | string | *(空)* | `tmd object create` 未指定 type 時的預設 type |

### TUI

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `tui.debounce_ms` | int | `200` | 檔案監聽的 debounce 間隔（毫秒） |
| `tui.stats_type_layout` | string | `fullscreen` | 統計詳情佈局：`fullscreen` 或 `popup` |

### AI

AI 功能需要安裝並認證 `claude` CLI。

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `ai.enabled` | bool | `false` | 啟用 TUI 中的 AI 功能（`g` 和 `ctrl+e` 快捷鍵） |
| `ai.model` | string | *(claude 預設)* | 覆蓋 Claude 模型（例如 `claude-haiku-4-5-20251001`） |
| `ai.prompts.describe` | string | *(內建)* | 自訂 AI 描述產生的系統提示 |
| `ai.prompts.tag` | string | *(內建)* | 自訂 AI 標籤建議的系統提示 |
| `ai.prompts.explore` | string | *(內建)* | 自訂 AI schema 探索的系統提示 |
| `ai.explore.sample_count` | int | `10` | Schema 探索取樣的 object 數量 |
| `ai.explore.body_truncate` | int | `500` | 探索提示中包含的 object 正文最大字元數 |

Keys 使用 dot-notation，對應到 YAML 巢狀結構：

```yaml
cli:
  default_type: idea
tui:
  debounce_ms: 200
  stats_type_layout: fullscreen
ai:
  enabled: true
  model: claude-sonnet-4-6-20250627
  prompts:
    describe: "自訂描述的系統提示"
  explore:
    sample_count: 20
    body_truncate: 1000
```

## 錯誤處理

| 情境 | 行為 |
|------|------|
| `set` 或 `get` 未知的 key | 錯誤訊息並列出已知的 keys |
| `get` 未設定的 key | 空輸出，exit code 0 |
| `list` 沒有 config 檔 | 空輸出，exit code 0 |

## 另見

- [設定](/zh-tw/basics/configuration/) — 設定概念說明與完整參考
- [檔案結構](/zh-tw/advanced/file-structure/) — `.typemd/config.yaml` 的存放位置
- [tmd init](/zh-tw/cli/init/) — 選取 starter types 時建立初始 config
- [TUI AI 輔助](/zh-tw/tui/tui/#ai-assist) — 由 `ai.enabled` 啟用的 AI 功能
