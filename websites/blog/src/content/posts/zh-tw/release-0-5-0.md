---
title: "TypeMD 0.5.0 — 你的資料，你的視角"
description: "v0.5.0 帶來完整的 View 系統：list 與 table 兩種佈局、篩選、排序、分組，以及全新的 tmd stats 指令。"
date: 2026-03-20
tags: [發布, 開發日誌]
---

0.4.0 結尾我們提到 aliases 和搜尋體驗在雷達上。但實際用了一陣子之後，發現「怎麼看資料」比「怎麼找資料」是更迫切的痛點——物件越來越多，側邊欄已經不夠用了。

所以 0.5.0 先做了 View 系統：**讓你用自己的方式看自己的資料**。

## View 系統：一個型別，多種看法

每個型別現在可以定義多個 View。一個 View 就是一組「我想怎麼看這些物件」的設定——用哪種佈局、篩選什麼、怎麼排序、怎麼分組。

View 存在 `types/<name>/views/<view>.yaml`，跟型別 schema 一樣是純文字檔，可以 git 追蹤。

每個型別都有一個隱含的 default view（列表佈局、按名稱排序）。想客製？在 TUI 裡按 `v` 就能進入視圖模式。

## 列表 vs. 表格

兩種佈局可以依場景切換：

- **List**（`layout: list`）— 精簡的名稱列表，適合快速瀏覽。
- **Table**（`layout: table`）— 欄位表格，把屬性攤開看。用 `columns` 指定要顯示哪些屬性：

```yaml
name: reading-list
layout: table
columns: [status, rating, author]
sort:
  - property: rating
    order: desc
```

表格模式讓你一眼看到每本書的狀態、評分和作者，不用點進去。

## 篩選：只看你在乎的

View 支援結構化的 filter rule。每種屬性型別有各自的運算子——字串可以 `contains`、`starts_with`；數字可以 `gt`、`lte`；日期可以 `before`、`after`：

```yaml
filter:
  - property: status
    operator: is
    value: reading
  - property: rating
    operator: gte
    value: "4"
```

篩選在查詢管線層級套用，不是在顯示層過濾——效能有保障。

## 分組：巢狀結構，一目了然

`group_by` 支援多層分組，幫你把物件按屬性歸類：

```yaml
group_by:
  - property: status
  - property: author
```

先按狀態分、再按作者分，結構清楚。

## View 編輯器：直接在 TUI 裡改

不想手動寫 YAML？在視圖模式按 `e`，右邊會打開 view 編輯器。裡面有屬性選取器和運算子選取器，改完自動存檔。

## tmd stats：一個指令看全貌

想知道你的 vault 裡有多少物件？哪個型別最多？哪些屬性最常用？

```bash
tmd stats
```

它會印出彙總統計，幫你快速掌握 vault 的全貌。

## ⚠️ 移除 tmd query

`tmd query` 指令在這個版本移除了。篩選功能現在整合進 View 系統的 filter rule，全文搜尋請用 `tmd search`。如果你有腳本用到 `tmd query`，需要遷移。

## 還有什麼

- **增量索引同步** — TUI 的資料重新整理從完整 reindex 改為增量同步（`Projector.SyncFiles()`），速度更快 (#231)
- **View 模式 session 持久化** — 關掉 TUI 再打開，你還在上次的視圖裡 (#265)
- **TUI overlay 重構** — 說明覆蓋層改用 lipgloss Layer/Compositor，顯示更穩定 (#276)
- **修正** — CJK 文字對齊、側邊欄名稱顯示、view 編輯器的重複欄位等多項修正

## 下一步

View 系統打好了底，接下來我們會往更進階的方向走——更多佈局選項、跨型別查詢、還有 Web UI。

升級方式：

```bash
brew upgrade typemd/tap/typemd-cli
```

或從 [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.5.0) 下載預編譯的執行檔。
