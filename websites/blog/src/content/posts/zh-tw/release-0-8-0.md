---
title: "TypeMD 0.8.0 — 你的 Vault，一目了然"
description: "v0.8.0 把 type schema 和 properties 搬到看得見的地方，加上物件歸檔、圖形匯出、模板 CLI——讓你的知識庫更透明、更好管。"
date: 2026-04-05
tags: [發布, 開發日誌]
---

0.7.0 讓你在 TUI 裡直接改所有東西。但有些東西你可能根本不知道放在哪裡——type schema 和 shared properties 藏在 `.typemd/` 目錄裡，不翻設定檔根本看不到。這版我們把它們搬出來，讓 vault 的結構一目了然。

## Type Schema 搬到 Vault 根目錄

這是一個 **破壞性變更**。`types/` 和 `properties/` 目錄從 `.typemd/` 搬到 vault 根目錄：

```
# 之前
my-vault/
├── .typemd/
│   ├── types/book/schema.yaml
│   └── properties/favorite.yaml
└── objects/

# 現在
my-vault/
├── types/book/schema.yaml
├── properties/favorite.yaml
├── objects/
└── .typemd/          ← 只剩內部狀態
```

原則很簡單：**你會編輯的檔案放根目錄，TypeMD 管理的內部狀態放 `.typemd/`**。開啟 vault 時自動遷移，不需要手動搬檔案。(#362)

## Shared Properties 拆成獨立檔案

之前所有 shared properties 擠在一個 `properties.yaml`。現在每個 property 是獨立檔案：

```yaml
# properties/due_date.yaml
type: date
emoji: "📅"
```

```yaml
# properties/priority.yaml
type: select
options:
  - value: high
  - value: medium
  - value: low
```

一個屬性一個檔案，更容易找、更好管、更適合版本控制。舊格式自動遷移。(#363)

## 物件歸檔：軟刪除不心疼

不再需要的物件，又捨不得永久刪除？現在可以歸檔：

```bash
tmd object archive book/old-notes
```

歸檔的物件從預設查詢和 TUI 中隱藏，但檔案還在。想看全部？加 `--archived` 就好。想恢復？`tmd object unarchive`。(#34)

## Relation Graph 匯出

想看 vault 裡物件之間的關係圖？

```bash
tmd graph | dot -Tpng -o graph.png
```

`tmd graph` 把所有 relation 和 wiki-link 匯出成 DOT 格式，丟給 Graphviz 就能畫圖。支援 `--type` 篩選特定類型。(#25)

## Template CLI

模板管理不用再手動開檔案了：

```bash
tmd template list                    # 列出所有模板
tmd template show book/quick-note    # 顯示模板內容
tmd template create book new-review  # 建立新模板
tmd template delete book/quick-note  # 刪除模板
```

(#248)

## `tmd log`：物件的 Git 歷史

想知道一個物件被改了什麼？

```bash
tmd log book/golang-in-action
```

直接看這個物件檔案的 git commit 歷史，不用自己拼路徑。(#24)

## Markdown 渲染

TUI body 面板現在會渲染 Markdown——標題、粗體、斜體、程式碼區塊、引用、清單都有對應的樣式和顏色。(#103)

## 還有什麼

- **`tmd type validate --watch`** — 持續監控 schema 和物件，檔案一改就重新驗證 (#200)
- **SQLite 故障 fallback** — 索引檔損壞或鎖定時，自動退回檔案系統查詢，不會中斷工作 (#180)
- **TUI 設定頁** — 按 `,` 開啟設定頁面，不用手動編輯 config.yaml (#355)

## 遷移注意事項

`types/` 和 `properties/` 目錄搬到 vault root 是自動遷移——開啟 vault 時 TypeMD 會自動搬移。但如果你的 `.gitignore` 有特別排除這些路徑，記得更新。

## 下一步

結構打好了。0.9.0 我們要擴充型別系統——computed properties、aliases、和更豐富的 relation 語意。

升級方式：

```bash
brew install typemd/tap/typemd-cli
```

或從 [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.8.0) 下載預編譯的執行檔。
