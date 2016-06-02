# ActionBar

ActionBar is dependant on [QOR Admin](https://github.com/qor/admin). It provides an action bar on the top of frontend page. The bar contains:

* Switcher of `Preview` and `Edit` mode
* Login/Logout links
* Additional links in a menu

[![GoDoc](https://godoc.org/github.com/qor/action_bar?status.svg)](https://godoc.org/github.com/qor/action_bar)

## Usage

```go
import "github.com/qor/admin"
import "github.com/qor/action_bar"

func main() {
  Admin := admin.New(&qor.Config{DB: db})

  // Register Global ActionBar object
  // Auth is admin.Auth interface, you need to define a struct and implements interface's functions
  ActionBar = action_bar.New(Admin, Auth)
  ActionBar.RegisterAction(&action_bar.Action{Name: "Admin Dashboard", Link: "/admin"})

  // Then use Render to render action bar in view
  ActionBar.Render(writer, request)
}

```

[Online Demo](http://demo.getqor.com/), you will see a bar at the top of homepage.

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).

