package jsonschema

import (
	"fmt"
	"regexp"
)

type Pattern struct {
	regexp *regexp.Regexp
	parent Validator
}

func (p *Pattern) Path() string {
	return p.parent.Path()
}

func (p *Pattern) Validate(c *ValidateCtx,value interface{}) {
	str, ok := value.(string)
	if !ok {
		return
	}
	if !p.regexp.MatchString(str) {
		c.AddError(&Error{
			Path: c.Path(),
			Info: appendString(str, "value does not match pattern"),
		})
	}
}

func NewPattern(i interface{},parent Validator) (Validator, error) {
	str, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("%s is not a string when assign regexp", str)
	}
	reg, err := regexp.Compile(str)
	if err != nil {
		return nil, fmt.Errorf("regexp compile error:%w", err)
	}
	return &Pattern{regexp: reg}, nil
}


