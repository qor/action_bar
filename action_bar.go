package action_bar

import (
	"bytes"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ActionBar struct {
	auth  admin.Auth
	admin *admin.Admin
}

func (bar *ActionBar) Register(admin *admin.Admin) {
	bar.admin = admin
	router := admin.GetRouter()
	router.Get("/switch_mode", SwitchMode)
}

func (bar *ActionBar) SetAuth(auth admin.Auth) {
	bar.auth = auth
}

func (bar *ActionBar) RenderIncludedTag(w http.ResponseWriter, r *http.Request) template.HTML {
	var file string
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		file = path.Join(gopath, "src/github.com/qor/action_bar/views/themes/action_bar/assets/action_bar.tmpl")
	}

	var checked bool
	if cookie, err := r.Cookie("qor-action-bar"); err == nil {
		checked = cookie.Value == "true"
	}

	var result = bytes.NewBufferString("")
	context := bar.admin.NewContext(w, r)
	if tmpl, err := template.New(filepath.Base(file)).ParseFiles(file); err == nil {
		context := struct {
			Checked     bool
			Auth        admin.Auth
			Context     *admin.Context
			CurrentUser qor.CurrentUser
		}{
			Checked:     checked,
			Auth:        bar.auth,
			Context:     context,
			CurrentUser: bar.auth.GetCurrentUser(context),
		}
		if err = tmpl.Execute(result, context); err == nil {
			return template.HTML(result.String())
		}
	}
	return template.HTML("")
}
