---
title: tmd log
description: 顯示特定 Object 的 git 提交歷史。
sidebar:
  order: 11.6
---

顯示特定 Object 檔案的 git 提交歷史。使用 `git log --follow` 來追蹤跨重新命名的歷史紀錄。

## 用法

```bash
tmd log <object-id>
```

支援前綴比對——可以省略 ULID 後綴，只要前綴能唯一識別一個 Object。如果前綴匹配多個 Object，會顯示互動式選擇器。

### 輸出範例

```
commit a1b2c3d (HEAD -> main)
Author: Alice <alice@example.com>
Date:   Sat Mar 29 2026

    update book properties

commit e4f5g6h
Author: Alice <alice@example.com>
Date:   Fri Mar 28 2026

    add clean-code book
```

## 精簡輸出

```bash
tmd log --oneline book/clean-code
```

每個 commit 顯示為一行，包含縮短的 hash 和標題。

## 邊界情況

- **Vault 不在 git repository 內** — 報錯：「vault is not inside a git repository」
- **Object 沒有任何 commit** — 顯示：「no commits found for this object」
- **已重新命名的 Object** — `--follow` 會追蹤跨檔案重新命名的歷史
