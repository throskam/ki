package ki

import (
	"fmt"
	"maps"
	"net/url"
	"regexp"
	"strings"
)

// Location represents a resolved route.
type Location struct {
	prefix     string
	method     string
	pattern    string
	pathParams []string
	query      url.Values
}

// NewLocation returns a new Location.
func NewLocation(method, pattern string) Location {
	return Location{
		method:     method,
		pattern:    pattern,
		pathParams: []string{},
		query:      url.Values{},
	}
}

// Method returns the method of the route.
func (l Location) Method() string {
	return l.method
}

// Pattern returns the pattern of the route.
func (l Location) Pattern() string {
	return l.pattern
}

// URL returns the parameterized URL.
// It panics if the URL is invalid.
func (l Location) URL() *url.URL {
	path := l.pattern

	if l.prefix != "" {
		path = l.prefix + path
	}

	path = strings.ReplaceAll(path, "{$}", "")

	re := regexp.MustCompile(`\{[^}]+\}`)
	index := 0

	path = re.ReplaceAllStringFunc(path, func(match string) string {
		if index < len(l.pathParams) {
			replacement := l.pathParams[index]
			index++
			return replacement
		}
		return match
	})

	if len(l.query) > 0 {
		path = path + "?" + l.query.Encode()
	}

	u, err := url.Parse(path)
	if err != nil {
		panic(fmt.Sprintf("cannot parse URL (%v)", path))
	}

	return u
}

// WithPrefix returns a new Location with the prefix.
func (l Location) WithPrefix(prefix string) Location {
	l.prefix = prefix + l.prefix

	return l
}

// WithPathParams returns a new Location with the path parameters.
func (l Location) WithPathParams(params ...string) Location {
	l.pathParams = params

	return l
}

// WithQuery returns a new Location with the query.
func (l Location) WithQuery(query url.Values) Location {
	l.query = query

	return l
}

// WithQueryParam returns a new Location with the query parameter.
func (l Location) WithQueryParam(key, value string) Location {
	clone := url.Values{}

	maps.Copy(clone, l.query)

	l.query = clone

	l.query.Add(key, value)

	return l
}