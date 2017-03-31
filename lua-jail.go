package main

import (
	"context"
	"fmt"
	"github.com/yuin/gopher-lua"
	"os"
	"time"
)

func jailPrint(L *lua.LState) int {
	top := L.GetTop()
	fmt.Printf("%d ", time.Now().Unix())
	for i := 1; i <= top; i++ {
		fmt.Print(L.ToStringMeta(L.Get(i)).String())
		if i != top {
			fmt.Print("\t")
		}
	}
	fmt.Println("")
	return 0
}

var jailFunctions = map[string]lua.LGFunction{
	"print": jailPrint,
}

func openBase(L *lua.LState) int {
	global := L.Get(lua.GlobalsIndex).(*lua.LTable)
	L.SetGlobal("_G", global)
	basemod := L.RegisterModule("_G", jailFunctions)
	L.Push(basemod)
	return 1
}

func newState(jail bool) *lua.LState {
	if jail == false {
		return lua.NewState()
	}

	L := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})

	if err := L.CallByParam(lua.P{
		Fn:      L.NewFunction(openBase),
		NRet:    0,
		Protect: true,
	}, lua.LString(lua.BaseLibName)); err != nil {
		panic(err)
	}

	return L
}

func main() {
	if len(os.Args) < 2 {
		println("lua-jail <script> <off>")
		return
	}

	jail := true
	if len(os.Args) > 2 {
		jail = false
	}

	L := newState(jail)
	if jail {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		L.SetContext(ctx)
	}
	err := L.DoFile(os.Args[1])
	if err != nil {
		panic(err)
	}
}
