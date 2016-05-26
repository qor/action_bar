package action_bar

import (
	"github.com/qor/admin"
	"net/http"
)

func SwitchMode(context *admin.Context) {
	cookie := http.Cookie{Name: "qor-action-bar", Value: context.Request.URL.Query().Get("checked"), Path: "/", HttpOnly: true}
	http.SetCookie(context.Writer, &cookie)
	http.Redirect(context.Writer, context.Request, "/", http.StatusFound)
}
