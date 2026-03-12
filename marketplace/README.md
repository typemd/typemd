# typemd Marketplace

Community plugins for [typemd](https://github.com/typemd/typemd) — a local-first CLI knowledge management tool.

This marketplace uses [Claude Code's native plugin system](https://code.claude.com/docs/en/plugins). Plugins extend Claude Code with skills, commands, and agents designed for typemd workflows.

## Installation

Add the marketplace to Claude Code:

```
/plugin marketplace add typemd/marketplace
```

Browse and install plugins:

```
/plugin install <plugin-name>@typemd-marketplace
```

## Available Plugins

| Plugin | Description |
|--------|-------------|
| `markdown-import` | Convert existing markdown files into typemd objects with frontmatter, type detection, and relation discovery |
| `vault-guide` | typemd reference guide — teaches Claude how to work with vaults, CLI commands, type schemas, and objects |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to submit your own plugin.

## License

MIT
