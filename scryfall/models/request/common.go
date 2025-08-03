package request

import "net/url"

type RequestParams interface {
	ToURLValues() url.Values
	Validate() error
}
