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
	admin   *admin.Admin
	auth    admin.Auth
	Actions []*Action
}

type Action struct {
	Name string
	Link string
}

var root = "."

func init() {
	if path := os.Getenv("WEB_ROOT"); path != "" {
		root = path
	}
}

func New(admin *admin.Admin, auth admin.Auth) *ActionBar {
	bar := &ActionBar{admin: admin, auth: auth}
	router := admin.GetRouter()
	router.Get("/switch_mode", SwitchMode)
	admin.RegisterViewPath("github.com/qor/action_bar/views")
	return bar
}

func (bar *ActionBar) RegisterAction(action *Action) {
	bar.Actions = append(bar.Actions, action)
}

func (bar *ActionBar) RenderIncludedTag(w http.ResponseWriter, r *http.Request) template.HTML {
	var result = bytes.NewBufferString("")
	context := bar.admin.NewContext(w, r)
	file := appendPackageToPath("views/themes/action_bar/action_bar.tmpl")
	if tmpl, err := template.New(filepath.Base(file)).ParseFiles(file); err == nil {
		context := struct {
			Checked      bool
			Auth         admin.Auth
			Context      *admin.Context
			CurrentUser  qor.CurrentUser
			Actions      []*Action
			RouterPrefix string
		}{
			Checked:      bar.IsChecked(w, r),
			Auth:         bar.auth,
			Context:      context,
			CurrentUser:  bar.auth.GetCurrentUser(context),
			Actions:      bar.Actions,
			RouterPrefix: bar.admin.GetRouter().Prefix,
		}
		if err = tmpl.Execute(result, context); err == nil {
			return template.HTML(result.String())
		}
	}
	return template.HTML("")
}

func (bar *ActionBar) IsChecked(w http.ResponseWriter, r *http.Request) bool {
	context := bar.admin.NewContext(w, r)
	if bar.auth.GetCurrentUser(context) == nil {
		return false
	}
	if cookie, err := r.Cookie("qor-action-bar"); err == nil {
		return cookie.Value == "true"
	}
	return false
}

func appendPackageToPath(f string) string {
	var file string

	file = filepath.Join(root, "vendor", "github.com/qor/action_bar/"+f)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
			file = path.Join(gopath, "src/github.com/qor/action_bar/"+f)
		}
	}
	return file
}
