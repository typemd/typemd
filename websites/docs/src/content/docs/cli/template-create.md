---
title: tmd template create
description: Create a new template file.
sidebar:
  order: 9.8
---

Creates a new template file and opens it in your editor.

```bash
tmd template create book/review
tmd template create note/meeting
```

The argument must be in `type/name` format. The command:

1. Creates the file at `templates/<type>/<name>.md` (creating the type directory if needed)
2. Opens the file in your editor (`$EDITOR`, `$VISUAL`, or `vi`)

If the template already exists, the command returns an error.
