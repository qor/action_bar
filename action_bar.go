package action_bar

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/qor/admin"
	"github.com/qor/qor/utils"
)

// ActionBar stores configuration about a action bar.
type ActionBar struct {
	Admin   *admin.Admin
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
	bar := &ActionBar{Admin: admin}
	ctr := &controller{ActionBar: bar}
	admin.GetRouter().Get("/action_bar/switch_mode", ctr.SwitchMode)
	admin.GetRouter().Get("/action_bar/inline_edit", ctr.InlineEdit)
	return bar
}

// RegisterAction registered a new action
func (bar *ActionBar) RegisterAction(action *Action) {
	bar.Actions = append(bar.Actions, action)
}

// Render will return the HTML of the bar, used this function to render the bar in frontend page's template or layout
func (bar *ActionBar) Render(w http.ResponseWriter, r *http.Request, configs ...Config) template.HTML {
	context := bar.Admin.NewContext(w, r)
	result := map[string]interface{}{
		"EditMode":     bar.EditMode(w, r),
		"Auth":         bar.Admin.Auth,
		"CurrentUser":  bar.Admin.Auth.GetCurrentUser(context),
		"Actions":      bar.Actions,
		"RouterPrefix": bar.Admin.GetRouter().Prefix,
	}
	if len(configs) > 0 {
		result["InlineActions"] = configs[0].InlineActions
	}
	return context.Render("action_bar/action_bar", result)
}

// EditMode return whether current mode is `Preview` or `Edit`
func (bar *ActionBar) EditMode(w http.ResponseWriter, r *http.Request) bool {
	context := bar.Admin.NewContext(w, r)
	if bar.Admin.Auth.GetCurrentUser(context) == nil {
		return false
	}
	if cookie, err := r.Cookie("qor-action-bar"); err == nil {
		return cookie.Value == "true"
	}
	return false
}

// FuncMap will return helper to render inline edit button
func (bar *ActionBar) FuncMap(w http.ResponseWriter, r *http.Request) template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_edit_button"] = func(value interface{}, resources ...*admin.Resource) template.HTML {
		return bar.RenderEditButtonWithResource(w, r, value, resources...)
	}

	return funcMap
}

func (bar *ActionBar) RenderEditButtonWithResource(w http.ResponseWriter, r *http.Request, value interface{}, resources ...*admin.Resource) template.HTML {
	var (
		admin        = bar.Admin
		context      = admin.NewContext(w, r)
		resourceName = utils.ModelType(value).String()
		editURL, _   = utils.JoinURL(context.URLFor(value, resources...), "edit")
	)

	if res := admin.GetResource(resourceName); res != nil {
		resourceName = string(admin.T(context.Context, fmt.Sprintf("%v.name", res.ToParam()), res.Name))
	}

	for _, res := range resources {
		resourceName = string(admin.T(context.Context, fmt.Sprintf("%v.name", res.ToParam()), res.Name))
	}

	title := string(admin.T(context.Context, "qor_action_bar.action.edit_resource", "Edit {{$1}}", resourceName))
	return bar.RenderEditButton(w, r, title, editURL)
}

func (bar *ActionBar) RenderEditButton(w http.ResponseWriter, r *http.Request, title string, link string) template.HTML {
	if bar.EditMode(w, r) {
		var (
			prefix   = bar.Admin.GetRouter().Prefix
			jsURL    = fmt.Sprintf("<script data-prefix=\"%v\" src=\"%v/assets/javascripts/action_bar_check.js?theme=action_bar\"></script>", prefix, prefix)
			frameURL = fmt.Sprintf("%v/action_bar/inline_edit", prefix)
		)

		return template.HTML(fmt.Sprintf(`%v<a target="blank" data-iframe-url="%v" data-url="%v" href="#" class="qor-actionbar-button">%v</a>`, jsURL, frameURL, link, title))
	}
	return template.HTML("")
}
