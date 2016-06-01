package action_bar

import (
	"html/template"
	"net/http"

	"github.com/qor/admin"
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

func init() {
	admin.RegisterViewPath("github.com/qor/action_bar/views")
}

func New(admin *admin.Admin, auth admin.Auth) *ActionBar {
	bar := &ActionBar{admin: admin, auth: auth}
	admin.GetRouter().Get("/action_bar/switch_mode", SwitchMode)
	return bar
}

func (bar *ActionBar) RegisterAction(action *Action) {
	bar.Actions = append(bar.Actions, action)
}

func (bar *ActionBar) Render(w http.ResponseWriter, r *http.Request) template.HTML {
	context := bar.admin.NewContext(w, r)
	result := map[string]interface{}{
		"Checked":      bar.IsChecked(w, r),
		"Auth":         bar.auth,
		"Context":      context,
		"CurrentUser":  bar.auth.GetCurrentUser(context),
		"Actions":      bar.Actions,
		"RouterPrefix": bar.admin.GetRouter().Prefix,
	}
	return context.Render("action_bar", result)
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
