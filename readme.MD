# go-sveltespa

Just a simple router for svelte spa apps. has fallbacks. Embedded and regular file system. 

Redirects trailing slash routes as these gave me trouble. 

## ⚠️ Warning
When using the embedded filesystem with Svelte applications, make sure to use the `//go:embed all:build` directive instead of just `//go:embed build`. Svelte's build process creates files with leading underscores (e.g., `_app/chunks/_23fij2.js`), and Go's default embed behavior ignores files that start with underscores. Using `all:build` ensures all files are included.