package ratelimiter

import(
	"time"
	"github.com/benbjohnson/clock"

)

type Limiter interface{
	Take() time.Time //what this does it returns the time at which the request was allowed to proceed
}

type Clock interface{
	Now() time.Time
	Sleep(time.Duration)
}
type config struct{
	//this config package is initiated by caller through options rather than using a constructor
	clock Clock
	slack int // slack controls how much unused time can accumulate
	per time.Duration //this is time window over which the rate is measured 
}
func New(rate int,opts ...Option) Limiter{
	//rate int is supplied by caller it is no of requests per configured time window
	return newAtomicInt64Based(rate,opts...)
	//...paramter meaning variadic parameter it allows  0 or more options
	//internally opts are like slice of options
	//New(100)
	//New(100,withoutslack)
	//New(100,withslack(20))
	//New(100,Per(time.Minute))

}
func buildConfig(opts[]Option) config{
	c:=config{
		clock:clock.New(),
		slack: 10,
		per:time.Second,
	}
	for _,opt:=range opts{
		opt.apply(&c)
	}
	return c
}

type Option interface{
	apply(*config)
}
type clockOption struct{
	clock Clock
}

func(o clockOption) apply(c *config){
	c.clock=o.clock
}

func WithClock(clock Clock) Option{
	return clockOption{
		clock: clock,
	}
}

type slackOption int

func(o slackOption) apply(c *config){
	c.slack=int(o)
}
var WithoutSlack Option=slackOption(0)

func WithSlack(slack int) Option{
	return slackOption(slack)
}

type perOption time.Duration

func(p perOption) apply(c *config){
	c.per=time.Duration(p)
}
func Per(per time.Duration) Option{
	return perOption(per)
}

type unlimited struct{}

func NewUnlimted() Limiter{
	return unlimited{}
}

func(unlimited) Take() time.Time{
	return time.Now()
}
