package poker

import (
	"fmt"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type BlindAlerterFunc func(duration time.Duration, amount int)

func (fn BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	fn(duration, amount)
}

func StdOutBlindAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %v\n", amount)
	})
}
