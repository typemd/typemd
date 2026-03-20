---
title: "TypeMD 0.2.0 — 屬性型別系統、共用屬性，以及更好看的 TUI"
description: "v0.2.0 帶來 9 種正式定義的屬性型別、共用屬性機制，以及支援 emoji、置頂欄位與 session 持久化的全新 TUI 體驗。"
date: 2026-03-11
tags: [發布, 開發日誌]
---

TypeMD 從第一個版本發布後，我們下一個想完成的目標是：「屬性到底支援哪些型別？」v0.1.0 中，屬性的 `type` 欄位雖然存在，但從來沒有正式定義過。你可以寫 `type: number`，大致上也能運作，但背後沒有真正的約束。

0.2.0 改變了這件事。這個版本的核心，是讓 TypeMD 的型別系統變成你真的可以依賴的東西。

## 九種屬性型別，正式定義

型別 schema 現在支援 9 種明確的屬性型別：

| 型別 | 說明 |
|------|------|
| `string` | 短文字 |
| `number` | 數值 |
| `date` | 日期 |
| `datetime` | 日期加時間 |
| `url` | 網址 |
| `checkbox` | 勾選框（是／否） |
| `select` | 單選 |
| `multi_select` | 多選 |
| `relation` | 連結到另一個型別的物件 |

舉例來說，一個書籍型別可以長這樣：

```yaml
name: book
emoji: 📚
properties:
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
  - name: rating
    type: number
  - name: published
    type: date
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
```

有了明確的型別，驗證、遷移、顯示都能更精準。

## 不用再重複定義：共用屬性

如果你有多個型別需要相同的屬性——例如每種資源都有 `tags` 和 `favorite`——以往只能在每個型別 schema 裡各寫一遍。現在不用了。

在 `.typemd/properties.yaml` 集中定義一次：

```yaml
properties:
  - name: tags
    type: multi_select
    emoji: 🔖
    options:
      - value: programming
      - value: design
  - name: favorite
    type: checkbox
    emoji: ❤️
```

再透過 `use` 在任何型別中參照：

```yaml
properties:
  - name: rating
    type: number
  - use: tags
  - use: favorite
```

改一個地方，所有型別自動更新。

## 長得像你的 vault 的 TUI

### Emoji（如果你想要的話）

型別和屬性現在都可以加 `emoji` 欄位：

```yaml
name: book
emoji: 📚
properties:
  - name: rating
    type: number
    emoji: ⭐
```

型別 emoji 會出現在物件列表的群組標頭，屬性 emoji 會出現在詳細面板。你的 vault 開始有自己的樣子。

### 專屬標題列

瀏覽物件時，TUI 頂部現在有一條標題列，顯示型別 emoji 與物件名稱。改動很小，但定向感差很多。

### 置頂重要屬性

在屬性上標記 `pin` 數值，它就會按順序固定顯示在詳細面板最上方。不用再滾過一堆欄位才找到你最在意的那個：

```yaml
  - name: status
    type: select
    pin: 1
  - name: rating
    type: number
    pin: 2
```

### 從你離開的地方繼續

關掉 TUI 再打開，游標位置、選取物件、面板狀態全部恢復。TUI 從一個瀏覽器，開始變得更像一個工作空間。

## 其他改善

- **`name` 現在是系統屬性** — 自動從 slug 填入，不需要也不能在型別 schema 裡手動定義。
- **`--readonly` 旗標** — 以唯讀模式啟動 TUI，在不想誤觸編輯的情境下很好用。
- **`--reindex` 旗標** — 取代原本的 `tmd reindex` 子指令，在啟動時傳入即可強制重建索引。
- **前綴比對** — 輸入物件 ID 時，不需要完整的 ULID，短前綴就夠了。
- **Homebrew** — `brew install typemd/tap/typemd-cli`，macOS 使用者現在有最簡單的安裝方式。

## 升級方式

```bash
brew upgrade typemd/tap/typemd-cli
```

或從 [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.2.0) 下載預編譯的執行檔。

## 接下來

0.3.0 的重點會放在 Web UI——讓你不用開終端機也能瀏覽 vault。一樣的 local-first、以檔案為本，上面多一個圖形介面。

有任何問題或想法，歡迎[開 issue](https://github.com/typemd/typemd/issues)，我們都會看。
