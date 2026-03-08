# 變更日誌

本檔案記錄本專案所有值得注意的變更。

格式依循 [Keep a Changelog](https://keepachangelog.com/)。

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
