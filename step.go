package pipeline

import (
	"fmt"

	"github.com/fatih/color"
)

// Result is returned by a step to dispatch data to the next step or stage
type Result struct {
	Error error
	// dispatch any type
	Data interface{}
	// dispatch key value pairs
	KeyVal map[string]interface{}
}

// Request is the result dispatched in a previous step.
type Request struct {
	Data   interface{}
	KeyVal map[string]interface{}
}

// Step is the unit of work which can be concurrently or sequentially staged with other steps
type Step interface {
	out
	// Exec is invoked by the pipeline when it is run
	Exec(*Request) *Result
	// Cancel is invoked by the pipeline when one of the concurrent steps set Result{Error:err}
	Cancel() error
}

type out interface {
	Status(line string)
	getCtx() *stepContextVal
	setCtx(ctx *stepContextVal)
}

type stepContextVal struct {
	pipelineKey string
	name        string
	index       int
	concurrent  bool
}

// StepContext type is embedded in types which need to statisfy the Step interface
type StepContext struct {
	ctx *stepContextVal
}

func (sc *StepContext) getCtx() *stepContextVal {
	return sc.ctx
}

func (sc *StepContext) setCtx(ctx *stepContextVal) {
	sc.ctx = ctx
}

// Status is used to log status from a step
func (sc *StepContext) Status(line string) {
	stepText := fmt.Sprintf("[step-%d]", sc.getCtx().index)
	blue := color.New(color.FgBlue).SprintFunc()
	line = blue(stepText) + "[" + sc.getCtx().name + "]: " + line
	send(sc.getCtx().pipelineKey, line)
}


/// func
type HandleStepExec func(request *Request) *Result

func (h HandleStepExec) Exec(request *Request) *Result {
	return h(request)
}

func (w HandleStepExec) Cancel() error {
	w.Status("cancel step")
	return nil
}

func (h HandleStepExec) Status(line string) {
	fmt.Println("into custom status", line)
}

func (h HandleStepExec) getCtx() *stepContextVal {
	return nil
}

func (h HandleStepExec) setCtx(*stepContextVal) {
}
