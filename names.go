package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 名字系统(Name System)

import (
	"fmt"
)

// 数据充值周期
type retsetCycle uint

// 所有类型数据Set时都会触发事件
const (
	// 临时数据
	RSC_TEMP retsetCycle = iota
	// 纯事件数据，仅发送时有效，且无法Get
	RSC_EVENT
	// 永久
	RSC_PERMANENT
	RSC_MINUTE
	RSC_HOUR
	RSC_DAY
	RSC_WEEK
	RSC_MONTH
	RSC_YEAR
)

const maxType = 0xFF
const typeIdBit = 8
const maxNameOfType = 0x00FFFFFF
const RawIDMark = 0x00FFFFFF
const OrderIdNameType = "OrderIDName"

type setFunc = func(store *Storehouse, id uint32, operSymbol OperSymbol, value float64) float64
type getFunc = func(store *Storehouse, id uint32) float64

type nameCfg struct {
	id      uint32
	name    string
	rsc     retsetCycle
	init    float64
	max     float64
	min     float64
	setFunc setFunc
	getFunc getFunc
}

type typeCfg struct {
	typeId      uint32
	flagOfRawID map[uint32]bool
	nameIdCount uint32
	setFunc     setFunc
	getFunc     getFunc
}

type names struct {
	orderIdCount   uint32
	typeIdCount    uint32
	typeCfgOfName  map[string]*typeCfg
	nameCfgOfName  map[string]*nameCfg
	nameCfgOfId    map[uint32]*nameCfg
	nameCfgsOfType map[string][]*nameCfg
}

type INames interface {
	GetIdByName(name string) uint32
	GetNameById(id uint32) string
	GetCfgById(id uint32) *nameCfg
}

func newNames() *names {
	ret := &names{}
	ret.init()
	return ret
}

func (this *names) init() {
	this.orderIdCount = uint32(0)
	this.typeIdCount = uint32(0)
	this.typeCfgOfName = make(map[string]*typeCfg)
	this.nameCfgOfName = make(map[string]*nameCfg)
	this.nameCfgOfId = make(map[uint32]*nameCfg)
	this.nameCfgsOfType = make(map[string][]*nameCfg)
}

func (this *names) GetCfgById(id uint32) *nameCfg {
	return this.nameCfgOfId[id]
}

func (this *names) registerType(typ string) *typeCfg {
	if typ == "" {
		typ = "empty+nil"
	}

	tCfg := this.typeCfgOfName[typ]
	if tCfg == nil {
		if this.typeIdCount >= maxType {
			panic("[dmp]registerType => 名字类型超过规定限数：" + typ)
		}
		this.typeIdCount++
		tCfg = &typeCfg{
			typeId:      this.typeIdCount << (32 - typeIdBit),
			nameIdCount: 0,
			flagOfRawID: make(map[uint32]bool),
		}
		this.typeCfgOfName[typ] = tCfg
	}

	return tCfg
}

func (this *names) RegisterType(typ string) uint32 {
	return this.registerType(typ).typeId
}

func (this *names) allocNameId(typ, name string, rawId uint32) (ret uint32, typeCfg *typeCfg) {
	if rawId >= maxNameOfType {
		panic("[dmp]allocNameId => 该类型原始ID大于规定数值：" + typ)
	}

	typeCfg = this.registerType(typ)
	if rawId != 0 {
		if typeCfg.flagOfRawID[rawId] {
			panic("[dmp]allocNameId => 该类型原始ID重复或者已经被注册：" + name)
		}
		if typeCfg.nameIdCount < rawId {
			typeCfg.nameIdCount = rawId
		}
		typeCfg.flagOfRawID[rawId] = true
		ret = typeCfg.typeId + rawId
	} else {
		if typeCfg.nameIdCount >= maxNameOfType {
			panic("[dmp]allocNameId => 该类型名字数量超过规定限数：" + typ)
		}
		typeCfg.nameIdCount++
		typeCfg.flagOfRawID[typeCfg.nameIdCount] = true
		ret = typeCfg.typeId + typeCfg.nameIdCount
	}
	return
}

func (this *names) RegisterNameOfOrderId(name string, rawId uint32, init, min, max float64) uint32 {
	if this.orderIdCount >= maxNameOfType {
		panic(fmt.Sprintf("[dmp]RegisterNameOfOrderId => 有序ID的名字数量已经达到设定的个数:%d", maxNameOfType))
	}

	id := this.GetIdByName(name)
	if id != 0 {
		panic(fmt.Sprintf("[dmp]RegisterNameOfOrderId => 名字重复：%s", name))
	}

	if rawId == 0 {
		this.orderIdCount++
		rawId = this.orderIdCount
	} else {
		if rawId > this.orderIdCount {
			this.orderIdCount = rawId
		}
	}

	typeCfg := this.registerType(OrderIdNameType)
	this.addName(OrderIdNameType, rawId, name, RSC_TEMP, init, min, max, typeCfg)

	return rawId
}

func (this *names) RegisterName(typ string, name string, rawId uint32) uint32 {
	return this.RegisterNameByInfo(typ, name, rawId, RSC_TEMP, 0, 0, 0)
}

const StrTypeName = "stringtypename"

func (this *names) RegisterNameByCycle(typ string, name string, rawId uint32, rsc retsetCycle) uint32 {
	return this.RegisterNameByInfo(typ, name, rawId, rsc, 0, 0, 0)
}

func (this *names) addName(typ string, id uint32, name string, rsc retsetCycle, init, min, max float64, typeCfg *typeCfg) {
	cfg := &nameCfg{
		id:      id,
		name:    name,
		rsc:     rsc,
		init:    init,
		min:     min,
		max:     max,
		getFunc: typeCfg.getFunc,
		setFunc: typeCfg.setFunc,
	}

	this.nameCfgOfName[name] = cfg
	this.nameCfgOfId[id] = cfg

	if typ != "" {
		cfgs := this.nameCfgsOfType[typ]
		cfgs = append(cfgs, cfg)
		this.nameCfgsOfType[typ] = cfgs
	}
}

func (this *names) RegisterNameByInfo(typ string, name string, rawId uint32, rsc retsetCycle, init, min, max float64) uint32 {
	if this.nameCfgOfName[name] != nil {
		panic("[dmp]RegisterNameByInfo => 名字重复: " + name)
	}

	id, typeCfg := this.allocNameId(typ, name, rawId)
	if this.nameCfgOfId[id] != nil {
		panic("[dmp]RegisterNameByInfo => 名字重复: " + name)
	}

	this.addName(typ, id, name, rsc, init, min, max, typeCfg)

	return id
}

func (this *names) GetIdByName(name string) uint32 {
	cfg := this.nameCfgOfName[name]
	if cfg == nil {
		return 0
	}
	return cfg.id
}

func (this *names) GetNameById(id uint32) string {
	cfg := this.nameCfgOfId[id]
	if cfg == nil {
		return ""
	}
	return cfg.name
}

func (this *names) RegisterSetFuncByType(typeName string, value setFunc) {
	tcfg := this.typeCfgOfName[typeName]
	if tcfg == nil {
		panic(fmt.Sprintf("[dmp]names.RegisterSetFuncByType => 类型名称“%s”尚未注册", typeName))
	}
	tcfg.setFunc = value

	cfgs := this.nameCfgsOfType[typeName]
	if cfgs != nil {
		for _, cfg := range cfgs {
			this.RegisterSetFuncById(cfg.id, value)
		}
	}
}

func (this *names) RegisterSetFuncByName(name string, value setFunc) {
	id := this.GetIdByName(name)
	if id == 0 {
		panic(fmt.Sprintf("[dmp]names.RegisterSetFuncByName => 数据名称“%s”尚未注册", name))
	}

	this.RegisterSetFuncById(id, value)
}

func (this *names) RegisterSetFuncById(id uint32, value setFunc) {
	if id <= maxNameOfType {
		panic(fmt.Sprintf("[dmp]names.RegisterSetFuncById => 有序数据“%s”不能设置Set函数", this.GetNameById(id)))
		return
	}

	cfg := this.nameCfgOfId[id]
	cfg.setFunc = value
}

func (this *names) RegisterGetFuncByType(typeName string, value getFunc) {
	tcfg := this.typeCfgOfName[typeName]
	if tcfg == nil {
		panic(fmt.Sprintf("[dmp]names.RegisterGetFuncByType => 类型名称“%s”尚未注册", typeName))
	}
	tcfg.getFunc = value

	cfgs := this.nameCfgsOfType[typeName]
	if cfgs != nil {
		for _, cfg := range cfgs {
			this.RegisterGetFuncById(cfg.id, value)
		}
	}
}

func (this *names) RegisterGetFuncByName(name string, value getFunc) {
	id := this.GetIdByName(name)
	if id == 0 {
		panic(fmt.Sprintf("[dmp]names.SetGetFunc => 数据名称“%s”尚未注册", name))
	}

	this.RegisterGetFuncById(id, value)
}

func (this *names) RegisterGetFuncById(id uint32, value getFunc) {
	if id <= maxNameOfType {
		panic(fmt.Sprintf("[dmp]names.SetGetFunc => 有序数据“%s”不能设置Get函数", this.GetNameById(id)))
		return
	}

	cfg := this.nameCfgOfId[id]
	cfg.getFunc = value
}

func strGetFunc(store *Storehouse, id uint32) float64 {
	// 字符串无数据值，默认用其ID来做比较
	return float64(id)
}

var Names = newNames()

func init() {
	Names.registerType(StrTypeName)
	Names.RegisterGetFuncByType(StrTypeName, strGetFunc)
}
