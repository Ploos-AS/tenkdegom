import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://tenkdegom.no',
  build: {
    inlineStylesheets: 'never',
  },
  vite: {
    build: {
      assetsInlineLimit: 0,
    },
  },
});
