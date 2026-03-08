---
title: tmd relation unlink
description: 移除兩個 Object 之間的 Relation。
sidebar:
  order: 7
---

移除 Relation。使用 `--both` 來同時移除反向端。

```bash
tmd relation unlink book/golang-in-action author person/alan-donovan --both
```

## 選項

| 選項 | 說明 |
|------|------|
| `--both` | 同時移除 Relation 的反向端 |
