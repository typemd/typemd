---
title: "TypeMD 0.7.0 — 什麼都能直接改"
description: "v0.7.0 讓你在 TUI 裡直接編輯所有屬性——不用離開終端機，所有資料隨手可改。"
date: 2026-03-28
tags: [發布, 開發日誌]
---

0.6.0 讓 AI 走進你的 vault。但就算 AI 幫你產生了描述和標籤，你還是得開編輯器才能改屬性值。這次我們把最後一道牆拆了——**在 TUI 裡直接編輯所有東西**。

## 屬性？按 Enter 就能改

選中一個屬性，按 Enter，直接改。不用開檔案、不用記 YAML 格式。

所有屬性型別都支援：

- **string / text** — 行內文字輸入
- **number** — 數字輸入
- **select** — 選項清單
- **checkbox** — 切換開關
- **url** — 網址輸入
- **date / datetime** — 分段輸入搭配日曆選取器

改完按 Enter 儲存，Esc 取消。就這麼簡單。

## 日期有了專屬的日曆選取器

日期屬性不只是文字輸入——我們做了一個完整的日期選取器：

- 分段輸入（年/月/日），Tab 在各段之間移動
- 彈出行內日曆，用方向鍵瀏覽
- 支援 `date_format` 和 `datetime_format` 設定，自訂顯示格式：

```yaml
# .typemd/config.yaml
date_format: "2006/01/02"
datetime_format: "2006/01/02 15:04"
```

## 關聯選取：模糊搜尋就對了

編輯關聯屬性時，彈出模糊搜尋選取器。輸入幾個字，找到目標物件，Enter 選取。多值關聯可以重複加入。

把 book 的 author 關聯到 person？兩秒搞定。

## 表格視圖也能改

上一版加入了表格視圖，這版讓它可以編輯。在表格中選中一個儲存格，按 Enter，用跟屬性面板一模一樣的編輯元件改值。

批次檢視和修改資料，表格視圖現在是真正的工作介面。

## 鎖定物件，防手殘

重要的物件不想被不小心改到？在 frontmatter 加一行：

```yaml
locked: true
```

TUI 會顯示鎖頭圖示，所有編輯操作都會被擋下。解鎖只要把 `locked` 改回 `false`。

## 本地屬性不再消失

如果你在 frontmatter 裡加了 schema 沒定義的屬性，之前同步時它們會被靜默忽略。現在這些「本地屬性」會顯示在一個獨立的區段，並且在同步時完整保留。

自由加欄位，不用改 schema。

## 用自己的 LLM

0.6.0 的 AI 功能綁定 Claude CLI。現在你可以接任何 OpenAI 相容的 provider——Ollama、LM Studio、vLLM 都行：

```yaml
# .typemd/config.yaml
ai:
  enabled: true
  default: ollama
  providers:
    ollama:
      type: openai-compatible
      base_url: http://localhost:11434
      model: llama3.2
```

想切換 provider？改一行 `default` 就好。既有的 `ai.enabled` / `ai.model` 設定會自動遷移。

## 還有什麼

- **`tmd serve`** — 用 `tmd serve` 啟動 HTTP 伺服器，在瀏覽器裡看你的 vault，內建三色主題（warm/dark/light）
- **結構化日誌** — 全套件改用 `slog`，TUI 日誌寫到 `.typemd/logs/`，CLI 用 `--debug` 輸出
- **專注模式** — 按 `.` 切換單欄全寬 body，專心寫作

## 下一步

編輯體驗到位了。接下來 0.8.0 我們要打磨開發者體驗——template CLI、MCP 擴充、markdown syntax highlighting。

升級方式：

```bash
brew upgrade typemd/tap/typemd-cli
```

或從 [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.7.0) 下載預編譯的執行檔。
