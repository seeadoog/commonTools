package resource

import "github.com/seeadoog/commonTools/ngcfg"

type ResLoader func(name string , cfg *ngcfg.Elem)error

/**
res{
	redis{
		dx {
			mode
		}
	}
}
 */