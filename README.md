# Dev server for React

Standalone esbuild based app for building, watching and serving React projects.  
Written on Go, it uses net/http module to serve web project and Go implementation of esbuild.  
However, to get maximum performance and best stability, standalone original compilers are used for [sass](https://sass-lang.com/dart-sass) (Dart) and [less](https://github.com/rus-sharafiev/less-compiler) (ES6 source code wrapped inside Deno)

## Install
Using PowerShell (Windows x64 only)
```
irm https://github.com/rus-sharafiev/dev-server/releases/latest/download/install.ps1 | iex
```
## Usage

Start development server, build project and watch for changes
```powershell
dev start
```

Build minified production build
```powershell
dev build
```
