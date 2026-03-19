---
title: 術語表
description: TypeMD 中使用的關鍵術語定義。
sidebar:
  order: 10
---

TypeMD 文件、CLI 和設定檔中使用的術語快速參考。

## Vault

由 TypeMD 管理的目錄。包含 `.typemd/`（設定）、`objects/`（內容），以及選填的 `templates/`（object template）。每個 Vault 都是獨立的，可以用版本控制管理。

## Object（物件）

TypeMD 中知識的基本單位。你不需要用「檔案」或「筆記」來思考，而是用物件——一本書、一個人、一個想法、一場會議。每個物件是一個 Markdown 檔案，包含 YAML frontmatter（屬性）和自由格式的內文。

[深入了解 →](/zh-tw/concepts/objects)

## Type（類型）

定義你要建立什麼樣的物件。每個物件只屬於一個類型（例如 `book`、`person`、`note`）。類型由 `.typemd/types/` 中的 YAML schema 檔案定義。

[深入了解 →](/zh-tw/concepts/types)

## Property（屬性）

物件上的具名欄位，由 Type schema 定義。屬性有型別（`string`、`number`、`date`、`select`、`relation` 等），儲存在 YAML frontmatter 中。屬性可以有 `emoji`、`pin`、`default` 等附加設定。

[深入了解 →](/zh-tw/basics/properties)

## Relation（關聯）

兩個物件之間具名且帶型別的連結。不同於簡單的超連結，關聯定義在 Type schema 中，可以設定為雙向——更新一端會自動更新另一端。

[深入了解 →](/zh-tw/concepts/relations)

## Link（連結）

在 Markdown 內文中使用 `[[type/slug-ulid]]` 語法引用其他物件的方式。TypeMD 會追蹤這些連結及其反向連結（backlink）。

[深入了解 →](/zh-tw/concepts/links)

## Backlink（反向連結）

從連結自動計算的反向參照。如果物件 A 連結到物件 B，則 B 有一個來自 A 的反向連結。反向連結是唯讀的，由索引自動維護。

[深入了解 →](/zh-tw/concepts/links#backlinks)

## Tag（標籤）

內建 `tag` 類型的物件。標籤強制 name 唯一性，可以透過 `tags` 系統屬性連結到任何物件。標籤支援 `color` 和 `icon` 屬性。

[深入了解 →](/zh-tw/basics/tags)

## Template（模板）

建立新物件時提供預設 frontmatter 和內文。Object template 以 Markdown 檔案儲存在 `templates/<type>/` 中。Name template 使用 `{{ date:FORMAT }}` 語法自動產生物件名稱。

[深入了解 →](/zh-tw/basics/templates)

## View（檢視）

已儲存的設定，控制某個 Type 的物件如何顯示 — 包含排序順序、篩選規則、分組和佈局。View 以 YAML 檔案儲存在 `.typemd/types/<name>/views/`。每個 Type 都有一個隱含的預設 View（list 佈局，按 name 排序）。

[深入了解 →](/zh-tw/advanced/file-structure#views)

## System Property（系統屬性）

由 TypeMD 管理、存在於每個 Object 上的屬性：`name`、`description`、`created_at`、`updated_at` 和 `tags`。系統屬性在 frontmatter 中依固定順序排在最前面。使用者撰寫的系統屬性（`name`、`description`、`tags`）可以被 template 覆蓋；自動管理的（`created_at`、`updated_at`）不能覆蓋。

[深入了解 →](/zh-tw/basics/properties#系統屬性)

## Shared Property（共用屬性）

定義在 `.typemd/properties.yaml` 中的可重複使用屬性定義，type schema 可透過 `use` 關鍵字引用。共用屬性確保跨 type 的定義一致。引用時可以依 type 覆寫 `pin`、`emoji` 和 `description`。

[深入了解 →](/zh-tw/basics/properties#共用屬性)

## Query（查詢）

對索引進行的結構化篩選與排序請求。查詢使用篩選規則（屬性 + 運算子 + 值）和排序規則（屬性 + 方向）來擷取指定 Type 的符合物件。透過 CLI（`tmd query`）公開，View 也在內部使用。

[深入了解 →](/zh-tw/basics/queries)

## Search（搜尋）

跨物件名稱、描述和內文的全文搜尋。不同於查詢（Query）是針對結構化屬性篩選，搜尋會將自由文字輸入與搜尋索引比對。透過 CLI（`tmd search`）和 TUI 搜尋列公開。

## Slug

物件檔名中人類可讀的部分。例如 `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md` 中的 `golang-in-action`。來自建立物件時提供的名稱。

## ULID

26 字元的通用唯一且可按字典序排列的識別碼（Universally Unique Lexicographically Sortable Identifier），由 CLI 自動附加在 slug 後方。確保即使 slug 相同也不會衝突。手動建立的物件不需要 ULID。

## Object ID

物件的完整識別碼，格式為 `type/slug-ulid`。例如：`book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`。用於連結目標和關聯參照。

## Frontmatter

物件檔案頂部的 YAML 中繼資料區塊，以 `---` 分隔。依固定順序包含屬性：系統屬性在前，schema 定義的屬性在後。

[深入了解 →](/zh-tw/advanced/frontmatter)

## Index（索引）

SQLite 資料庫（`.typemd/index.db`），快取物件中繼資料以供快速查詢和全文搜尋。自動建立和更新，不需要手動編輯。

## Sync（同步）

讀取物件檔案並更新索引的過程。開啟 Vault 時自動發生。同步時會補上缺少的 `name` 屬性、解析標籤參照、並重新整理搜尋索引。

## Validation（驗證）

檢查 schema、物件、關聯、連結和 name 唯一性約束是否正確。使用 `tmd type validate` 執行。

[深入了解 →](/zh-tw/basics/validation)
