package print

import (
	"os"
)

type Printer struct {
	prefix string
}

func WithPrefix(prefix string) *Printer {
	return &Printer{
		prefix: prefix,
	}
}

func (p *Printer) Print(v ...interface{}) {
	print(prefixPrint("", p.prefix), v...)
}

func (p *Printer) Debug(v ...interface{}) {
	if !isDebug {
		return
	}
	print(prefixDebug(DebugPrePrefix, p.prefix), v...)
}

func (p *Printer) Info(v ...interface{}) {
	print(prefixInfo(InfoPrePrefix, p.prefix), v...)
}

func (p *Printer) Warn(v ...interface{}) {
	print(prefixWarn(WarnPrePrepix, p.prefix), v...)
}

func (p *Printer) Error(v ...interface{}) {
	print(prefixError(ErrorPrePrefix, p.prefix), v...)
}

func (p *Printer) Fatal(err error) {
	if err == nil {
		return
	}
	p.Error(err)
	os.Exit(1)
}
