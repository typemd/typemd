---
title: 標籤
description: 標籤如何為你的 Object 提供跨類型分類。
sidebar:
  order: 2
---

標籤（Tag）讓你可以跨 Type 分類 Object。一本 `book` 和一篇 `note` 可以共用同一個標籤，非常適合用於跨類型的主題，例如「golang」、「important」或「to-review」。

## 標籤是一級 Object

不同於簡單的字串標記，TypeMD 中的標籤是由內建 `tag` type 支援的完整 Object。每個標籤在 `objects/tag/` 底下都有自己的 Markdown 檔案，包含兩個可選屬性：

| 屬性 | 型別 | 說明 |
|------|------|------|
| `color` | string | 標籤的顯示顏色 |
| `icon` | string | 標籤的顯示圖示 |

因為標籤是 Object，它們可以有 description、內文和 wiki-link——就跟 vault 中的其他 Object 一樣。

## 標籤唯一性

內建的 `tag` type 在 schema 中設有 `unique: true`，因此標籤名稱必須唯一。你不能有兩個 `name` 值相同的 tag Object。此唯一性約束在建立時強制執行，並由 `tmd type validate` 驗證。

任何使用者自訂的 type 也可以透過在 schema 中加入 `unique: true` 來啟用 name 唯一性。詳情請參閱[資料模型](/zh-tw/developers/data-model#唯一性約束機制)。

## 使用標籤

每個 Object 都支援 `tags` [系統屬性](/zh-tw/basics/properties#系統屬性)，存放對 tag Object 的參照。`tags` 屬性只在使用者明確設定後才會出現在 frontmatter 中。在 frontmatter 中，標籤參照可以使用完整的 Object ID 或標籤名稱：

```yaml
tags:
  - tag/go
  - tag/programming-01jqr3k5mpbvn8e0f2g7h9txyz
```

兩種形式都有效。名稱參照（如 `tag/go`）會在同步時被解析為完整 ID。

## 標籤參照解析

同步時，TypeMD 會將名稱參照解析為完整的 Object ID。如果參照的標籤不存在，TypeMD 會自動建立一個空內文的標籤。這表示你可以在 frontmatter 中自由加入標籤，不需要先建立 tag Object——TypeMD 會幫你處理。

## 標籤 vs. Type

標籤和 Type 有不同的用途：

- **Type** 定義結構——一本 `book` 永遠有 title、status 和 rating。
- **標籤** 跨 Type 分類——一本 `book` 和一篇 `note` 可以共用同一個標籤。

使用 Type 確保結構一致性；使用標籤進行跨類型分類。

## 相關頁面

- [屬性](/zh-tw/basics/properties#系統屬性) — `tags` 系統屬性
- [Type](/zh-tw/concepts/types#內建-type) — 內建的 `tag` type
- [驗證](/zh-tw/basics/validation) — 標籤的 name 唯一性驗證
