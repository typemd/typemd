// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://docs.typemd.io',
	integrations: [
		starlight({
			title: 'TypeMD Docs',
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
					label: 'Guides',
					translations: { 'zh-TW': '指南' },
					items: [
						{ label: 'Objects & Types', slug: 'guides/objects-and-types' },
						{ label: 'Relations', slug: 'guides/relations' },
						{ label: 'Type Schemas', slug: 'guides/type-schemas' },
						{ label: 'Querying', slug: 'guides/querying' },
						{ label: 'Translating', slug: 'guides/translating' },
					],
				},
				{
					label: 'CLI Reference',
					translations: { 'zh-TW': 'CLI 參考' },
					autogenerate: { directory: 'reference' },
				},
				{
					label: 'Architecture',
					translations: { 'zh-TW': '架構' },
					items: [
						{ label: 'Data Model', slug: 'architecture/data-model' },
					],
				},
			],
		}),
	],
});
