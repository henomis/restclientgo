package restclientgo

import (
	"fmt"
	"net/url"
)

type URLValues url.Values

// NewURLValues creates a new URLValues.
func NewURLValues() *URLValues {
	return &URLValues{}
}

// Add adds a string value to the URLValues.
func (p *URLValues) Add(key string, value *string) {
	if value != nil && *value != "" {
		(*url.Values)(p).Add(key, *value)
	}
}

// AddInt adds an int value to the URLValues.
func (p *URLValues) AddInt(key string, value *int) {
	if value != nil {
		(*url.Values)(p).Add(key, fmt.Sprintf("%d", *value))
	}
}

// AddBool adds a bool value to the URLValues.
func (p *URLValues) AddBool(key string, value *bool) {
	if value != nil {
		(*url.Values)(p).Add(key, fmt.Sprintf("%t", *value))
	}
}

// AddBoolAsInt adds a bool value to the URLValues as an int (0=false, 1=true).
func (p *URLValues) AddBoolAsInt(key string, value *bool) {
	if value != nil {
		if *value {
			(*url.Values)(p).Add(key, "1")
		} else {
			(*url.Values)(p).Add(key, "0")
		}
	}
}

// AddFloat adds a float value to the URLValues.
func (p *URLValues) AddFloat(key string, value *float64) {
	if value != nil {
		(*url.Values)(p).Add(key, fmt.Sprintf("%f", *value))
	}
}

// Del deletes the URLValues associated with key.
func (p *URLValues) Del(key string) {
	(*url.Values)(p).Del(key)
}

// Encode encodes the URLValues into "URL encoded" form ("bar=baz&foo=quux").
func (p *URLValues) Encode() string {
	return (*url.Values)(p).Encode()
}
