package stdlib

import (
	"bytes"
	"io"
	"net/http"

	"github.com/d5/tengo/v2"
)

var httpModule = map[string]tengo.Object{
	"do": &tengo.UserFunction{
		Name: "do",
		Value: tengo.CheckOptArgs(func(args ...tengo.Object) (tengo.Object, error) {
			// build req from method, url [, headers[, body]]
			method := args[0].(*tengo.String).Value
			url := args[1].(*tengo.String).Value
			var body io.Reader
			if len(args) > 3 {
				body = bytes.NewBuffer(args[3].(*tengo.Bytes).Value)
			}
			req, err := http.NewRequest(method, url, body)
			if err != nil {
				return wrapError(err), nil
			}
			// add headers
			if len(args) > 2 {
				for k, v := range args[2].(*tengo.Map).Value {
					s, err := tengo.ToString(0, v)
					if err != nil {
						return nil, err
					}
					req.Header.Add(k, s)
				}
			}

			// do req
			res, err := http.DefaultClient.Do(req)
			if res != nil && res.Body != nil {
				// ensure to always close body no matter what
				defer res.Body.Close()
			}
			if err != nil {
				return wrapError(err), nil
			}
			if res.ContentLength > int64(tengo.MaxBytesLen) {
				// don't allow going over byte limit
				return nil, tengo.ErrBytesLimit
			}

			// read full body, with byte limit on it
			bs, err := io.ReadAll(io.LimitReader(res.Body, int64(tengo.MaxBytesLen)))
			if err != nil {
				return wrapError(err), nil
			}
			resHeaders := &tengo.Map{Value: map[string]tengo.Object{}}
			for k := range res.Header {
				resHeaders.Value[k] = &tengo.String{Value: res.Header.Get(k)}
			}
			return &tengo.Map{
				Value: map[string]tengo.Object{
					"code":    &tengo.Int{Value: int64(res.StatusCode)},
					"status":  &tengo.String{Value: res.Status},
					"headers": resHeaders,
					"body": &tengo.Bytes{
						Value: bs,
					},
				},
			}, nil
		}, 2, 4, tengo.StringTN, tengo.StringTN, tengo.MapTN, tengo.BytesTN),
	},
}
