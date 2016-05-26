package action_bar

import (
	"bytes"
	"github.com/qor/admin"
	"html/template"
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

func (bar *ActionBar) RenderIncludedTag() template.HTML {
	var file string
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		file = path.Join(gopath, "src/github.com/qor/action_bar/views/themes/action_bar/assets/action_bar.tmpl")
	}

	var result = bytes.NewBufferString("")
	if tmpl, err := template.New(filepath.Base(file)).ParseFiles(file); err == nil {
		if err = tmpl.Execute(result, nil); err == nil {
			return template.HTML(result.String())
		}
	}
	return template.HTML("")
}
