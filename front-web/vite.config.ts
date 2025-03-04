import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Remova ou atualize esta configuração se estiver causando problemas
  // css: {
  //   postcss: './postcss.config.js',
  // },
});
