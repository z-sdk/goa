package breaker

import (
	"fmt"
	"goa/lib/container"
	"goa/lib/proc"
	"goa/lib/stat"
)

type (
	internalPromise interface {
		Accept()
		Reject()
	}

	internalThrottle interface {
		allow() (internalPromise, error)
		doReq(req Request, fallback Fallback, acceptable Acceptable) error
	}

	promiseWithReason struct {
		promise internalPromise
		errWin  *container.ErrorWindow
	}

	loggedThrottle struct {
		name string
		internalThrottle
		errWin *container.ErrorWindow
	}
)

func newLoggedThrottle(name string, t internalThrottle) loggedThrottle {
	return loggedThrottle{
		name:             name,
		internalThrottle: t,
		errWin:           container.NewErrorWindow(),
	}
}

func (t loggedThrottle) allow() (Promise, error) {
	promise, err := t.internalThrottle.allow()
	return promiseWithReason{
		promise: promise,
		errWin:  t.errWin,
	}, t.logError(err)
}

func (t loggedThrottle) doReq(req Request, fallback Fallback, acceptable Acceptable) error {
	return t.logError(t.internalThrottle.doReq(req, fallback, func(err error) bool {
		accept := acceptable(err)
		if !accept {
			t.errWin.Add(err.Error())
		}
		return accept
	}))
}

func (t loggedThrottle) logError(err error) error {
	if err == ErrServiceUnavaliable {
		stat.Report(fmt.Sprintf(
			"proc(%s/%d), caller: %s, 断路器已打开，请求被丢弃\n最新错误：\n%s",
			proc.ProcessName(), proc.Pid(), t.name, t.errWin))
	}
	return err
}

func (p promiseWithReason) Accept() {
	p.promise.Accept()
}

func (p promiseWithReason) Reject(reason string) {
	p.errWin.Add(reason)
	p.promise.Reject()
}
