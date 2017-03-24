package main

import (
	"bufio"
	"errors"
	"github.com/yuin/gopher-lua"
	"os"
)

type Filter struct {
	state                    *lua.LState
	exceptionHandlerFunction *lua.LFunction
}

func NewFilter() *Filter {
	state := lua.NewState()

	filter := &Filter{
		state: state,
	}
	filter.exceptionHandlerFunction = state.NewFunction(
		filter.exceptionHandler)
	return filter
}

func (f *Filter) exceptionHandler(L *lua.LState) int {
	panic("exception in lua code")
	return 0
}

func (f *Filter) LoadScript(filename string) error {
	return f.state.DoFile(filename)
}

func (f *Filter) ValidateScript() error {
	fn := f.state.GetGlobal("filter")
	if fn.Type() != lua.LTFunction {
		return errors.New("Function 'filter' not found")
	}
	return nil
}

func (f *Filter) ValidateEvent(event string) (bool, error) {
	fn := f.state.GetGlobal("filter")

	f.state.Push(fn.(*lua.LFunction))
	f.state.Push(lua.LString(event))

	// one argument and one return value
	err := f.state.PCall(1, 1, f.exceptionHandlerFunction)
	if err != nil {
		return false, err
	}

	top := f.state.GetTop()
	returnValue := f.state.Get(top)
	if returnValue.Type() != lua.LTBool {
		return false, errors.New("Invalid return value")
	}

	return lua.LVAsBool(returnValue), err
}

func main() {
	if len(os.Args) != 2 {
		println("provide filter script")
		return
	}

	filter := NewFilter()
	err := filter.LoadScript(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	err = filter.ValidateScript()
	if err != nil {
		panic(err.Error())
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		event := scanner.Text()
		isValid, err := filter.ValidateEvent(event)
		if err != nil {
			panic(err.Error())
		}

		if isValid {
			println(event)
		}
	}
}
