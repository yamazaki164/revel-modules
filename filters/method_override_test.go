package filters

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/revel/revel"
)

func dummyFilter1(c *revel.Controller, fc []revel.Filter) {
	fc[0](c, fc[1:])
}

func dummyFilter2(c *revel.Controller, fc []revel.Filter) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	fc[0](c, fc[1:])
}

func TestMethodOverrideFilterIsOverride(t *testing.T) {
	req := revel.NewRequest(httptest.NewRequest("POST", "http://localhost/", &bytes.Buffer{}))
	req.Form = url.Values{}
	req.Form.Add("_method", "DELETE")

	res := revel.NewResponse(httptest.NewRecorder())
	ctlr := revel.NewController(req, res)
	fc := []revel.Filter{
		dummyFilter1,
		dummyFilter2,
	}

	if ctlr.Request.Method != "POST" {
		t.Error("method error")
	}

	MethodOverrideFilter(ctlr, fc)

	if ctlr.Request.Method == "POST" {
		t.Error("method override error")
	}
}

func TestMethodOverrideFilterIsNotOverride(t *testing.T) {
	req := revel.NewRequest(httptest.NewRequest("PUT", "http://localhost/", &bytes.Buffer{}))
	res := revel.NewResponse(httptest.NewRecorder())
	ctlr := revel.NewController(req, res)
	fc := []revel.Filter{
		dummyFilter1,
		dummyFilter2,
	}

	if ctlr.Request.Method != "PUT" {
		t.Error("method error")
	}

	MethodOverrideFilter(ctlr, fc)

	if ctlr.Request.Method != "PUT" {
		t.Error("method error")
	}

	MethodOverrideFilter(ctlr, []revel.Filter{})
	if ctlr.Request.Method != "PUT" {
		t.Error("method error")
	}
}
