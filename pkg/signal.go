package pkg

import (
	"os"
	"os/signal"
)

type handleSlice []func()

var handleMap = make(map[string]handleSlice)

// registerHandler registers signal disposition
// Not a concurrent safe function
func registerHandler(
	sig os.Signal,
	handler ...func()) {

	if hslice, ok := handleMap[sig.String()]; ok {
		handleMap[sig.String()] = append(hslice, handler...)
		return
	}

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, sig)
	handleMap[sig.String()] = append(handleSlice{}, handler...)

	go func() {
		for {
			<-sigchan

			for _, h := range handleMap[sig.String()] {
				go h()
			}
		}
	}()
}
