package common

import (
	"time"

	"github.com/briandowns/spinner"
)

func StartProgressSpinner(prefix string) (stop func()) {
	s := spinner.New(spinner.CharSets[9], time.Millisecond*200)
	s.Prefix = prefix + " "
	s.Start()

	return s.Stop
}
