package jsonschema

import jsonscpt "github.com/seeadoog/json_script"

type Script struct {
	script jsonscpt.Exp
}

func (s Script) Validate(path *pathTree, value interface{}, errs *[]Error) {
	m,ok:=value.(map[string]interface{})
	if !ok{
		return
	}
	ctx:=jsonscpt.NewVmWithContext(m)
	err:=s.script.Exec(ctx)
	if err,ok:=jsonscpt.IsExitError(err);ok{
		*errs = append(*errs,Error{
			Path: path.String(),
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

