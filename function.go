package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 可自定义函数系统（User Defined Function System）

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type FuncExec = func(store *Storehouse, params []FrmlExp) float64

type funcExp struct {
	name   string
	params []FrmlExp
	exec   FuncExec
}

func (this *funcExp) Float64(store *Storehouse) float64 {
	return this.exec(store, this.params)
}

func (this *funcExp) NameExp() string {
	ret := ""
	for _, param := range this.params {
		if ret != "" {
			ret = ret + ","
		}
		ret = ret + param.NameExp()
	}
	return this.name + "(" + ret + ")"
}

func (this *funcExp) ValueExp(store *Storehouse) string {
	ret := ""
	for _, param := range this.params {
		if ret != "" {
			ret = ret + ","
		}
		ret = ret + param.ValueExp(store)
	}
	return this.name + "(" + ret + ")"
}

func (this *funcExp) EachId(fn func(id uint32)) {
	for _, param := range this.params {
		param.EachId(fn)
	}
}

type FuncParser interface {
	doParse(names INames, params []string) FrmlExp
}

type funcParser struct {
	name   string
	pcount int
	exec   FuncExec
}

func (this *funcParser) doParse(names INames, params []string) FrmlExp {
	len := len(params)
	if len != this.pcount {
		panic(fmt.Sprintf("%s函数参数必须为%d个，当前为：%d", this.name, this.pcount, len))
	}

	ret := &funcExp{
		name:   this.name,
		exec:   this.exec,
		params: make([]FrmlExp, len),
	}
	for i, p := range params {
		ret.params[i] = parseFrmlExpByNames(p, names)
	}

	return ret
}

var funcParserOfName = make(map[string]FuncParser)

func RegisterFunc(name string, exec FuncExec, paramCount int) {
	if funcParserOfName[name] != nil {
		panic(fmt.Sprintf("函数“%s”已经被注册", name))
	}
	funcParserOfName[name] = &funcParser{
		name:   name,
		exec:   exec,
		pcount: paramCount,
	}
}

func getFuncParser(name string) FuncParser {
	return funcParserOfName[name]
}

type ifFuncExp struct {
	funcExp
	cond CondExp
}

func (this *ifFuncExp) EachId(fn func(id uint32)) {
	this.cond.EachId(fn)
	this.funcExp.EachId(fn)
}

func (this *ifFuncExp) NameExp() string {
	return "if(" + this.cond.NameExp() + "," + this.params[0].NameExp() + "," + this.params[1].NameExp() + ")"
}

func (this *ifFuncExp) ValueExp(store *Storehouse) string {
	return "If(" + this.cond.ValueExp(store) + "," + this.params[0].ValueExp(store) + "," + this.params[1].ValueExp(store) + ")"
}

func (this *ifFuncExp) Float64(store *Storehouse) float64 {
	if store.Check(this.cond) {
		return this.params[0].Float64(store)
	}
	return this.params[1].Float64(store)
}

type ifFuncParser struct {
	funcParser
}

func (this *ifFuncParser) doParse(names INames, params []string) FrmlExp {
	len := len(params)
	if len != 3 {
		panic(fmt.Sprintf("%s函数参数必须为%d个，当前为：%d", this.name, this.pcount, len))
	}

	ret := &ifFuncExp{}
	ret.name = this.name
	ret.cond = parseCondExpByNames(params[0], names)
	ret.params = make([]FrmlExp, len-1)
	for i, p := range params[1:] {
		ret.params[i] = parseFrmlExpByNames(p, names)
	}

	return ret
}

func init() {
	parser := &ifFuncParser{}
	parser.name = "If"
	funcParserOfName["If"] = parser
}

type random1FuncExp struct {
	funcExp
}

func (this *random1FuncExp) Float64(store *Storehouse) float64 {
	return float64(rand.Intn(int(this.params[0].Float64(store))))
}

type random2FuncExp struct {
	funcExp
}

func (this *random2FuncExp) Float64(store *Storehouse) float64 {
	return float64(int(this.params[0].Float64(store)) + rand.Intn(int(this.params[1].Float64(store)-this.params[0].Float64(store))))
}

type randomFuncParser struct {
	funcParser
}

func (this *randomFuncParser) doParse(names INames, params []string) FrmlExp {
	switch len(params) {
	case 1:
		{
			ret := &random1FuncExp{}
			ret.name = this.name
			ret.params = make([]FrmlExp, 1)
			ret.params[0] = parseFrmlExpByNames(params[0], names)
			return ret
		}
	case 2:
		{
			ret := &random2FuncExp{}
			ret.name = this.name
			ret.params = make([]FrmlExp, 2)
			ret.params[0] = parseFrmlExpByNames(params[0], names)
			ret.params[1] = parseFrmlExpByNames(params[1], names)
			return ret
		}
	default:
		{
			panic(fmt.Sprintf("不合法的random函数参数个数：%d", len(params)))
		}

	}
}

func init() {
	parser := &randomFuncParser{}
	parser.name = "Random"
	funcParserOfName["Random"] = parser
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		v1 := params[0].Float64(store)
		v2 := params[1].Float64(store)
		if v1 < v2 {
			return v1
		}
		return v2
	}
	RegisterFunc("Min", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		v1 := params[0].Float64(store)
		v2 := params[1].Float64(store)
		if v1 > v2 {
			return v1
		}
		return v2
	}
	RegisterFunc("Max", funcExec, 2)
}
