package hitrix

import (
	"flag"
)

type FlagsRegistry struct {
	flags map[string]interface{}
}

type Flags struct {
	registry *FlagsRegistry
}

func (r *FlagsRegistry) Bool(name string, value bool, usage string) {
	r.flags[name] = flag.Bool(name, value, usage)
}

func (r *FlagsRegistry) String(name string, value string, usage string) {
	r.flags[name] = flag.String(name, value, usage)
}

func (f *Flags) Bool(name string) bool {
	v, has := f.registry.flags[name]
	if !has {
		return false
	}
	return *v.(*bool)
}

func (f *Flags) String(name string) string {
	v, has := f.registry.flags[name]
	if !has {
		return ""
	}
	return *v.(*string)
}
