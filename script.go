package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"go.starlark.net/starlark"
)

func (r *RequestData) ToStarlark() *starlark.Dict {
	res := starlark.NewDict(10)
	res.SetKey(starlark.String("Date"), starlark.MakeInt64(r.Date))
	res.SetKey(starlark.String("ContentLength"), starlark.MakeInt64(r.ContentLength))
	res.SetKey(starlark.String("Headers"), headerToStarDict(r.Header))
	res.SetKey(starlark.String("Body"), starlark.String(r.Body))
	res.SetKey(starlark.String("Method"), starlark.String(r.Method))
	res.SetKey(starlark.String("Path"), starlark.String(r.Path))
	res.SetKey(starlark.String("Query"), starlark.String(r.Query))
	return res
}

func headerToStarDict(h http.Header) *starlark.Dict {
	res := starlark.NewDict(20)
	for k, vals := range h {
		headers := []starlark.Value{}
		for _, v := range vals {
			headers = append(headers, starlark.String(v))
		}
		res.SetKey(starlark.String(k), starlark.NewList(headers))
	}
	return res
}

func scriptResponse(bucket, script string, req *RequestData) (string, error) {
	out := new(bytes.Buffer)
	thread := &starlark.Thread{
		Name:  bucket,
		Print: func(_ *starlark.Thread, msg string) { fmt.Fprintln(out, msg) },
	}
	_, err := starlark.ExecFile(thread, "test.star", []byte(script), starlark.StringDict{
		"request": req.ToStarlark(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return out.String(), err
}
