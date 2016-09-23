package action_bar

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/test/utils"

	"github.com/headzoo/surf"
)

var Server *httptest.Server
var CurrentUser *User
var actionBar *ActionBar

type User struct {
	Name string
	Role string
}

func (user *User) DisplayName() string {
	return user.Name
}

type AdminAuth struct {
}

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/auth/logout"
}

func (AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	if CurrentUser.Role != "admin" {
		return nil
	}
	return CurrentUser
}

// Init
func init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(actionBar.Render(w, req)))
	})
	db := utils.TestDB()
	Server = httptest.NewServer(mux)
	Admin := admin.New(&qor.Config{DB: db})
	Admin.MountTo("/admin", mux)

	actionBar = New(Admin, AdminAuth{})
}

func TestAuth(t *testing.T) {
	bow := surf.NewBrowser()
	CurrentUser = &User{Name: "Kaea", Role: ""}
	bow.Open(Server.URL + "/")
	if !strings.Contains(bow.Body(), "LOGIN") {
		t.Errorf(color.RedString("Should Get `LOGIN` link if current user is normal user"))
	}

	CurrentUser = &User{Name: "Admin", Role: "admin"}
	bow.Open(Server.URL + "/")
	if !strings.Contains(bow.Body(), "LOGOUT") {
		t.Errorf(color.RedString("Should Get `LOGOUT` link if current user is admin user"))
	}
}
