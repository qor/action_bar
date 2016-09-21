package action_bar

import (
	"fmt"
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

// FuncMap will return helper to render inline edit button
func (bar *ActionBar) FuncMap(w http.ResponseWriter, r *http.Request) template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_edit_button"] = func(value interface{}, resources ...*admin.Resource) template.HTML {
		if bar.EditMode(w, r) {
			context := bar.admin.NewContext(nil, nil)
			url := context.URLFor(value, resources...)
			return template.HTML(fmt.Sprintf(`<a target="blank" href="%v/edit" class="qor-actionbar-button">Edit</a>`, url))
		}
		return template.HTML("")
	}

	return funcMap
}
