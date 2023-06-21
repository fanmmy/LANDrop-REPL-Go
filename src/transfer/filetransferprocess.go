package transfer

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Process struct {
	DoingFile string
	AllNum    int
	DoneNum   int
	Receive   bool
}

func NewProcessBar(total int64, process *Process) (*mpb.Progress, *mpb.Bar) {
	p := mpb.New(mpb.PopCompletedMode(), mpb.WithOutput(l.Stdout()), mpb.WithAutoRefresh())

	bar := p.New(total,
		mpb.BarStyle().Tip(`-`, `\`, `|`, `/`),
		mpb.PrependDecorators(
			decor.Name("process:"),
		),
		mpb.AppendDecorators(
			decor.NewPercentage("%d , "),

			decor.Name("ETA: "),
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO), "done",
			),
		),
		//mpb.BarExtender(mpb.BarFillerFunc(efn), true),
		//mpb.BarExtender(mpb.BarFillerFunc(processDoingTipsFunc(process)), true),
		//mpb.BarOptional(mpb., true),
		mpb.BarOptional(mpb.BarRemoveOnComplete(), true),
		mpb.BarOptional(mpb.BarFillerClearOnComplete(), true),
	)

	return p, bar
}
