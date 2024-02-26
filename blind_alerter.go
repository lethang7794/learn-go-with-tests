package poker

import (
	"fmt"
	"io"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (fn BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	fn(duration, amount, to)
}

func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %v\n", amount) // TODO: remove hard-code os.Stdout
	})
}
