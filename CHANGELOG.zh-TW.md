# 變更日誌

本檔案記錄本專案所有值得注意的變更。

格式依循 [Keep a Changelog](https://keepachangelog.com/)。

## [v0.8.0] - 2026-04-05

### 破壞性變更

- types/ 和 properties/ 搬到 Vault 根目錄 — type schema 和 shared property 檔案從 `.typemd/` 搬到 vault 根目錄的 `types/` 和 `properties/`。開啟 vault 時自動遷移。如有需要請更新 `.gitignore`。(#362)

### 新增

- 物件歸檔 — 使用 `tmd object archive` 軟刪除物件；歸檔的物件從預設查詢中隱藏但檔案保留；用 `tmd object unarchive` 恢復 (#34)
- 圖形匯出 — `tmd graph` 將 relation 和 wiki-link 匯出為 DOT 格式，可用 Graphviz 視覺化；支援 `--type` 篩選 (#25)
- Template CLI — `tmd template list/show/create/delete` 指令，不用手動開檔案即可管理模板 (#248)
- `tmd log` — 顯示特定物件的 git commit 歷史 (#24)
- `tmd type validate --watch` — 持續監控，檔案變更時自動重新驗證 (#200)
- Markdown 渲染 — TUI body 面板渲染標題、粗體、斜體、程式碼區塊、引用和清單，支援語法高亮 (#103)
- 設定頁面 — 在 TUI 中按 `,` 開啟互動式設定編輯器 (#355)
- 獨立屬性檔案 — shared properties 從單一 `properties.yaml` 拆為個別的 `properties/<name>.yaml` 檔案；自動遷移 (#363)
- SQLite Fallback — 索引不可用時自動退回檔案系統查詢，不中斷工作 (#180)

### 變更

- Reconciler/Projector 拆分 — 內部重構，將檔案正規化（Reconciler）與索引寫入（Projector）分離，職責更清楚 (#339)
- 驗證邏輯整合 — 將重複的驗證邏輯整合至單一 `type_schema_validate.go` (#340)

[v0.8.0]: https://github.com/typemd/typemd/releases/tag/v0.8.0

## [v0.7.0] - 2026-03-28

### 破壞性變更

- 移除 `--reindex` 旗標 — 全域 `--reindex` 旗標已移除；索引現在在每次開啟 vault 時自動同步，不再需要手動重建索引 (#325)
- 移除舊版單檔型別 Schema — 不再支援 `.typemd/types/<name>.yaml` 格式；僅接受目錄格式 `.typemd/types/<name>/schema.yaml`。若仍有單檔 schema，請先在 v0.6.0 執行 `tmd migrate` (#273)
- 面板調整快捷鍵 — 從 `[`/`]` 改為 `-`/`=`

### 新增

- 行內屬性編輯 — 在 TUI 的屬性上按 Enter 即可原地編輯；支援所有屬性型別，包含 string、number、date、select、checkbox、url、text (#87)
- 關聯選取器 — 關聯屬性開啟模糊搜尋選取器，瀏覽並選取目標物件；支援單值與多值關聯 (#88)
- 表格儲存格編輯 — 在表格視圖中按 Enter 即可行內編輯儲存格的屬性值，使用與屬性面板相同的編輯元件 (#316)
- 日期選取器 — date 和 datetime 屬性使用分段輸入搭配行內日曆彈出視窗，精確選取日期 (#324)
- 物件鎖定 — 在 frontmatter 設定 `locked: true` 防止意外編輯；TUI 顯示鎖定指示並阻擋對已鎖定物件的編輯 (#157)
- 可設定的日期格式 — 在 `.typemd/config.yaml` 設定 `date_format` 和 `datetime_format`，自訂 TUI 中的日期顯示格式 (#323)
- 本地屬性區段 — 未在型別 schema 中定義的屬性，在 TUI 中以視覺區隔的「Local Properties」區段顯示，同步時予以保留 (#285)
- 本地 LLM 支援 — 透過 `ai.providers` 設定 OpenAI 相容的 provider（Ollama、LM Studio、vLLM）；用 `ai.default` 切換 provider；既有的 `ai.enabled`/`ai.model` 設定自動遷移 (#319)
- `tmd serve` — 新指令，啟動 HTTP 伺服器，提供 REST API 與 Vue 3 Web 前端，透過瀏覽器存取 vault；內建三色主題系統（warm/dark/light）(#125)
- 結構化日誌 — 所有套件使用 `slog` 進行結構化 JSON 記錄；TUI 寫入 `.typemd/logs/`，CLI 使用 `--debug` 輸出至 stderr (#326)
- 專注模式 — 在 TUI 中按 `.` 切換單欄全寬的 body 面板

### 修正

- 範本刪除 — 在型別編輯器中按 `d` 鍵現在能正確刪除範本 (#317)
- 屬性縮排 — 屬性值現在與 TUI 面板標題正確對齊

[v0.7.0]: https://github.com/typemd/typemd/releases/tag/v0.7.0

## [v0.6.0] - 2026-03-22

### 新增

- AI MVP — TUI 中的 AI 功能：自動產生描述、建議標籤、schema 探索；在物件上按 `g` 開啟 AI 動作、在型別上按 `Ctrl+E` 進入 schema 探索 (#294)
- `tmd instructions` 指令 — 輸出嵌入的 skill 指令並附帶 vault context（型別摘要），以 JSON 格式呈現；支援 `--skill` 輸出原始 SKILL.md、`--json` 列出所有 skill (#311)
- Skill 指令覆寫 — 在 `.typemd/instructions/<skill>.md` 放置檔案即可覆寫嵌入的 skill (#311)
- Marketplace Guides Plugin — `vault-guide` skill 教 AI 管理 vault；`instructions-guide` skill 教 AI 使用 `tmd instructions`
- Shell 自動補全 — 在所有 CLI 指令中 tab 補全物件 ID、型別名稱和關聯名稱 (#119)
- 互動式消歧義 — 當前綴匹配多個物件時彈出模糊選取器 (#147)
- 簡寫 Wiki-Links — `[[name]]` 和 `[[type/name]]` 語法在同步時自動解析為完整 ID (#176)
- Stats 儀表板 — TUI 中的統計模式，顯示 vault 總覽和各型別統計 (#280)
- Toast 通知 — TUI 中的暫時性覆蓋訊息，用於同步警告、AI 錯誤等事件 (#303)
- `tmd format` 指令 — 正規化物件 frontmatter 的排序和格式 (#281)
- 同步警告 — TUI 透過 toast 顯示找不到和模糊的關聯參照警告 (#297)
- 物件重新命名 — 在 TUI 中按 `r` 重新命名物件

### 變更

- 關聯中的前綴匹配 — 關聯屬性值在同步時支援前綴匹配 (#74)
- Wiki-Link 驗證 — `ValidateWikiLinks` 改從檔案解析而非資料庫，提升準確性 (#306)

### 修正

- 修正所有文件中的 brew formula 名稱為 `typemd-cli`

[v0.6.0]: https://github.com/typemd/typemd/releases/tag/v0.6.0

## [v0.5.0] - 2026-03-20

### 破壞性變更

- 移除 `tmd query` — 篩選功能改由 view filter（結構化 `FilterRule` 物件）與 `tmd search` 取代；使用 `tmd query` 的腳本需遷移 (#267)

### 新增

- View 系統 — 在 `.typemd/types/<name>/views/<view>.yaml` 定義型別的儲存視圖，含佈局、篩選、排序、分組與欄位設定 (#95, #256)
- 列表視圖佈局 — `layout: list` 以名稱列表呈現物件，可附帶行內值；為每個型別的預設視圖 (#256)
- 表格視圖佈局 — `layout: table` 以欄位表格呈現物件，可透過 `columns` 設定顯示欄位 (#97, #282)
- View 編輯器 — 在視圖模式按 `e` 開啟側邊面板編輯器，支援篩選規則、排序規則、分組規則的屬性／運算子選取器與自動儲存 (#258, #279)
- View 篩選 — 型別感知的篩選運算子（`is`、`contains`、`before`、`gt` 等）在查詢管線層級套用；每種屬性型別有各自的有效運算子集 (#263, #264)
- 多層分組 — `group_by` 支援 `{property: string}` 規則陣列，可進行巢狀分組；舊版單一字串格式自動遷移 (#279)
- View 模式 Session 持久化 — TUI 記住開啟中的視圖（型別、視圖名稱、游標、捲動、展開分組）跨重啟還原 (#265)
- View 選取篩選 — 選取視圖的彈出視窗支援文字篩選 (#259)
- `tmd stats` 指令 — 彙總統計：總物件數、各型別物件數、屬性使用率 (#30)
- 增量索引同步 — TUI 使用 `Projector.SyncFiles()` 進行基於路徑的增量同步，取代完整重建索引 (#231)

### 變更

- 型別 Schema 目錄格式 — `MigrateSchemas` 現在同時處理目錄格式與舊版單檔格式的型別 schema (#260)
- 統一屬性格式化 — `FormatValue` 在視圖模式表格與物件詳細頁提供一致的屬性值呈現 (#262)
- 結構化篩選模型 — 舊版篩選字串改為型別化的 `FilterRule` 物件，支援型別感知的運算子驗證 (#267)
- TUI Overlay 系統 — 說明覆蓋層與彈出視窗使用 lipgloss `Layer`/`Compositor` 正確分層 (#261, #276)

### 修正

- View 編輯器正確將變更轉發至視圖模式，並支援 Shift+Tab 導航
- View 編輯器不再顯示重複欄位
- 視圖模式佈局、CJK 文字對齊與物件詳細頁導航修正
- 側邊欄顯示物件名稱而非 slug

[v0.5.0]: https://github.com/typemd/typemd/releases/tag/v0.5.0

## [v0.4.0] - 2026-03-18

### 新增

- 內建 Page 型別 — `page`（📄）為新的內建型別，無需 YAML 檔案即可使用，適合自由格式內容 (#245)
- Vault 健康檢查 — `tmd doctor` 掃描孤立目錄、遺失型別等 vault 完整性問題 (#19)
- Vault 設定指令 — `tmd config get/set/list` 管理 `.typemd/config.yaml` 中的持久設定；`cli.default_type` 設定物件建立時的預設型別 (#241)
- 入門型別範本 — `tmd init` 提供可選的入門 schema（idea 💡、note 📝、book 📚）快速建立 vault (#235)
- 型別 Schema 版本號 — `version` 欄位（semver `"major.minor"`）追蹤 schema 演進 (#45)
- 型別 Schema 色彩 — `color` 欄位支援預設色名與 hex 色碼，用於視覺主題 (#228)
- 型別與屬性描述 — 型別 schema 與屬性新增 `description` 欄位，用於行內文件 (#228)
- TUI 範本管理 — 在型別編輯器中瀏覽、編輯、建立和刪除物件範本 (#250)
- TUI 物件建立精靈 — 標題面板內行內輸入，支援即時範本預覽與智慧 slug 轉換 (#229, #239)
- TUI 型別建立精靈 — 標題面板中的多欄位精靈（emoji、名稱、複數），含即時 schema 預覽 (#230)
- 彈性物件建立 — `tmd object create` 的型別參數現為選填（預設讀取 `cli.default_type`），名稱自動從自然語言轉換為 slug (#236, #240)
- Frontmatter 身份優先排序 — 系統屬性以固定順序排列：`name`、`description`、`created_at`、`updated_at`、`tags` (#199)

[v0.4.0]: https://github.com/typemd/typemd/releases/tag/v0.4.0

## [v0.3.0] - 2026-03-14

### 破壞性變更

- 移除內建型別 — `tmd init` 不再建立 `book`、`person`、`note`；請自行定義所需型別 (#208)
- 保留系統屬性 — `description`、`created_at`、`updated_at`、`tags` 現為保留名稱；型別 schema 若定義這些屬性名稱，驗證將會失敗。升級前請先移除 (#193, #201, #204)

### 新增

- 物件範本 — 在 `templates/<type>/` 放置 Markdown 檔案，建立物件時自動套用 frontmatter 預設值與正文內容；單一範本自動套用，多個範本提示選擇 (#173)
- 名稱範本 — 在型別 schema 的 `name` 屬性定義 `template`，自動產生物件名稱（例如 `日記 {{ date:YYYY-MM-DD }}`）(#186)
- 複數顯示名稱 — 型別 schema 新增 `plural` 欄位，在 TUI 中使用文法正確的集合名稱 (#205)
- 唯一性約束 — 型別 schema 設定 `unique: true`，防止同一型別中出現重複名稱 (#79)
- 標籤名稱驗證 — `tmd type validate` 新增全 vault 標籤名稱唯一性檢查 (#215)
- 系統屬性 — `description`、`created_at`、`updated_at`、`tags` 現為每個物件自動擁有的內建系統屬性 (#193, #201, #204)
- 內建標籤型別 — `tag` 為內建型別，sync 時若物件參考不存在的標籤會自動建立 (#204)
- TUI 型別編輯器 — 在 TUI 中直接進行型別 schema 的完整 CRUD：瀏覽、編輯、新增／刪除屬性、調整順序 (#207)
- 領域事件 — 實體操作產生領域事件（`ObjectCreated`、`ObjectSaved`、`PropertyChanged`、`ObjectLinked`、`TagAutoCreated`），為擴充性打下基礎 (#226)
- CQRS 架構 — core 重構為讀寫分離，寫入走 `ObjectService`、查詢走 `QueryService`，底層由 `ObjectRepository` 與 `ObjectIndex` 介面支撐 (#224)

### 修正

- TUI Emoji 對齊 — 修正含有 variation selector 的 emoji 寬度不一致問題

[v0.3.0]: https://github.com/typemd/typemd/releases/tag/v0.3.0

## [v0.2.0] - 2026-03-11

### 破壞性變更

- `name` 屬性 — 現為保留系統屬性；型別 schema 若手動定義 `name` 屬性，升級後驗證將會失敗。升級前請先移除型別 schema 中的 `name` 定義 (#187)

### 新增

- 屬性型別系統 — 在型別 schema 中定義 9 種屬性型別（`string`、`text`、`number`、`bool`、`date`、`datetime`、`url`、`enum`、`relation`）(#8)
- 共用屬性 — 在 `.typemd/properties.yaml` 定義可重用的屬性，並透過 `use` 在型別 schema 中參照 (#188)
- 型別 Emoji — 在型別 schema 加入 `emoji` 欄位，於 TUI 中視覺化識別型別 (#145)
- 屬性 Emoji — 在屬性 schema 加入 `emoji` 欄位，用於緊湊顯示 (#144)
- TUI 標題面板 — 瀏覽物件時顯示型別 emoji 與物件名稱的專用標題列 (#169)
- TUI 置頂屬性 — 在 schema 中標記 `pinned: true`，使屬性在 TUI 詳細檢視中突出顯示 (#168)
- TUI Session 持久化 — 游標位置、選取物件與面板狀態在 TUI 重新啟動後恢復 (#82)
- `--readonly` 旗標 — 以唯讀模式啟動 TUI，停用所有編輯功能 (#107)
- `--reindex` 旗標 — 全域旗標，啟動時強制重建 SQLite 索引，取代原本的 `tmd reindex` 子指令 (#159)
- 前綴比對 — 可用 ULID 後綴的短前綴解析物件，不需輸入完整 ID (#72)
- Homebrew 安裝 — 透過 `brew install typemd/tap/typemd-cli` 安裝 (#140)

### 變更

- `name` 屬性 — 現為必要系統屬性，自動從物件 slug 填入；型別 schema 不可自行定義名為 `name` 的屬性 (#187)
- TUI 物件列表 — 群組標頭中顯示型別 emoji (#163)
- 未定義屬性 — 型別 schema 未宣告的屬性在同步時會被靜默過濾 (#174)

### 修正

- Relation 顯示 — 移除 relation 屬性顯示值中的 ULID 後綴

[v0.2.0]: https://github.com/typemd/typemd/releases/tag/v0.2.0

## [v0.1.0] - 2026-03-08

### 新增

- 物件與型別 — 在 YAML 中定義型別 schema，透過 `tmd object create` 建立 Markdown 物件檔案 (#18)
- ULID 檔名 — 唯一的 ULID 後綴，避免物件命名衝突 (#48)
- Relation — 透過 `tmd relation link` / `tmd relation unlink` 建立雙向連結，支援單值覆寫與多值附加
- Wiki-links 與反向連結 — 在 Markdown 內文中使用 `[[target]]` 語法，自動追蹤反向連結 (#10)
- 查詢 — `tmd query` 依型別與屬性篩選，`tmd search` 全文搜尋，皆支援 `--json` 輸出
- 驗證 — `tmd type validate` 檢查 schema 完整性、屬性型別、孤立 relation 與壞掉的 wiki-links (#20)
- 遷移 — `tmd migrate` 在 schema 演進時更新既有物件 (#22)
- 自動重建索引 — SQLite 索引為空或遺失時自動重建 (#41)
- 孤立清理 — 重新索引時偵測並移除過期的 relation (#21)
- CLI 重組 — 指令依資源類型分組：`tmd object`、`tmd type`、`tmd relation` (#141)
- TUI — 三面板介面 (#47)、原地內文編輯 (#85)、編輯模式視覺指示 (#84)、退出時自動儲存 (#86)、快捷鍵說明 (#104)
- TUI 顯示 — 移除顯示名稱中的 ULID (#75)、縮減縮排 (#57)、群組化物件列表 (#43)
- MCP Server — `tmd mcp` 將 vault 開放給 AI 助手使用
- `.gitignore` 初始化 — `tmd init` 建立 `.typemd/.gitignore` 排除 `index.db` (#1)
- `tmd` 執行檔 — `go install` 產生 `tmd` binary (#61)
- 支援英文與繁體中文的文件網站 (#50, #54)
- 使用 Godog 與 Gherkin feature 檔案的 BDD 測試框架 (#111, #112)
- GitHub Actions 跨平台編譯發布流程 (#39)
- 程式碼重構 — 統一命名慣例、抽取 helper、改善錯誤處理 (#56)
- Vault 結構重構 — 移除 `objects/` 目錄層 (#117)

[v0.1.0]: https://github.com/typemd/typemd/releases/tag/v0.1.0
