package filters

import "github.com/revel/revel"

// for method override filter
// e.g.)
// <form method="post" action="your_action">
//   <input type="hidden" name="_method" value="DELETE">
//   <input type="submit" value="send">
// </form>
var MethodOverrideFilter = func(c *revel.Controller, fc []revel.Filter) {
	if c.Request.Method == "POST" && c.Request.FormValue("_method") != "" {
		c.Request.Method = c.Request.FormValue("_method")
	}

	if len(fc) > 1 {
		fc[0](c, fc[1:])
	}
}
