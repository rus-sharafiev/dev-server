# Dev server for React

The standalone application based on the [esbuild](https://esbuild.github.io/) for building, watching, serving and deploying React projects.  
Written on Go, it uses net/http module to serve web project and Go implementation of the esbuild, but, in order to get maximum performance and better stability, standalone original compilers are used for processing [Sass](https://sass-lang.com/dart-sass) and [Less](https://github.com/rus-sharafiev/less-compiler) files.

## Install
Using PowerShell (Windows x64 only)
```
irm https://github.com/rus-sharafiev/dev-server/releases/latest/download/install.ps1 | iex
```
## Usage

Build the project, start the development server and watch for the changes
```powershell
dev start
```

Create a minified production build
```powershell
dev build
```

Create a minified production build and start the web server
```powershell
dev serve
```

Build the project and then copy the resulting files to the remote server via scp, specifying the path in the configuration file named `dev.conf`...

```json
{
    "deployPath": "root@1.1.1.1:/var/www/html/",
    "jsPath": "root@1.1.1.1:/var/www/html/js/",
    "cssPath": "root@1.1.1.1:/var/www/html/css/"
}
```
...and then
```bash
dev deploy
```

`deployPath` - a path to copy whole build dir content  
`jsPath` - a path to copy .js files only  
`cssPath` - a path to copy .css files only  

A `deployPath` or both `jsPath` and `cssPath` should be provided!  
If all fiels present, then `deployPath` will be ignored
