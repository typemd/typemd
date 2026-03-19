---
title: tmd config
description: 從 CLI 管理 vault 設定。
sidebar:
  order: 12
---

從命令列管理 `.typemd/config.yaml`。提供 `get`、`set`、`list` 子命令，讓你不需直接編輯檔案就能讀寫 vault 設定。

## 子命令

### `tmd config set`

設定一個 config 值。如果 `.typemd/config.yaml` 不存在，會自動建立。

```bash
tmd config set cli.default_type idea
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
```

## 已知的 Keys

| Key | 說明 |
|-----|------|
| `cli.default_type` | `tmd object create` 未指定 type 時的預設 type |
| `tui.debounce_ms` | 檔案監聽的 debounce 間隔（毫秒，預設：200） |

Keys 使用 dot-notation，對應到 YAML 巢狀結構。例如 `cli.default_type` 對應：

```yaml
cli:
  default_type: idea
tui:
  debounce_ms: 200
```

## 錯誤處理

| 情境 | 行為 |
|------|------|
| `set` 或 `get` 未知的 key | 錯誤訊息並列出已知的 keys |
| `get` 未設定的 key | 空輸出，exit code 0 |
| `list` 沒有 config 檔 | 空輸出，exit code 0 |

## 另見

- [檔案結構](/zh-tw/advanced/file-structure/) — `.typemd/config.yaml` 的存放位置
- [tmd init](/zh-tw/cli/init/) — 選取 starter types 時建立初始 config
- [tmd object create](/zh-tw/cli/create/) — 使用 `cli.default_type` 作為預設 type
