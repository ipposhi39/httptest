package errof

import (
	"fmt"
	"runtime"
)

// PanicToErr :
func PanicToErr(p interface{}) (err error) {
	var msgs []string
	for depth := 0; ; depth++ {
		pc, src, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		msg := fmt.Sprintf("Panic Trace --> depth: %d, pc: %s, src: %s, line: %d", depth, runtime.FuncForPC(pc).Name(), src, line)
		msgs = append(msgs, msg)
	}
	var ok bool
	err, ok = p.(error)
	if !ok {
		err = fmt.Errorf("panic: %v, trace: %v", p, msgs)
	}
	return err
}
