mist
====
*mi*ni *st*atic site generator written in Go.

mist generates static HTML pages from markdown files in a flat directory
structure. A simple, more lightweight solution for basic static sites.

mist includes a built-in web server.

Usage
-----
Generate static site:
	go run mist.go [-i input_dir] [-o output_dir] [-t template_file]

Serve generated site:
	go run mist.go -s port_number

Options:
	-i string   input directory (default ".")
	-o string   output directory (default "out")  
	-t string   template file (default "page.tmpl")
	-s string   serve on port (-s 8080) 

Template Variables
------------------
provides the following template variables:

	.Title   string       // page title (filename without .md extension)
	.Content string       // HTML content from markdown
	.Nav     []string     // array of other markdown files in directory

Example template usage:
	<h1>{{.Title}}</h1>
	<div>{{.Content}}</div>
	<nav>
	{{range .Nav}}
		<a href="{{.}}.html">{{.}}</a>
	{{end}}
	</nav>

Features
--------
- Flat directory structure
- Built-in development server
- Simple Go template system
- Minimal dependencies
