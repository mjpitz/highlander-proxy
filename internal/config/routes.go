package config

import (
	"strings"

	"github.com/spf13/pflag"
)

type RouteSlice struct {
	value []*Route
}

func (r *RouteSlice) Routes() []*Route {
	return r.value
}

func (r *RouteSlice) String() string {
	if r == nil {
		return ""
	}

	strs := make([]string, len(r.value))
	for i, val := range r.value {
		strs[i] = val.String()
	}
	return strings.Join(strs, ",")
}

func (r *RouteSlice) Set(value string) error {
	parts := strings.Split(value, ",")
	routes := make([]*Route, len(parts))

	for i, part := range parts {
		routes[i] = &Route{}

		if err := routes[i].Set(part); err != nil {
			return err
		}
	}

	r.value = append(r.value, routes...)
	return nil
}

func (r *RouteSlice) Type() string {
	return "RouteSlice"
}

var _ pflag.Value = &RouteSlice{}
