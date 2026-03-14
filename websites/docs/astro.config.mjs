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
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/typemd/typemd' }],
			sidebar: [
				{
					label: 'Getting Started',
					translations: { 'zh-TW': '快速入門' },
					items: [
						{ label: 'Introduction', slug: 'getting-started/introduction' },
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Quick Start', slug: 'getting-started/quick-start' },
					],
				},
				{
					label: 'Concepts',
					translations: { 'zh-TW': '設計概念' },
					items: [
						{ label: 'Overview', slug: 'concepts/overview' },
						{ label: 'Objects', slug: 'concepts/objects' },
						{ label: 'Types', slug: 'concepts/types' },
						{ label: 'Relations', slug: 'concepts/relations' },
						{ label: 'Wiki-links', slug: 'concepts/wiki-links' },
						{ label: 'Data Model', slug: 'concepts/data-model' },
					    { label: 'Architecture', slug: 'concepts/architecture' },
					],
				},
				{
					label: 'CLI Reference',
					translations: { 'zh-TW': 'CLI 參考' },
					autogenerate: { directory: 'reference' },
				},
			],
		}),
	],
});
