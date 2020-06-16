package schema

import (
	"fmt"
	"regexp"
)

type Pattern struct {
	regexp *regexp.Regexp
	Path string
}



func (p *Pattern) Validate(c *ValidateCtx,value interface{}) {
	str, ok := value.(string)
	if !ok {
		return
	}
	if !p.regexp.MatchString(str) {
		c.AddError(Error{
			Path: p.Path,
			Info: appendString(str, "value does not match pattern"),
		})
	}
}

func NewPattern(i interface{},path string,parent Validator) (Validator, error) {
	str, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("%s is not a string when assign regexp", str)
	}
	reg, err := regexp.Compile(str)
	if err != nil {
		return nil, fmt.Errorf("regexp compile error:%w", err)
	}
	return &Pattern{regexp: reg,Path: path}, nil
}


