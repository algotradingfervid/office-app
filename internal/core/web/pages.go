package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"

	"github.com/pocketbase/pocketbase/core"
)

// View is what every page template receives.
type View struct {
	AppName string
	Title   string
	User    *core.Record // nil when signed out
	Data    any          // the page's own data
}

// Pages holds one module's page templates, each parsed together with the shared layout.
type Pages struct {
	set map[string]*template.Template
}

// NewPages parses every templates/*.html in fsys. A page defines "title" and "content";
// the layout defines "base". Page names are file names without ".html".
func NewPages(fsys fs.FS) *Pages {
	files, err := fs.Glob(fsys, "templates/*.html")
	if err != nil || len(files) == 0 {
		panic(fmt.Sprintf("web: no templates found: %v", err))
	}
	p := &Pages{set: map[string]*template.Template{}}
	for _, f := range files {
		t := template.Must(template.New("").ParseFS(layoutFS, "templates/*.html"))
		template.Must(t.ParseFS(fsys, f))
		p.set[path.Base(f[:len(f)-len(".html")])] = t
	}
	return p
}

// Render writes page name inside the layout with status 200.
func (p *Pages) Render(e *core.RequestEvent, name, title string, data any) error {
	return p.RenderStatus(e, http.StatusOK, name, title, data)
}

// RenderStatus writes page name inside the layout with the given status.
func (p *Pages) RenderStatus(e *core.RequestEvent, status int, name, title string, data any) error {
	t, ok := p.set[name]
	if !ok {
		return fmt.Errorf("web: unknown page %q", name)
	}
	var buf bytes.Buffer
	v := View{AppName: AppName, Title: title, User: e.Auth, Data: data}
	if err := t.ExecuteTemplate(&buf, "base", v); err != nil {
		return err
	}
	return e.HTML(status, buf.String())
}
