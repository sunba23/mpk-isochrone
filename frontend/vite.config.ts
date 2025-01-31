import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { nodePolyfills } from "vite-plugin-node-polyfills";
import path from "path";

export default defineConfig({
  server: {
    proxy: {
      '/stops': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path
      }
    }
  },
  plugins: [
    react(),
    nodePolyfills({
      protocolImports: true,
    }),
  ],
  define: {
    "process.env": {},
  },
  build: {
    outDir: "dist",
    assetsDir: "assets",
    jsx: "react",
  },
});
