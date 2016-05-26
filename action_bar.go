package action_bar

import (
	"bytes"
	"github.com/qor/admin"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ActionBar struct {
}

func (bar *ActionBar) Register(admin *admin.Admin) {
	router := admin.GetRouter()
	router.Get("/switch_mode", SwitchMode)
}

func (bar *ActionBar) RenderIncludedTag(Request *http.Request) template.HTML {
	var file string
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		file = path.Join(gopath, "src/github.com/qor/action_bar/views/themes/action_bar/assets/action_bar.tmpl")
	}

	var checked bool
	if cookie, err := Request.Cookie("qor-action-bar"); err == nil {
		checked = cookie.Value == "true"
	}

	var result = bytes.NewBufferString("")
	if tmpl, err := template.New(filepath.Base(file)).ParseFiles(file); err == nil {
		if err = tmpl.Execute(result, struct{ Checked bool }{Checked: checked}); err == nil {
			return template.HTML(result.String())
		}
	}
	return template.HTML("")
}
