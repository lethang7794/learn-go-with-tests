package poker

import (
	"fmt"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type BlindAlert struct {
}

func (b BlindAlert) ScheduleAlertAt(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %v\n", amount)
	})
}
