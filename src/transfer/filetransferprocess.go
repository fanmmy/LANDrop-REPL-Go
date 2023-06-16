package transfer

import (
	"awesomeProject1/src/utils"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"io"
)

type Process struct {
	DoingFile string
	AllNum    int
	DoneNum   int
	Receive   bool
}

func (p *Process) doingNo() int {
	return utils.Min(p.DoneNum+1, p.AllNum)
}

func (p *Process) ReceiveTips(completed bool) string {
	if completed {
		return fmt.Sprintf("Received ( %d of %d files) , done ! \n", p.AllNum, p.AllNum)
	} else {

		return fmt.Sprintf("Receiving %s (its %d of %d files) \n", p.DoingFile, p.doingNo(), p.AllNum)
	}
}
func (p *Process) SendTips(completed bool) string {
	if completed {
		return fmt.Sprintf("Sent ( %d of %d files) , done ! \n", p.AllNum, p.AllNum)
	} else {

		return fmt.Sprintf("Sending %s (its %d of %d files) \n", p.DoingFile, p.doingNo(), p.AllNum)
	}
}

func processDoingTipsFunc(process *Process) func(w io.Writer, s decor.Statistics) (err error) {
	return func(w io.Writer, s decor.Statistics) (err error) {
		if process.Receive {
			_, err = fmt.Fprintf(w, "%s", process.ReceiveTips(s.Completed))
		} else {
			_, err = fmt.Fprintf(w, "%s", process.SendTips(s.Completed))
		}
		return err
	}
}

func NewProcessBar(total int64, process *Process) (*mpb.Progress, *mpb.Bar) {
	p := mpb.New()

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
		mpb.BarExtender(mpb.BarFillerFunc(processDoingTipsFunc(process)), true),
		mpb.BarOptional(mpb.BarRemoveOnComplete(), true),
	)
	return p, bar
}
