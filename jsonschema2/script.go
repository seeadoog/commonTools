package jsonschema

import jsonscpt "github.com/seeadoog/json_script"

type Script struct {
	script jsonscpt.Exp
}

func (s Script) Validate(c *ValidateCtx,value interface{}) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	ctx:=jsonscpt.NewVmWithContext(m)
	err:=s.script.Exec(ctx)
	if err,ok:=jsonscpt.IsExitError(err);ok{
		c.AddError(&Error{
			Path: c.Path(),
			Info: err.Message,
		})
	}
}

func NewScript(i interface{},parent Validator)(Validator,error){
	exp,err:=jsonscpt.CompileExpFromJsonObject(i)
	if err != nil{
		return nil, err
	}
	return &Script{script: exp},nil
}

