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

設定檔使用 YAML 格式，有三個頂層命名空間：

```yaml
cli:
  default_type: page
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

檔案監聽器會監控 `objects/` 和 `.typemd/types/` 的變更。debounce 間隔決定 TUI 對檔案修改的反應速度——較低的值讓更新感覺更即時，較高的值減少多餘的刷新。

## AI 設定

AI 功能需要安裝並認證 [Claude CLI](https://docs.anthropic.com/en/docs/claude-code)。

| Key | 型別 | 預設值 | 說明 |
|-----|------|--------|------|
| `ai.enabled` | bool | `false` | 啟用 TUI 中的 AI 功能 |
| `ai.model` | string | *(claude 預設)* | 覆蓋 Claude 模型（例如 `claude-haiku-4-5-20251001`） |
| `ai.prompts.describe` | string | *(內建)* | 自訂描述產生的系統提示 |
| `ai.prompts.tag` | string | *(內建)* | 自訂標籤建議的系統提示 |
| `ai.prompts.explore` | string | *(內建)* | 自訂 schema 探索的系統提示 |
| `ai.explore.sample_count` | int | `10` | Schema 探索取樣的 object 數量 |
| `ai.explore.body_truncate` | int | `500` | 探索提示中包含的 object 正文最大字元數 |

當 `ai.enabled` 為 `true` 時，TUI 會新增兩個 AI 快捷鍵：

- **`g`** 在 object 詳情頁 — 開啟 AI 動作選擇器（產生描述 / 建議標籤）
- **`ctrl+e`** 在側邊欄 — 進入 schema 探索模式

## 讀取與寫入設定

你可以直接編輯 `.typemd/config.yaml`，或使用 CLI：

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
- AI 功能預設關閉
- 所有 AI 操作使用內建提示
- Schema 探索取樣 10 個 object，正文上限 500 字元

## 相關頁面

- [tmd config](/zh-tw/cli/config/) — CLI 命令參考
- [檔案結構](/zh-tw/advanced/file-structure/) — `.typemd/config.yaml` 的存放位置
- [tmd init](/zh-tw/cli/init/) — 選取 starter types 時建立初始 config
