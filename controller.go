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

// InlineEdit using to make inline edit resource shown as slideout
func InlineEdit(context *admin.Context) {
	context.Writer.Write([]byte(context.Render("inline_edit")))
}

// FuncMap will return helper to render inline edit button
func (bar *ActionBar) FuncMap(w http.ResponseWriter, r *http.Request) template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_edit_button"] = func(value interface{}, resources ...*admin.Resource) template.HTML {
		if bar.EditMode(w, r) {
			context := bar.admin.NewContext(nil, nil)
			url := context.URLFor(value, resources...)
			prefix := bar.admin.GetRouter().Prefix
			js := fmt.Sprintf("<script data-prefix=\"%v\" src=\"%v/assets/javascripts/action_bar_check.js?theme=action_bar\"></script>", prefix, prefix)
			frameURL := fmt.Sprintf("%v/action_bar/inline-edit", prefix)
			return template.HTML(fmt.Sprintf(`%v<a target="blank" data-iframe-url="%v" data-url="%v/edit" href="#" class="qor-actionbar-button">Edit</a>`, js, frameURL, url))
		}
		return template.HTML("")
	}

	return funcMap
}
