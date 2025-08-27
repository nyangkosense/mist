/* mist - minimal static site generator */ 
/* see LICENSE for copyright and license details */

package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
)

type Page struct {
	Title   string
	Content string
	Nav     []string
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getTitle(path string) string {
	name := filepath.Base(path)
	return strings.TrimSuffix(name, ".md")
}

func mkNav(dir string) []string {
	files, err := os.ReadDir(dir)
	die(err)
	var nav []string
	for _, f := range files {
	if strings.HasSuffix(f.Name(), ".md") && f.Name() != "index.md" {
	name := strings.TrimSuffix(f.Name(), ".md")
	nav = append(nav, name)
		}
	}
	return nav
}

func copy(src, dst string) {
	in, err := os.Open(src)
	die(err)
	defer in.Close()
	out, err := os.Create(dst)
	die(err)
	defer out.Close()
	io.Copy(out, in)
}

func process(src, dst, tpl string) {
	os.MkdirAll(dst, 0755)
	
	copy("style.css", filepath.Join(dst, "style.css"))
	copy(tpl, filepath.Join(dst, "template.html"))
	files, err := os.ReadDir(src)
	die(err)	
	t, err := template.ParseFiles(tpl)
	die(err)
	nav := mkNav(src)
	
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			md, err := os.ReadFile(filepath.Join(src, f.Name()))
			die(err)
			
			page := Page{
				Title:   getTitle(filepath.Join(src, f.Name())),
				Content: string(markdown.ToHTML(md, nil, nil)),
				Nav:     nav,
			}
			
			name := strings.TrimSuffix(f.Name(), ".md")
			if name == "index" {
				name = "index.html"
			} else {
				name += ".html"
			}
			
			out, err := os.Create(filepath.Join(dst, name))
			die(err)
			die(t.Execute(out, page))
			out.Close()
		} else {
			copy(filepath.Join(src, f.Name()), filepath.Join(dst, f.Name()))
		}
	}
}

func serve(dir, port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		} else if !strings.HasSuffix(path, ".html") && !strings.Contains(path, ".") {
			if _, err := os.Stat(filepath.Join(dir, path[1:]+".html")); err == nil {
				path += ".html"
			}
		}
		http.ServeFile(w, r, filepath.Join(dir, path[1:]))
	})
	log.Printf("serving %s on :%s", dir, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	src := flag.String("i", ".", "input dir")
	dst := flag.String("o", "out", "output dir") 
	tpl := flag.String("t", "page.tmpl", "template")
	port := flag.String("s", "", "serve on port")
	flag.Parse()
	
	if *port != "" {
		serve(*dst, *port)
	} else {
		process(*src, *dst, *tpl)
	}
}
