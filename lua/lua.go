package lua

import (
	"lib4go/pool"
	l "github.com/yuin/gopher-lua"
    "errors"
)

type luaPoolObject struct {
	state *l.LState
}
type luaPoolFactory struct{
    script string
    count int
}
func (p *luaPoolObject) Close() {
	if p.state != nil {
		p.state.Close()
	}
}
func (f *luaPoolFactory) Create() pool.Object {
    f.count++
	o := &luaPoolObject{}
	o.state = l.NewState()   
    o.state.PreloadModule("sys",syslibLoader)
    er:= o.state.DoFile(f.script)
    if er!=nil{
       panic(er)
    }     
	return o
}
func (f *luaPoolFactory) registerFunc(name string, fun l.LGFunction, obj *luaPoolObject) {
	obj.state.SetGlobal(name, obj.state.NewFunction(fun))
}

//LuaPool  LUA对象池
type LuaPool struct {
	p *pool.ObjectPool
}

//NewLuaPool 构建LUA对象池
func NewLuaPool() *LuaPool {
	return &LuaPool{p: pool.New()}
}

//Load 预加载脚本
func (p *LuaPool) Load(script string, size int)(int) {
	return p.p.Register(script, &luaPoolFactory{script:script}, size)
}

//Call 执行脚本main函数
func (p *LuaPool) Call(script string, input ...string) ([]string, error) {
	o, er := p.p.Get(script)
	if er != nil {
		p.Load(script,10)
	}
	defer p.p.Recycle(script, o)
	L := o.(*luaPoolObject).state
	co := L.NewThread()         /* create a new thread */   
    main:= L.GetGlobal("main")
    if main == l.LNil{
        panic(errors.New("cant find main func"))
    }
	fn :=main.(*l.LFunction) /* get function from lua */
	var inputs [20]l.LValue
	for i, v := range input {
		inputs[i] =l.LString(v)
	}
	st, err, values := L.Resume(co, fn, inputs[0:len(input)]...)
	if st == l.ResumeError {
		return nil, err
	}
	var buffer [20]string
	for i, lv := range values {
		buffer[i] = lv.String()
	}
	return buffer[0:len(values)], nil
}
