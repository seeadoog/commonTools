package wait

import (
	"context"
	"time"
)

func Until(ctx context.Context,f func()error,duration time.Duration){
	var timer *time.Timer
	for{
		select {
		case <-ctx.Done():
		default:

		}

		if err:=f();err != nil{
			if timer == nil{
				timer = time.NewTimer(duration)
			}
			select {
			case <-timer.C:
			case <-ctx.Done():
			}
		}

	}
}
