package action_bar

import (
	"html/template"
	"net/http"

	"github.com/qor/admin"
	"github.com/qor/qor/utils"
)

// SwitchMode is handle to store switch status in cookie
func SwitchMode(context *admin.Context) {
	utils.SetCookie(http.Cookie{Name: "qor-action-bar", Value: context.Request.URL.Query().Get("checked")}, context.Context)

	referrer := context.Request.Referer()
	if referrer == "" {
		referrer = "/"
	}

	http.Redirect(context.Writer, context.Request, referrer, http.StatusFound)
}

func (bar *ActionBar) FuncMap() template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_edit_button"] = func(widgetName string, widgetGroupName ...string) template.HTML {
		return template.HTML(`<a target="blank" href="/admin/widget_contents/FeatureProducts/edit?widget_scope=default" class="qor-actionbar-button">Edit</a>`)
	}

	return funcMap
}
