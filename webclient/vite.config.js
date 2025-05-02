import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig(({ mode }) => {
	const isProduction = mode === 'production';

	return {
		define: {
			'process.env.PUBLIC_DEV_API': JSON.stringify(
				isProduction
					? undefined
					: 'http://localhost:3000/logger'
			),
		},
		plugins: [sveltekit()]
	};
});
