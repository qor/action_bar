package action_bar

import (
	"fmt"
	"html/template"

	"github.com/qor/admin"
	"github.com/qor/qor/utils"
)

type ActionInterface interface {
	InlineAction() bool
	ToHTML(*admin.Context) template.HTML
}

type Action struct {
	EditModeOnly bool
	Inline       bool
	Name         string
	Link         string
}

func (action Action) InlineAction() bool {
	return action.Inline
}

func (action Action) ToHTML(context *admin.Context) template.HTML {
	if action.EditModeOnly && !isEditMode(context) {
		return template.HTML("")
	}
	name := context.Admin.T(context.Context, "qor_action_bar.action."+action.Name, action.Name)
	return template.HTML(fmt.Sprintf("<a href='%v'>%v</a>", action.Link, name))
}

type EditResourceAction struct {
	EditModeOnly bool
	Inline       bool
	Value        interface{}
	Resource     *admin.Resource
}

func (action EditResourceAction) InlineAction() bool {
	return action.Inline
}

func (action EditResourceAction) ToHTML(context *admin.Context) template.HTML {
	if action.EditModeOnly && !isEditMode(context) {
		return template.HTML("")
	}

	var (
		admin          = context.Admin
		resourceParams = utils.ModelType(action.Value).Name()
		resourceName   = resourceParams
		editURL, _     = utils.JoinURL(context.URLFor(action.Value, action.Resource), "edit")
	)

	if action.Resource != nil {
		resourceParams = action.Resource.ToParam()
		resourceName = string(admin.T(context.Context, fmt.Sprintf("%v.name", resourceParams), action.Resource.Name))
	}

	name := context.Admin.T(context.Context, "qor_action_bar.action.edit_"+resourceParams, fmt.Sprintf("Edit %v", resourceName))

	return template.HTML(fmt.Sprintf("<a href='%v'>%v</a>", editURL, name))
}

type HTMLAction struct {
	EditModeOnly bool
	HTML         template.HTML
}

func (action HTMLAction) InlineAction() bool {
	return true
}

func (action HTMLAction) ToHTML(context *admin.Context) template.HTML {
	if action.EditModeOnly && !isEditMode(context) {
		return template.HTML("")
	}

	return action.HTML
}
