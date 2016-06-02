package action_bar

import (
	"net/http"

	"github.com/qor/admin"
)

// SwitchMode is handle to store switch status in cookie
func SwitchMode(context *admin.Context) {
	cookie := http.Cookie{Name: "qor-action-bar", Value: context.Request.URL.Query().Get("checked"), Path: "/", HttpOnly: true}
	http.SetCookie(context.Writer, &cookie)

	referrer := context.Request.Referer()
	if referrer == "" {
		referrer = "/"
	}

	http.Redirect(context.Writer, context.Request, referrer, http.StatusFound)
}
