package action_bar

import (
	"github.com/qor/admin"
	"html/template"
	"net/http"
)

// ActionBar stores configuration about a action bar.
type ActionBar struct {
	admin   *admin.Admin
	auth    admin.Auth
	Actions []*Action
}

// Action define a addition action(link), will append to the top-right menu.
type Action struct {
	Name string
	Link string
}

func init() {
	admin.RegisterViewPath("github.com/qor/action_bar/views")
}

// New will create a ActionBar object
func New(admin *admin.Admin, auth admin.Auth) *ActionBar {
	bar := &ActionBar{admin: admin, auth: auth}
	admin.GetRouter().Get("/action_bar/switch_mode", SwitchMode)
	return bar
}

// RegisterAction registered a new action
func (bar *ActionBar) RegisterAction(action *Action) {
	bar.Actions = append(bar.Actions, action)
}

// Render will return the HTML of the bar, used this function to render the bar in frontend page's template or layout
func (bar *ActionBar) Render(w http.ResponseWriter, r *http.Request) template.HTML {
	context := bar.admin.NewContext(w, r)
	result := map[string]interface{}{
		"EditMode":     bar.EditMode(w, r),
		"Auth":         bar.auth,
		"CurrentUser":  bar.auth.GetCurrentUser(context),
		"Actions":      bar.Actions,
		"RouterPrefix": bar.admin.GetRouter().Prefix,
	}
	return context.Render("action_bar", result)
}

// EditMode return whether current mode is `Preview` or `Edit`
func (bar *ActionBar) EditMode(w http.ResponseWriter, r *http.Request) bool {
	context := bar.admin.NewContext(w, r)
	if bar.auth.GetCurrentUser(context) == nil {
		return false
	}
	if cookie, err := r.Cookie("qor-action-bar"); err == nil {
		return cookie.Value == "true"
	}
	return false
}
