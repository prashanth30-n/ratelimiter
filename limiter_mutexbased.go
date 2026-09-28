package ratelimiter

import(
	"time"
	"sync"
)

type mutexLimiter struct{
	sync.Mutex
	last time.Time
	sleepFor time.Duration
	perRequest time.Duration
	maxSlack time.Duration
	clock Clock
}

func newMutexBased(rate int,opts ...Option) *mutexLimiter{
	config:=buildConfig(opts)
	perRequest:=config.per/time.Duration(rate)
	l:=&mutexLimiter{
		perRequest: perRequest,
		maxSlack: -1*time.Duration(config.slack)*perRequest,
		clock: config.clock,
	}
	return l
}


func(t*mutexLimiter)Take() time.Time{
	t.Lock()
	defer t.Unlock()
	now:=t.clock.Now()
	if(t.last.IsZero()){
		t.last=now
		return t.last
	}
	t.sleepFor+=t.perRequest-now.Sub(t.last)


	if t.sleepFor<t.maxSlack{
		t.sleepFor=t.maxSlack
	}
	if t.sleepFor>0{
		t.clock.Sleep(t.sleepFor)
		t.last=now.Add(t.sleepFor)
		t.sleepFor=0
	}else{
		t.last=now
	}
	return t.last
}