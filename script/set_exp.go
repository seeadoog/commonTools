package script

type SetExpression struct {
	key   string
	value interface{}
}

func (e *SetExpression) Exec(c *Context) *Error {
	c.Set(e.key, e.value)
	return nil
}

type IfExpression struct {
	ife   Value
	then  Expression
	elsee Expression
}

func (e *IfExpression) Exec(c *Context) *Error {
	ok := Bool(e.ife.Get(c))
	if ok && e.then != nil {
		return e.then.Exec(c)
	}
	if e.elsee != nil {
		return e.elsee.Exec(c)
	}
	return nil
}

type ForExpression struct {
	v  Value
	e  Expression
	kk string
	vk string
}

func (e *ForExpression) Exec(c *Context) *Error {
	c = c.Next()
	val := e.v.Get(c)
	switch val.(type) {
	case map[string]interface{}:
		for k, v := range val.(map[string]interface{}) {
			c.Set(e.kk, k)
			c.Set(e.vk, v)
			if err := e.e.Exec(c); err != nil {
				return err
			}
		}
	case []interface{}:
		for k, v := range val.([]interface{}) {
			c.Set(e.kk, k)
			c.Set(e.vk, v)
			if err := e.e.Exec(c); err != nil {
				return err
			}
		}
	case nil:
		return nil
	}
	return Errorf("%v: value is not ranged", val)
}

type Expressions []Expression

func (e Expressions) Exec(c *Context) *Error {
	for _, expression := range e {
		if err := expression.Exec(c); err != nil {
			return err
		}
	}
	return nil
}
