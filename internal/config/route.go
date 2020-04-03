package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/pflag"
)

type Route struct {
	from *url.URL
	to   *url.URL
}

func (r *Route) To() *url.URL {
	return r.to
}

func (r *Route) From() *url.URL {
	return r.from
}

func (r *Route) String() string {
	if r == nil {
		return ""
	}

	return fmt.Sprintf("%s|%s", r.from, r.to)
}

func (r *Route) Set(val string) error {
	parts := strings.SplitN(val, "|", 2)

	if len(parts) != 2 {
		return fmt.Errorf("route: missing pipe")
	}

	fromString := parts[0]
	toString := parts[1]

	from, err := url.Parse(fromString)
	if err != nil {
		return fmt.Errorf("route: bind part not a valid URI")
	}
	to, err := url.Parse(toString)
	if err != nil {
		return fmt.Errorf("route: target part not a valid URI")
	}

	r.from = from
	r.to = to

	return nil
}

func (r *Route) Type() string {
	return "Route"
}

var _ pflag.Value = &Route{}
