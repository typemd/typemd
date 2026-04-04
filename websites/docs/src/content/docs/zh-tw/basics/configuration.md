---
title: 設定
description: 如何透過 .typemd/config.yaml 設定你的 vault。
sidebar:
  order: 7
---

每個 vault 都有一個選用的設定檔 `.typemd/config.yaml`。它控制 CLI 預設值、TUI 行為和 AI 功能。當你執行 `tmd init` 或 `tmd config set` 時，檔案會自動建立。

## 檔案位置

```
my-vault/
└── .typemd/
    └── config.yaml
```

## 結構

設定檔使用 YAML 格式，有頂層顯示設定和三個命名空間：

```yaml
date_format: "YYYY-MM-DD"
datetime_format: "YYYY-MM-DD HH:mm:ss"
cli:
  default_type: page
tui:
  debounce_ms: 200
  stats_type_layout: fullscreen
  toast:
    duration_ms: 3000
    dismiss_key: esc
ai:
  enabled: true
  default: claude
  providers:
    claude:
      type: cli
      model: claude-sonnet-4-6-20250627
    ollama:
      type: openai-compatible
      base_url: http://localhost:11434
      model: qwen3-coder:30b
  prompts:
    describe: "自訂描述的系統提示"
  explore:
    sample_count: 20
    body_truncate: 1000
```

## 日期顯示格式

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `date_format` | string | `YYYY-MM-DD` | `date` 屬性的顯示格式 |
| `datetime_format` | string | `YYYY-MM-DD HH:mm:ss` | `datetime` 屬性及系統時間戳（`created_at`、`updated_at`）的顯示格式 |

這些設定控制日期在 TUI 和 CLI 中的顯示方式。儲存格式不變（ISO 8601 / RFC 3339）。

**支援的 token：** `YYYY`（年）、`MM`（月）、`DD`（日）、`HH`（24 小時制）、`mm`（分）、`ss`（秒）。無法辨識的字元會原樣保留。

```yaml
# 歐式格式
date_format: "DD/MM/YYYY"
datetime_format: "DD/MM/YYYY HH:mm:ss"

# 中文日期
date_format: "YYYY年MM月DD日"
```

## CLI 設定

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `cli.default_type` | string | *(空)* | `tmd object create` 未指定 type 時的預設 type |

設定後，建立 object 時可以省略 type：

```bash
# 沒有 default_type — 必須指定 type
tmd object create idea "My Idea"

# 設定了 cli.default_type: idea
tmd object create "My Idea"
```

## TUI 設定

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `tui.debounce_ms` | int | `200` | 檔案監聽的 debounce 間隔（毫秒） |
| `tui.stats_type_layout` | string | `fullscreen` | 統計詳情佈局：`fullscreen` 或 `popup` |
| `tui.toast.position` | string | `bottom-right` | Toast 顯示位置：`bottom-right` 或 `help-bar` |
| `tui.toast.duration_ms` | int | `3000` | 自動消失的時間（毫秒） |
| `tui.toast.dismiss_key` | string | `esc` | 手動關閉 toast 通知的按鍵 |
| `tui.toast.show_warnings` | bool | `true` | 顯示警告等級的 toast 通知 |
| `tui.toast.show_success` | bool | `false` | 顯示資訊/成功等級的 toast 通知 |
| `tui.theme.focus_border` | string | `63` | 聚焦面板的邊框顏色（ANSI 色碼或 hex） |
| `tui.theme.wiki_link` | string | `33` | Wiki-link 顯示文字顏色 |
| `tui.theme.heading` | string | `3` | Markdown 標題顏色 |
| `tui.theme.bold` | string | *(空)* | 粗體文字顏色（空值 = 僅粗體樣式） |
| `tui.theme.italic` | string | *(空)* | 斜體文字顏色（空值 = 僅斜體樣式） |
| `tui.theme.inline_code` | string | `245` | 行內程式碼顏色 |
| `tui.theme.code_block` | string | `245` | 程式碼區塊顏色 |
| `tui.theme.link` | string | `33` | Markdown 連結顏色 |
| `tui.theme.blockquote` | string | `8` | 引用區塊顏色 |
| `tui.theme.hrule` | string | `8` | 水平分隔線顏色 |

檔案監聽器會監控 `objects/` 和 `.typemd/types/` 的變更。debounce 間隔決定 TUI 對檔案修改的反應速度——較低的值讓更新感覺更即時，較高的值減少多餘的刷新。

Toast 通知顯示在 TUI 的右下角，用於暫時性訊息，例如同步警告（未解析的參照）和 AI 操作錯誤。通知會在 `duration_ms` 後自動消失，也可以按 `dismiss_key` 手動關閉。錯誤等級的 toast 通知無論設定如何都會顯示。

## AI 設定

AI 功能支援多個 provider。`cli` 類型使用 [Claude CLI](https://docs.anthropic.com/en/docs/claude-code)；`openai-compatible` 類型可搭配任何 OpenAI 相容 API（Ollama、LM Studio、vLLM、LocalAI 等）。

### Provider 設定

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `ai.enabled` | bool | `false` | 啟用 TUI 中的 AI 功能 |
| `ai.default` | string | *(空)* | 從 `ai.providers` 中選擇使用的 provider 名稱 |
| `ai.providers.<name>.type` | string | — | Provider 類型：`cli` 或 `openai-compatible` |
| `ai.providers.<name>.model` | string | *(空)* | 模型識別碼 |
| `ai.providers.<name>.base_url` | string | *(空)* | HTTP 端點（`openai-compatible` 必填） |
| `ai.providers.<name>.api_key` | string | *(空)* | 選用的 Bearer token 認證 |

多 provider 範例：

```yaml
ai:
  enabled: true
  default: ollama
  providers:
    claude:
      type: cli
      model: claude-sonnet-4-6-20250627
    ollama:
      type: openai-compatible
      base_url: http://localhost:11434
      model: qwen3-coder:30b
    my-server:
      type: openai-compatible
      base_url: http://192.168.1.100:8080
      model: llama3.2
      api_key: sk-my-key
```

只需更改 `ai.default` 即可切換 provider，無需修改其他設定。

### 提示與探索設定

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `ai.prompts.describe` | string | *(內建)* | 自訂描述產生的系統提示 |
| `ai.prompts.tag` | string | *(內建)* | 自訂標籤建議的系統提示 |
| `ai.prompts.explore` | string | *(內建)* | 自訂 schema 探索的系統提示 |
| `ai.explore.sample_count` | int | `10` | Schema 探索取樣的 object 數量 |
| `ai.explore.body_truncate` | int | `500` | 探索提示中包含的 object 正文最大字元數 |

當 `ai.enabled` 為 `true` 時，TUI 會新增兩個 AI 快捷鍵：

- **`g`** 在 object 詳情頁 — 開啟 AI 動作選擇器（產生描述 / 建議標籤）
- **`ctrl+e`** 在側邊欄 — 進入 schema 探索模式

## 讀取與寫入設定

你可以直接編輯 `.typemd/config.yaml`、使用 CLI，或使用 TUI 設定頁面（在 TUI 中按 `,` 開啟可瀏覽的設定編輯器）：

```bash
# 設定一個值
tmd config set ai.enabled true

# 取得一個值
tmd config get cli.default_type

# 列出所有已設定的值
tmd config list

# 列出所有已知的 key 及其預設值
tmd config list --all
```

只接受已知的 key。設定未知的 key 會回傳錯誤，並列出有效的 keys。

## 預設值

如果 `.typemd/config.yaml` 不存在或某個 key 未設定，TypeMD 會使用合理的預設值：

- 沒有預設 type — `tmd object create` 必須指定 type
- 檔案監聽 debounce 為 200ms
- 統計佈局使用全螢幕模式
- Toast 通知：右下角、3 秒自動消失、Esc 關閉、顯示警告、隱藏成功訊息
- AI 功能預設關閉
- 所有 AI 操作使用內建提示
- Schema 探索取樣 10 個 object，正文上限 500 字元

## 相關頁面

- [tmd config](/zh-tw/cli/config/) — CLI 命令參考
- [檔案結構](/zh-tw/advanced/file-structure/) — `.typemd/config.yaml` 的存放位置
- [tmd init](/zh-tw/cli/init/) — 選取 starter types 時建立初始 config
