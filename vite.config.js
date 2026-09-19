// For the most part we use Vite defaults. Mainly the source and output directories are a little different to support
// integration with the Go build.
export default {
  appType: 'custom',
  // Vite 8: declare the entry at the top level rather than under build.rolldownOptions.input, since the
  // top-level option is also picked up by the dev server's dependency optimizer, not just production builds.
  input: '/js/main.js',
  build: {
    outDir: 'build/dist',
    // Since we're already outputting to a build subdirectory and we want to embed the manifest file, output to
    // 'manifest.json' rather than the default '.vite/manifest.json'
    manifest: 'manifest.json',
  },
  css: {
    preprocessorOptions: {
      scss: {
        // NOTE(dk 2025/02/20): The latest Sass noisily flags the latest Bootstrap for some deprecated features.
        // For the time being, hide that noise with these.
        silenceDeprecations: ['mixed-decls', 'color-functions', 'global-builtin', 'import'],
      }
    }
  },
}
