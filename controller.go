package action_bar

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strings"

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
	context.Writer.Write([]byte(context.Render("action_bar/inline_edit")))
}

func (bar *ActionBar) RenderEditButtonWithResource(w http.ResponseWriter, r *http.Request, value interface{}, resources ...*admin.Resource) template.HTML {
	context := bar.admin.NewContext(w, r)
	editURL, _ := utils.JoinURL(context.URLFor(value, resources...), "edit")
	resourceName := "Resource"
	if res := bar.admin.GetResource(reflect.Indirect(reflect.ValueOf(value)).Type().String()); res != nil {
		resourceName = strings.ToUpper(res.Name)
	}
	title := string(bar.admin.T(context.Context, "qor_action_bar.action.edit_resource", "Edit {{$1}}", resourceName))
	return bar.RenderEditButton(w, r, title, editURL)
}

func (bar *ActionBar) RenderEditButton(w http.ResponseWriter, r *http.Request, title string, link string) template.HTML {
	if bar.EditMode(w, r) {
		var (
			prefix   = bar.admin.GetRouter().Prefix
			jsURL    = fmt.Sprintf("<script data-prefix=\"%v\" src=\"%v/assets/javascripts/action_bar_check.js?theme=action_bar\"></script>", prefix, prefix)
			frameURL = fmt.Sprintf("%v/action_bar/inline_edit", prefix)
		)

		return template.HTML(fmt.Sprintf(`%v<a target="blank" data-iframe-url="%v" data-url="%v" href="#" class="qor-actionbar-button">%v</a>`, jsURL, frameURL, link, title))
	}
	return template.HTML("")
}

// FuncMap will return helper to render inline edit button
func (bar *ActionBar) FuncMap(w http.ResponseWriter, r *http.Request) template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_edit_button"] = func(value interface{}, resources ...*admin.Resource) template.HTML {
		return bar.RenderEditButtonWithResource(w, r, value, resources...)
	}

	return funcMap
}
