package app

import (
	"flag"
)

type FlagsRegistry struct {
	Flags map[string]interface{}
}

type Flags struct {
	Registry *FlagsRegistry
}

func (r *FlagsRegistry) Bool(name string, value bool, usage string) {
	r.Flags[name] = flag.Bool(name, value, usage)
}

func (r *FlagsRegistry) String(name string, value string, usage string) {
	r.Flags[name] = flag.String(name, value, usage)
}

func (f *Flags) Bool(name string) bool {
	v, has := f.Registry.Flags[name]
	if !has {
		return false
	}

	return *v.(*bool)
}

func (f *Flags) String(name string) string {
	v, has := f.Registry.Flags[name]
	if !has {
		return ""
	}

	return *v.(*string)
}
