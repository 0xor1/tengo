package stdlib_test

import (
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

func TestHTTP(t *testing.T) {
	script := tengo.NewScript([]byte(`
http := import("http")

res := http.do("GET", "https://avatars.githubusercontent.com/u/1291934?s=48&v=4", {dnt: 1})

code := res.code
status := res.status
headers := res.headers
body := res.body
`))
	script.SetImports(stdlib.GetModuleMap("http"))
	executed, err := script.Run()
	if err != nil {
		t.Error(err)
	}
	check := func(name string, ok func(i interface{}) bool) {
		if !ok(executed.Get(name).Value()) {
			t.Errorf("unexpected %s value", name)
		}
	}
	check("code", func(i interface{}) bool { v, ok := i.(int64); return ok && v == 200 })
	check("status", func(i interface{}) bool { v, ok := i.(string); return ok && v == "200 OK" })
	check("headers", func(i interface{}) bool { v, ok := i.(map[string]interface{}); return ok && len(v) > 0 })
	check("body", func(i interface{}) bool { v, ok := i.([]byte); return ok && len(v) > 0 })
}
