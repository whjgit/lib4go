package lua

import (	
    m "crypto/md5"
	l "github.com/yuin/gopher-lua"
    "encoding/hex"
)

func syslibLoader(L *l.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]l.LGFunction{
	"md5": md5,
}

func md5(L *l.LState) int {
	input := L.ToString(1)
	md5Ctx := m.New()
	md5Ctx.Write([]byte(input))
	cipherStr := md5Ctx.Sum(nil)
	L.Push(l.LString(hex.EncodeToString(cipherStr)))
	return 1
}
