// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import mermaid from 'astro-mermaid';

// https://astro.build/config
export default defineConfig({
	site: 'https://apollogeddon.github.io',
	base: '/ignition-tfpl',
	integrations: [
		starlight({
			title: 'Ignition TF Plugin',
			// customCss: ['./src/styles/custom.css'],
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/apollogeddon/ignition-tfpl' },
			],
			sidebar: [
				{
					label: 'Guides',
					items: [
						{ label: 'Installation', link: '/guides/installation/' },
						{ label: 'Architecture', link: '/guides/architecture/' },
						{ label: 'Capabilities', link: '/guides/capabilities/' },
					],
				},
				{
					label: 'Reference',
					autogenerate: { directory: 'reference' },
				},
			],
		}),
		mermaid(),
	],
});