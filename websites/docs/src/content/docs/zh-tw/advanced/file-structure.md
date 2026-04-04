---
title: 檔案結構
description: TypeMD 如何在 Vault 中組織檔案與目錄。
sidebar:
  order: 1
---

TypeMD 將所有內容以純檔案儲存在磁碟上。本頁說明直接操作 Vault 時會遇到的目錄配置和檔案格式。

## Vault 目錄配置

Vault 是一個具有特定結構的一般目錄：

```
vault/
├── .typemd/
│   ├── types/              # 使用者定義的 type schema（目錄格式）
│   │   ├── book/
│   │   │   ├── schema.yaml # type schema 定義
│   │   │   └── views/      # 此 type 的已儲存 view（選填）
│   │   │       ├── default.yaml
│   │   │       └── by-rating.yaml
│   │   └── person/
│   │       └── schema.yaml
│   ├── config.yaml         # vault 設定（選填）
│   ├── properties.yaml     # 共用屬性定義（選填）
│   ├── index.db            # SQLite 索引（自動更新）
│   └── tui-state.yaml      # TUI 會話狀態（自動儲存）
├── templates/              # 依 type 分類的 object template（選填）
│   └── book/
│       └── review.md       # 新 object 的預設 frontmatter 和內文
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
```

- `.typemd/` — 設定與內部狀態（vault 設定、type schema、共用屬性、索引、TUI 狀態）
- `templates/` — 選填的 object template，依 type 分類
- `objects/` — 所有 Object 檔案，依 type 分類

## Object 檔案

Object 以 Markdown 檔案搭配 YAML frontmatter 儲存在 `objects/<type>/` 底下。每個檔案代表一個 Object。

### Object ID

每個 Object 都有唯一的 ID，格式為 `type/<slug>-<ulid>`，例如：

```
book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz
```

- **type** — Object 的型別（對應 `objects/` 底下的子目錄名稱）
- **slug** — 從名稱衍生的人類可讀識別碼
- **ulid** — 附加的 26 字元小寫 [ULID](https://github.com/ulid/spec)，確保唯一性

### ULID

透過 CLI 建立 Object 時（`tmd object create`），ULID 會自動附加在 slug 後方。如果你手動建立檔案（不使用 CLI），ULID 是選填的——TypeMD 有沒有 ULID 都能正常運作。這讓你可以輕鬆將既有的 Markdown 檔案加入 Vault。

### 系統屬性

所有 Object 都有七個由 TypeMD 管理的系統屬性。這些屬性在 frontmatter 中永遠依以下固定順序排在最前面：

| 屬性 | 說明 | 可變性 |
|------|------|--------|
| `name` | 顯示名稱，建立時從 slug 自動填入 | 使用者撰寫 |
| `description` | 選填的單行摘要，用於列表顯示和搜尋結果 | 使用者撰寫 |
| `created_at` | 建立時間戳記，RFC 3339 格式（僅設定一次，不會修改） | 自動管理 |
| `updated_at` | 最後修改時間戳記，RFC 3339 格式（每次儲存時更新） | 自動管理 |
| `tags` | 標籤參照的陣列（關聯到內建 `tag` 型別，支援多值） | 使用者撰寫 |
| `locked` | 鎖定物件，防止編輯 | 使用者撰寫 |
| `archived` | 軟刪除標記——從預設查詢中隱藏物件 | 使用者撰寫 |

系統屬性之後才是 schema 定義的屬性。這些名稱是保留的，不能在 type schema 或 shared properties 中使用。

**使用者撰寫**的屬性（`name`、`description`、`tags`、`locked`、`archived`）可以被 object template 覆蓋。**自動管理**的屬性（`created_at`、`updated_at`）無法被覆蓋——它們永遠反映實際的建立和修改時間。

關於 frontmatter 格式和手動編輯的詳細說明，請參閱 [Frontmatter](/zh-tw/advanced/frontmatter)。

## Type schema 檔案

每個 type 以目錄形式儲存在 `.typemd/types/<name>/` 下，包含 `schema.yaml` 檔案：

```yaml
# .typemd/types/book/schema.yaml
name: book
plural: books
unique: false
properties:
  - name: status
    type: select
    options:
      - value: reading
      - value: completed
      - value: want-to-read
  - name: author
    type: relation
    target: person
```

有兩個內建 type：`tag`（支援 `tags` 系統屬性，複數形式「tags」，`unique: true`）和 `page`（通用內容容器，複數形式「pages」，emoji 📄）。內建 type 無法刪除，但可以透過自訂 type schema 覆蓋。

完整的 type schema 格式請參閱 [Type](/zh-tw/concepts/types)。

## Views

每個 type 可以有多個已儲存的 view，定義物件的排序、篩選和顯示方式。View 以 YAML 檔案儲存在 `.typemd/types/<name>/views/`：

```yaml
# .typemd/types/book/views/by-rating.yaml
name: by-rating
layout: list
columns: [status, rating]
sort:
  - property: rating
    direction: desc
filter:
  - property: status
    operator: is
    value: reading
group_by:
  - property: genre
```

有兩種佈局可用：

- **`list`** — 顯示 object 名稱，搭配選填的行內屬性值。預設欄位：無（僅名稱）。
- **`table`** — 顯示欄位式表格，包含 NAME 欄位以及帶標題的屬性欄位。預設欄位：所有屬性。

選填的 `columns` 欄位指定要顯示哪些屬性。設定後，兩種佈局都會顯示指定的屬性。

每個 type 都有一個隱含的預設 view（list 佈局，按 name 升序排列）。自訂預設 view 後會儲存為 `views/default.yaml`。若沒有 view 檔案，type 目錄不需要 `views/` 子目錄。

## 共用屬性檔案

選填的 `.typemd/properties.yaml` 定義可重複使用的屬性定義，type schema 可透過 `use` 引用：

```yaml
# .typemd/properties.yaml
properties:
  - name: status
    type: select
    options:
      - value: active
      - value: archived
```

關於共用屬性的詳細說明，請參閱[共用屬性](/zh-tw/basics/properties#共用屬性)。

## Object Template

每個 type 可以在 `templates/<type>/<name>.md` 定義一個或多個 template。Template 是一般的 Markdown 檔案（frontmatter + 內文），用於提供建立 Object 時的預設內容。

```bash
# 單一 template — 自動套用
tmd object create book clean-code

# 明確指定 template
tmd object create book clean-code -t review

# 多個 template，未指定 flag — 互動式選擇
tmd object create book clean-code
```

Template 的 frontmatter 屬性會覆蓋 schema 預設值。Template 的內文會成為 Object 的初始內文。Template 中的自動管理系統屬性（`created_at`、`updated_at`）會被忽略。
