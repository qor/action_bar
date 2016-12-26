package action_bar

import (
	"html/template"
	"net/http"

	"github.com/qor/admin"
)

// ActionBar stores configuration about a action bar.
type ActionBar struct {
	admin   *admin.Admin
	Actions []*Action
}

// Config stores configuration for a render
type Config struct {
	// some inline edit actions that will placed on the bar
	InlineActions []template.HTML
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
func New(admin *admin.Admin) *ActionBar {
	bar := &ActionBar{admin: admin}
	admin.GetRouter().Get("/action_bar/switch_mode", SwitchMode)
	admin.GetRouter().Get("/action_bar/inline_edit", InlineEdit)
	return bar
}

// RegisterAction registered a new action
func (bar *ActionBar) RegisterAction(action *Action) {
	bar.Actions = append(bar.Actions, action)
}

// Render will return the HTML of the bar, used this function to render the bar in frontend page's template or layout
func (bar *ActionBar) Render(w http.ResponseWriter, r *http.Request, configs ...Config) template.HTML {
	context := bar.admin.NewContext(w, r)
	result := map[string]interface{}{
		"EditMode":     bar.EditMode(w, r),
		"Auth":         bar.admin.Auth,
		"CurrentUser":  bar.admin.Auth.GetCurrentUser(context),
		"Actions":      bar.Actions,
		"RouterPrefix": bar.admin.GetRouter().Prefix,
	}
	if len(configs) > 0 {
		result["InlineActions"] = configs[0].InlineActions
	}
	return context.Render("action_bar/action_bar", result)
}

// EditMode return whether current mode is `Preview` or `Edit`
func (bar *ActionBar) EditMode(w http.ResponseWriter, r *http.Request) bool {
	context := bar.admin.NewContext(w, r)
	if bar.admin.Auth.GetCurrentUser(context) == nil {
		return false
	}
	if cookie, err := r.Cookie("qor-action-bar"); err == nil {
		return cookie.Value == "true"
	}
	return false
}
