package ratelimiter

import(
	"sync/atomic"
	"time"
)

type atomicInt64Limiter struct{
	prepadding [64]byte
	state int64 //the timestamp at which next request will be processed
	postpadding [56]byte
	perRequest time.Duration
	maxSlack time.Duration
	clock Clock
}

func newAtomicInt64Based(rate int,opt ...Option) *atomicInt64Limiter{
	config:=buildConfig(opt)
	perRequest:=config.per/time.Duration(rate)
	l:=&atomicInt64Limiter{
		perRequest: perRequest,
		maxSlack: time.Duration(config.slack)*perRequest,
		clock: config.clock,
	}
	atomic.StoreInt64(&l.state,0)
	return l
}


func(t*atomicInt64Limiter) Take() time.Time{
	var(
		newTimeOfNextPermissionIssue int64
		now int64
	)
	for{
		now=t.clock.Now().UnixNano()
		timOfNextPermissionIssue:=atomic.LoadInt64(&t.state)

		switch{
		case timOfNextPermissionIssue==0 ||(t.maxSlack==0 && now-newTimeOfNextPermissionIssue>int64(t.perRequest)):
				newTimeOfNextPermissionIssue=now
			
		case t.maxSlack>0 && now-timeOfNextPermissionIssue>int64(t.maxSlack)+int64(t.perRequest):
			newTimeOfNextPermissionIssue=now-int64(t.maxSlack)
		default:
			newTimeOfNextPermissionIssue=timOfNextPermissionIssue+int64(t.perRequest)
		}
		if atomic.CompareAndSwapInt64(&t.state, timeOfNextPermissionIssue, newTimeOfNextPermissionIssue) {
			break
		}

	}
	
	sleepDuration := time.Duration(newTimeOfNextPermissionIssue - now)
	if sleepDuration > 0 {
		t.clock.Sleep(sleepDuration)
		return time.Unix(0, newTimeOfNextPermissionIssue)
	}
	// return now if we don't sleep as atomicLimiter does
	return time.Unix(0, now)
}
