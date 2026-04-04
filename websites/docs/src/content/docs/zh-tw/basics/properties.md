---
title: 屬性
description: 屬性如何定義 Object 的結構。
sidebar:
  order: 1
---

屬性（Property）是 Object 上的具名欄位，由 [Type schema](/zh-tw/concepts/types) 定義。它們儲存在 Markdown 檔案的 YAML frontmatter 中，讓每個 Object 在自由格式的內文之外，還擁有結構化、可查詢的資料。

## 系統屬性

每個 Object 都支援七個由 TypeMD 管理的系統屬性。這些屬性提供每個知識管理工具都需要的基礎中繼資料——身份識別、描述、時間追蹤、分類、保護和歸檔——使用者不需要在每個 type schema 中自行定義。並非所有系統屬性都會出現在每個 Object 的 frontmatter 中：`name` 會在建立時自動填入，`created_at` 和 `updated_at` 在使用 CLI 時設定，而 `description`、`tags`、`locked` 和 `archived` 只在使用者明確設定後才會出現。

| 屬性 | 說明 | 可變性 | 為什麼需要 |
|------|------|--------|-----------|
| `name` | 顯示名稱，建立時從 slug 自動填入 | 使用者撰寫 | 將顯示名稱與檔名解耦——支援空格、大小寫和重新命名而不需搬動檔案 |
| `description` | 選填的單行摘要，用於列表顯示和搜尋結果 | 使用者撰寫 | 提供跨所有 type 的一致摘要欄位，用於列表顯示、搜尋結果和 API 回應 |
| `created_at` | 建立時間戳記，RFC 3339 格式（僅設定一次，不會修改） | 自動管理 | 支援依建立日期排序，理解 vault 的時間軸 |
| `updated_at` | 最後修改時間戳記，RFC 3339 格式（每次儲存時更新） | 自動管理 | 支援依最近修改排序，追蹤 Object 的演變 |
| `tags` | 標籤參照的陣列（關聯到內建 `tag` 型別，支援多值） | 使用者撰寫 | 跨 type 的橫切分類 |
| `locked` | 鎖定物件，防止編輯 | 使用者撰寫 | 保護已完成或歸檔的 Object 不被意外修改 |
| `archived` | 軟刪除標記——從預設查詢中隱藏物件 | 使用者撰寫 | 讓不再活躍的 Object 退場而不永久刪除 |

**使用者撰寫**的屬性（`name`、`description`、`tags`、`locked`、`archived`）可以被 [object template](/zh-tw/basics/templates) 覆蓋。**自動管理**的屬性（`created_at`、`updated_at`）無法被覆蓋——它們永遠反映實際的建立和修改時間。

這些名稱是保留的，不能在 type schema 或[共用屬性](#共用屬性)中使用。唯一的例外是 `name`，可以出現在 `properties` 中並搭配 `template` 欄位來設定 [name template](/zh-tw/basics/templates#name-template)。

## 屬性型別

Type schema 中的每個屬性都有一個 `type`，決定它接受什麼值以及如何驗證。

| 型別 | 說明 | 值範例 |
|------|------|--------|
| `string` | 純文字 | `"Go in Action"` |
| `number` | 整數或浮點數 | `42`、`3.14` |
| `date` | `YYYY-MM-DD` 格式的日期 | `"2026-03-09"` |
| `datetime` | ISO 8601 日期時間 | `"2026-03-09T10:30:00+08:00"` |
| `url` | 網址（僅限 http/https） | `"https://example.com"` |
| `checkbox` | 布林值 | `true`、`false` |
| `select` | 從選項中單選 | `"reading"` |
| `multi_select` | 從選項中多選 | `["fiction", "sci-fi"]` |
| `relation` | 連結到另一個 Object | `"person/alan-donovan"` |

### 基本型別

- **string** — 接受任何文字值。
- **number** — 接受整數和浮點數。
- **date** — 接受 `YYYY-MM-DD` 格式的日期。YAML 自動解析的日期（例如不加引號的 `2026-03-09`）也會被接受。
- **datetime** — 接受 ISO 8601 日期時間值，必須至少包含時和分。
- **url** — 接受以 `http://` 或 `https://` 開頭的網址。不支援其他協定（ftp、ssh 等）。
- **checkbox** — 僅接受布林值 `true` 或 `false`。字串值如 `"true"` 會被拒絕。

### 選擇型別

**select** — 從預定義清單中單選。需要 `options` 陣列，每個項目有 `value`（必填）和可選的 `label` 顯示名稱：

```yaml
- name: status
  type: select
  options:
    - value: to-read
      label: To Read
    - value: reading
    - value: done
```

省略 `label` 時，預設等於 `value`。

**multi_select** — 從預定義清單中多選。`options` 格式與 `select` 相同。屬性值為字串列表。單一字串值（例如 `genres: fiction`）會自動視為單元素列表（`["fiction"]`）。

### Relation 型別

將 Object 連結到另一個 Object。`target`、`bidirectional`、`inverse` 和 `multiple` 欄位的詳細說明，請參閱 [Relation](/zh-tw/concepts/relations) 頁面。

```yaml
- name: author
  type: relation
  target: person
  bidirectional: true
  inverse: books
```

## 屬性附加設定

除了 `name` 和 `type` 之外，屬性還支援可選的附加設定：

| 設定 | 型別 | 說明 |
|------|------|------|
| `emoji` | string | 簡潔顯示用的視覺圖示（在同一 type 內必須唯一） |
| `description` | string | 屬性用途的自由文字說明 |
| `pin` | integer | 正整數，讓屬性突出顯示在 TUI body 面板頂端。數值越小排越前面。已置頂的屬性不會出現在 Properties 面板中。 |
| `default` | any | 建立新 object 時指定的預設值 |

**為什麼需要 emoji？** 在緊湊的 UI 環境中（TUI 屬性面板、表格欄位），屬性名稱會佔用大量水平空間。Emoji 提供空間效率高的視覺識別，一眼就能辨認。

**為什麼需要 pin？** 沒有置頂功能時，所有屬性在 Properties 面板中平等顯示。使用者必須掃描整個列表才能找到 status 或 rating 等關鍵中繼資料。置頂讓 type 作者標記特定屬性，突出顯示在 body 面板頂端。

### 置頂範例

```yaml
- name: status
  type: select
  emoji: 📋
  pin: 1
  options:
    - value: to-read
    - value: reading
    - value: done
```

`pin` 值必須是正整數，且在同一 type schema 內唯一。沒有 `pin`（或 `pin: 0`）的屬性會照常顯示在 Properties 面板中。

## 共用屬性

當同一個屬性出現在多個 type 中，在每個 schema 中獨立定義會導致重複和不一致——`project` 的 `due_date` 可能用 `date` 型別，而 `task` 的卻用了 `datetime`。共用屬性讓你定義一次、到處引用，確保跨 type 的定義一致。

如果同一個屬性出現在多個 type 中（例如 `due_date` 同時在 `project` 和 `task` 裡），你可以在 `properties/properties.yaml` 中定義一次，然後用 `use` 關鍵字引用。

### 定義共用屬性

```yaml
# properties/properties.yaml
properties:
  - name: due_date
    type: date
    emoji: "\U0001F4C5"
  - name: priority
    type: select
    options:
      - value: high
      - value: medium
      - value: low
```

每個項目支援與 type schema 屬性相同的所有欄位：`name`、`type`、`emoji`、`pin`、`options`、`target`、`default`、`multiple`、`bidirectional` 和 `inverse`。在 type schema 中透過 `use:` 引用時，只有 `pin`、`emoji` 和 `description` 可以被覆寫——其他欄位皆繼承自共用屬性定義。

此檔案為選填。如果 `properties/properties.yaml` 不存在或沒有 `properties` 陣列，TypeMD 會將其視為空的共用屬性集合。

### 在 type schema 中引用

使用 `use` 關鍵字依名稱引用共用屬性：

```yaml
# types/project/schema.yaml
name: project
emoji: "\U0001F4CB"
properties:
  - name: title
    type: string
  - use: due_date
  - use: priority
  - name: budget
    type: number
```

載入 type 時，每個 `use` 項目會被解析為完整的屬性定義。你可以混合使用 `use` 項目和一般的屬性定義——順序會被保留。

### 覆寫欄位

透過 `use` 引用共用屬性時，你可以覆寫這些欄位：

| 欄位 | 說明 |
|------|------|
| `pin` | 覆寫此 type 中的置頂位置 |
| `emoji` | 覆寫此 type 中的顯示 emoji |
| `description` | 覆寫此 type 的屬性說明 |

其他所有欄位不能覆寫——它們來自共用定義。

```yaml
# types/task/schema.yaml
name: task
properties:
  - use: due_date
    pin: 1
    emoji: "\U0001F5D3\uFE0F"
  - use: priority
```

### 共用屬性 vs. 內嵌屬性

| | 共用屬性 | 內嵌屬性 |
|---|---------|---------|
| 定義位置 | `properties/properties.yaml` | `types/<type>/schema.yaml` |
| 可重複使用 | 是——透過 `use` 引用 | 否——僅限單一 type |
| 可依 type 自訂 | `pin`、`emoji` 和 `description` | 完全可自訂 |
| 適用情境 | 多個 type 共用的屬性 | 單一 type 獨有的屬性 |

**經驗法則**：如果某個屬性以相同的定義出現在兩個以上的 type 中，就把它做成共用屬性。

## 相關頁面

- [Type](/zh-tw/concepts/types) — Type 如何定義 Object 結構
- [Relation](/zh-tw/concepts/relations) — 連結 Object
- [驗證](/zh-tw/basics/validation) — 屬性值如何被驗證
