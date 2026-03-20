---
title: "TypeMD 0.3.0 — 你的物件，你的規則"
description: "v0.3.0 帶來物件範本、名稱範本、唯一性約束、系統屬性，以及能在 TUI 裡直接編輯型別 schema 的視覺化編輯器。"
date: 2026-03-14
tags: [發布, 開發日誌]
---

如果說 0.2.0 是「幫屬性定義型別」，那 0.3.0 就是「幫你的 vault 定義規則」。

上個版本你可以描述一本書有哪些欄位、什麼型別。這個版本，你可以告訴 TypeMD：「書名不能重複」、「日記自動用日期命名」、「新的書評一建立就有固定格式」。規則由你寫，TypeMD 幫你守。

## 物件範本：每次建立都從對的地方開始

你有沒有每次建書評，都要手動填 `status: draft`、貼上同一段 Markdown 結構？現在不用了。

在 `templates/<type>/` 放一個 Markdown 檔案，裡面寫好 frontmatter 預設值和正文骨架：

```yaml
# templates/book/review.md
---
status: to-read
---

## Why I Want to Read This

## Notes
```

建立物件時自動套用：

```bash
tmd object create book "Clean Code" --template review
```

如果一個型別只有一個範本，連 `--template` 都不需要，直接套用。多個範本的話，會提示你選擇。

範本可以覆蓋使用者可編輯的系統屬性（`name`、`description`），但 `created_at` 和 `updated_at` 由系統管理，範本中的值會被忽略。

## 名稱範本：讓 TypeMD 幫你命名

有些物件的名稱是有規律的。每日日記就是「日記 2026-03-14」，會議紀錄就是「會議 2026-03-14」。手打很累，也容易出錯。

現在你可以在型別 schema 裡定義名稱範本：

```yaml
name: journal
properties:
  - name: name
    template: "日記 {{ date:YYYY-MM-DD }}"
```

建立物件時不需要給名稱，TypeMD 會自動填入。當然，你也可以手動指定來覆寫。

## 唯一性約束：防止重複

加上 `unique: true`，同一型別中就不能有兩個同名物件：

```yaml
name: person
unique: true
```

第一次建立 `person/john-doe` 沒問題，第二次就會被擋下。不同型別之間不受限制——`person/john-doe` 和 `character/john-doe` 可以並存。

`tmd type validate` 也新增了標籤名稱唯一性檢查。標籤天生就是唯一的，驗證器會幫你抓到重複的。

## 複數顯示名稱

型別 schema 現在可以定義 `plural` 欄位：

```yaml
name: book
plural: books
```

在需要顯示集合名稱的地方（例如 TUI 群組標題），TypeMD 會使用 `books` 而不是 `book`。小改動，但每次看到都會順眼一點。

## 系統屬性：TypeMD 幫你管的事

從 0.3.0 開始，每個物件都自動擁有五個系統屬性：

| 屬性 | 說明 | 可編輯？ |
|------|------|----------|
| `name` | 顯示名稱 | 是 |
| `description` | 描述文字 | 是 |
| `created_at` | 建立時間 | 否（自動設定） |
| `updated_at` | 最後更新時間 | 否（自動更新） |
| `tags` | 標籤（連結到內建 `tag` 型別） | 是 |

這些屬性不需要也不能在型別 schema 中定義——它們是保留名稱。如果你的 schema 有手動定義 `description`、`created_at`、`updated_at` 或 `tags`，升級後驗證會失敗，請先移除。

標籤是全新的內建型別。在物件的 frontmatter 中引用一個不存在的標籤，sync 時會自動建立對應的標籤物件。

## 再見，內建型別

0.2.0 預設附帶 `book`、`person`、`note` 三個型別。我們發現這些預設型別反而讓新使用者困惑——「這些是必要的嗎？可以刪嗎？」

答案是：**你的 vault，你決定裡面有什麼型別**。0.3.0 移除了所有預設型別，`tmd init` 現在只會建立一個乾淨的空 vault。唯一保留的內建型別是 `tag`，因為它是標籤系統的基礎。

如果你已經在用 `book`、`person`、`note`，不用擔心——你的物件和型別 schema 不受影響。只是新建的 vault 不會再自帶這些型別。

## TUI 型別編輯器

以前改型別 schema 只能手動編輯 YAML 檔案。現在打開 TUI，在 sidebar 移到任何型別標題上，就能進入視覺化的型別編輯器：

- **瀏覽**屬性清單，看到型別、emoji、pin 等中繼資料
- **新增**屬性，包含名稱、型別、emoji、pin 等完整設定
- **編輯**既有屬性的所有欄位
- **刪除**屬性（有確認步驟）
- **排序**屬性，拖曳調整順序

不用離開 TUI，不用記 YAML 語法。

## 升級方式

```bash
brew upgrade typemd/tap/typemd-cli
```

或從 [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.3.0) 下載預編譯的執行檔。

升級前請檢查：
1. 移除型別 schema 中手動定義的 `description`、`created_at`、`updated_at`、`tags` 屬性
2. 執行 `tmd type validate` 確認無誤

## 接下來

0.4.0 的重點是「搬進來」——vault 健康檢查（`tmd doctor`）、持續驗證（`tmd type validate --watch`），以及更好的 TUI 建立體驗。讓你的 vault 適合日常使用。

有任何問題或想法，歡迎[開 issue](https://github.com/typemd/typemd/issues)，我們都會看。
