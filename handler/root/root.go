package root

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/russross/blackfriday"
	"github.com/unders/docit/cli"
	"github.com/unders/docit/template"
)

const notFoundMsg = "<pre style='word-wrap: break-word;" +
	"white-space: pre-wrap;'>404 page not found</pre>"

// Handle renders pages unders "/"
func Handle(arg cli.Arg) func(w http.ResponseWriter, req *http.Request) {
	fileServer := http.FileServer(http.Dir(arg.Root))

	return func(w http.ResponseWriter, req *http.Request) {
		upath := req.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			req.URL.Path = upath
		}

		// Redirect to index.md page.
		if req.URL.Path == "/" {
			http.Redirect(w, req, "/"+arg.Index, http.StatusSeeOther)
			return
		}

		// Redirect to help.md page.
		if req.URL.Path == "/help" {
			http.Redirect(w, req, "/"+arg.Help, http.StatusSeeOther)
			return
		}

		// Parse and serve given Markdown file relative to root dir.
		if strings.HasSuffix(req.URL.Path, ".md") {
			html, code := loadPage(arg.Root + req.URL.Path)
			template.Render(w, html, code)
			return
		}

		// Serve file relative to root dir.
		fileServer.ServeHTTP(w, req)
	}
}

func loadPage(filename string) ([]byte, int) {
	markdown, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Markdown file %s not found. Error: %v\n", filename, err)
		return []byte(notFoundMsg), http.StatusNotFound
	}

	html := blackfriday.MarkdownCommon(markdown)

	return html, http.StatusOK
}