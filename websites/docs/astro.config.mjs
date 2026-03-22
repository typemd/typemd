// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
// https://astro.build/config
export default defineConfig({
	site: 'https://docs.typemd.io',
	integrations: [
		starlight({
			title: 'TypeMD Docs',
			head: [
				{
					tag: 'script',
					attrs: { type: 'module' },
					content: `
						import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs';
						document.querySelectorAll('code.language-mermaid').forEach(el => {
							const pre = el.parentElement;
							pre.classList.add('mermaid');
							pre.textContent = el.textContent;
						});
						mermaid.initialize({ startOnLoad: true, theme: 'neutral' });
					`,
				},
			],
			defaultLocale: 'root',
			locales: {
				root: {
					label: 'English',
					lang: 'en',
				},
				'zh-tw': {
					label: '繁體中文',
					lang: 'zh-TW',
				},
			},
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/typemd/typemd' },
				{ icon: 'pencil', label: 'Blog', href: 'https://blog.typemd.io' },
			],
			sidebar: [
				{
					label: 'Getting Started',
					translations: { 'zh-TW': '快速入門' },
					items: [
						{ label: 'Introduction', slug: 'getting-started/introduction' },
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Quick Start', slug: 'getting-started/quick-start' },
						{ label: 'AI Setup', slug: 'getting-started/ai-setup' },
					],
				},
				{
					label: 'Concepts',
					translations: { 'zh-TW': '核心概念' },
					items: [
						{ label: 'Overview', slug: 'concepts/overview' },
						{ label: 'Objects', slug: 'concepts/objects' },
						{ label: 'Types', slug: 'concepts/types' },
						{ label: 'Relations', slug: 'concepts/relations' },
						{ label: 'Links', slug: 'concepts/links' },
						{ label: 'Glossary', slug: 'concepts/glossary' },
					],
				},
				{
					label: 'Basics',
					translations: { 'zh-TW': '基本功能' },
					items: [
						{ label: 'Properties', slug: 'basics/properties' },
						{ label: 'Tags', slug: 'basics/tags' },
						{ label: 'Templates', slug: 'basics/templates' },
						{ label: 'Search', slug: 'basics/search' },
						{ label: 'Queries', slug: 'basics/queries' },
						{ label: 'Validation', slug: 'basics/validation' },
						{ label: 'Configuration', slug: 'basics/configuration' },
					],
				},
				{
					label: 'Advanced',
					translations: { 'zh-TW': '進階' },
					items: [
						{ label: 'File Structure', slug: 'advanced/file-structure' },
						{ label: 'Frontmatter', slug: 'advanced/frontmatter' },
					],
				},
				{
					label: 'TUI',
					translations: { 'zh-TW': 'TUI' },
					autogenerate: { directory: 'tui' },
				},
				{
					label: 'CLI',
					translations: { 'zh-TW': 'CLI' },
					autogenerate: { directory: 'cli' },
				},
				{
					label: 'MCP',
					translations: { 'zh-TW': 'MCP' },
					autogenerate: { directory: 'mcp' },
				},
				{
					label: 'Developers',
					translations: { 'zh-TW': '開發者' },
					items: [
						{ label: 'Data Model', slug: 'developers/data-model' },
						{ label: 'Architecture', slug: 'developers/architecture' },
						{ label: 'Design Decisions', slug: 'developers/design-decisions' },
						{ label: 'Internals', slug: 'developers/internals' },
					],
				},
			],
		}),
	],
});
