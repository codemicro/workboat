// Generated by ego.
// DO NOT EDIT

//line loginPage.ego:1

package views

import "fmt"
import "html"
import "io"
import "context"

import "github.com/codemicro/workboat/workboat/paths"

func LoginPage(ctx context.Context, w io.Writer) {

//line loginPage.ego:8
	_, _ = io.WriteString(w, "\n    ")
//line loginPage.ego:8
	{
		var EGO Page
		EGO.Title = "Login"
		EGO.Yield = func() {
//line loginPage.ego:9
			_, _ = io.WriteString(w, "\n        <h1>Login</h1>\n        <a href=\"")
//line loginPage.ego:10
			_, _ = io.WriteString(w, html.EscapeString(fmt.Sprint(paths.Make(paths.AuthOauthOutbound))))
//line loginPage.ego:10
			_, _ = io.WriteString(w, "\">Click here to login via Gitea</a>\n    ")
		}
		EGO.Render(ctx, w)
	}
//line loginPage.ego:12
	_, _ = io.WriteString(w, "\n")
//line loginPage.ego:12
}

var _ fmt.Stringer
var _ io.Reader
var _ context.Context
var _ = html.EscapeString