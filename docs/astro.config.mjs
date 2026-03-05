// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://docs.typemd.io',
	integrations: [
		starlight({
			title: 'TypeMD Docs',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/typemd/typemd' }],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Introduction', slug: 'getting-started/introduction' },
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Quick Start', slug: 'getting-started/quick-start' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'Objects & Types', slug: 'guides/objects-and-types' },
						{ label: 'Relations', slug: 'guides/relations' },
						{ label: 'Type Schemas', slug: 'guides/type-schemas' },
						{ label: 'Querying', slug: 'guides/querying' },
					],
				},
				{
					label: 'CLI Reference',
					autogenerate: { directory: 'reference' },
				},
				{
					label: 'Architecture',
					items: [
						{ label: 'Data Model', slug: 'architecture/data-model' },
					],
				},
			],
		}),
	],
});
