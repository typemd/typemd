---
title: tmd object list
description: 列出所有 Object。
sidebar:
  order: 4
---

列出 vault 中所有已建立的 Object。

```bash
tmd object list
```

輸出範例：

```
book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz
book/clean-code-01jqr4b7npcwq9f3h8k2m5vxya
person/robert-martin-01jqr4c8mqdzr0g4j9l3n6wyab
```

## Flags

| Flag | 說明 |
|------|------|
| `--json` | 以 JSON 格式輸出結果 |
| `--include-archived` | 在結果中包含已封存的 object |
