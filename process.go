package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 运算流程工艺控制系统(Operation Process Control System)

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ProcExp interface {
	NameExp() string
	ValueExp(store *Storehouse) string
	Steps() []OperExp
	LoadFrom(src *Storehouse)
	SaveTo(dest *Storehouse)
	MyStore() *Storehouse
}

var returnID = uint32(1)

type returnExp struct {
	names INames
	value FrmlExp
}

func (this *returnExp) DestId() uint32 {
	return 0
}

func (this *returnExp) NameExp() string {
	return "return(" + this.value.NameExp() + ")"
}

func (this *returnExp) ValueExp(store *Storehouse) string {
	return "return(" + strconv.FormatFloat(this.value.Float64(store), 'f', -1, 64) + ")"
}

func (this *returnExp) Exec(store *Storehouse) bool {
	store.Set(returnID, this.value.Float64(store))
	return true
}

type ifReturnExp struct {
	names INames
	cond  CondExp
	value FrmlExp
}

func (this *ifReturnExp) DestId() uint32 {
	return 0
}

func (this *ifReturnExp) NameExp() string {
	return "return(" + this.cond.NameExp() + "," + this.value.NameExp() + ")"
}

func (this *ifReturnExp) ValueExp(store *Storehouse) string {
	return "return(" + this.cond.ValueExp(store) + "," + strconv.FormatFloat(this.value.Float64(store), 'f', -1, 64) + ")"
}

func (this *ifReturnExp) Exec(store *Storehouse) bool {
	if store.Check(this.cond) {
		store.Set(returnID, this.value.Float64(store))
		return true
	}
	return false
}

type procExp struct {
	operSet
	store     *Storehouse
	rawIds    map[uint32]uint32
	rawValues []float64
}

func (this *procExp) MyStore() *Storehouse {
	return this.store
}

func (this *procExp) Steps() []OperExp {
	return this.operSet
}

func (this *procExp) LoadFrom(src *Storehouse) {
	for id, rawId := range this.rawIds {
		value := src.Get(rawId)
		this.store.datasOfOrderId[id] = value
		this.rawValues[id] = value
	}
}

func (this *procExp) SaveTo(dest *Storehouse) {
	for id, rawId := range this.rawIds {
		value := this.store.datasOfOrderId[id]
		if this.rawValues[id] != value {
			dest.Set(rawId, value)
		}
	}
}

func (this *procExp) NameExp() (ret string) {
	ret = ""
	for _, step := range this.operSet {
		if ret == "" {
			ret = step.NameExp()
		} else {
			ret = ret + string(stepSeparator) + "\r\n" + step.NameExp()
		}
	}
	return
}

func (this *procExp) ValueExp(src *Storehouse) (ret string) {
	this.LoadFrom(src)

	ret = ""
	for _, step := range this.operSet {
		if ret == "" {
			ret = step.ValueExp(this.store)
		} else {
			ret = ret + string(stepSeparator) + "\r\n" + step.ValueExp(this.store)
		}
	}
	return
}

func (this *procExp) recordRawId(names *procNames) {
	count := names.orderIdCount + 1
	this.rawIds = make(map[uint32]uint32)
	this.rawValues = make([]float64, count)
	this.store.datasOfOrderId = make([]float64, count)
	for id, rawId := range names.rawIds {
		this.rawIds[id] = rawId
	}
}

type procParser struct {
	operParser
}

type procNames struct {
	names
	rawIds map[uint32]uint32
}

func (this *procNames) GetIdByName(name string) uint32 {
	id := this.names.GetIdByName(name)
	if id != 0 {
		return id
	}

	ret := this.RegisterNameOfOrderId(name, 0, 0, 0, 0)
	rawId := Names.GetIdByName(name)
	if rawId > 0 {
		this.rawIds[ret] = rawId
	}
	return ret
}

func (this *procParser) doBuildExp(symbol rune, nameStart, nameEnd int) (ret OperExp) {
	if symbol == '(' {
		name := string(this.exp[nameStart:this.index])
		if name == "" {
			this.doError("此处不应该出现左括号")
		}
		if strings.ToLower(name) != "return" {
			this.doError("无效函数名称：" + name)
		}

		params := this.readFuncParams()
		len := len(params)
		switch len {
		case 0:
			{
				this.doError("return函数必须缺失参数")
			}
		case 1:
			{
				return &returnExp{value: parseFrmlExpByNames(params[0], this.names)}
			}
		case 2:
			{
				return &ifReturnExp{
					cond:  parseCondExpByNames(params[0], this.names),
					value: parseFrmlExpByNames(params[1], this.names),
				}
			}
		default:
			{
				this.doError("return函数的参数数量不能超过2个")
			}
		}
	}

	return this.operParser.doBuildExp(symbol, nameStart, nameEnd)
}

func parseProcExp(exp string) ProcExp {
	ret := procExp{}
	ret.store = &Storehouse{}
	ret.store.init(nil, nil)

	if exp == "" {
		return &ret
	}

	onGetOperExp := func(exp OperExp) {
		ret.operSet = append(ret.operSet, exp)
	}

	nms := &procNames{}
	nms.init()
	nms.rawIds = make(map[uint32]uint32)
	nms.RegisterNameOfOrderId("return", 1, 0, 0, 0)

	perser := &procParser{}
	perser.buildExp = perser.doBuildExp
	perser.init(exp, nms, "Process")
	perser.doParse(onGetOperExp)
	perser.checkEnd()

	if len(ret.operSet) == 0 {
		perser.doError("无效运算表达式")
	}

	ret.store.setNames(&nms.names)
	ret.recordRawId(nms)
	return &ret
}

func ParseProcExp(exp string) (ProcExp, error) {
	if err := recover(); err != nil {
		return nil, errors.New(fmt.Sprintf("%v", err))
	}
	ret := parseProcExp(exp)
	return ret, nil
}
