package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 数据仓库(Data Storehouse)

import (
	"fmt"
	"unsafe"
)

type OperSymbol uint

const (
	OS_INVAILD OperSymbol = iota
	OS_INC
	OS_DEC
	OS_MUL
	OS_DIV
	OS_SET
)

type Data struct {
	Value float64
	cfg   *nameCfg
}

func (this *Data) ResetCycle() retsetCycle {
	return this.cfg.rsc
}

func (this *Data) Name() string {
	return this.cfg.name
}

func (this *Data) Id() uint32 {
	return this.cfg.id
}

func (this *Data) RawID() uint32 {
	return this.cfg.id & RawIDMark
}

type Storehouse struct {
	allowTriggerChgEvt bool
	names              INames
	datasOfOrderId     []float64
	cfgsOfOrderId      []*nameCfg
	datasOfHashId      map[uint32]*Data
	datasOfCycle       map[retsetCycle][]*Data
	Owner              unsafe.Pointer
	Lister             *lister
	workstatLog        map[*Workstat]*workstatLog
}

func NewStorehouse(owner unsafe.Pointer) *Storehouse {
	ret := &Storehouse{}
	ret.init(owner, Names)
	return ret
}

func (this *Storehouse) init(owner unsafe.Pointer, names *names) {
	this.allowTriggerChgEvt = true
	this.datasOfHashId = make(map[uint32]*Data)
	this.datasOfCycle = make(map[retsetCycle][]*Data)
	this.Owner = owner
	this.Lister = newLister()
	this.workstatLog = make(map[*Workstat]*workstatLog)
	if names != nil {
		this.setNames(names)
	}
}

func (this *Storehouse) setNames(names *names) {
	this.names = names
	count := names.orderIdCount
	if count > 0 {
		count++
		this.datasOfOrderId = make([]float64, count)
		this.cfgsOfOrderId = make([]*nameCfg, count)
		for i := uint32(1); i < count; i++ {
			cfg := names.GetCfgById(i)
			if cfg == nil {
				panic(fmt.Sprintf("[dmp]NewStorehouse => ID为“%d”的有序数据名缺失", i))
			}
			this.cfgsOfOrderId[i] = cfg
		}
	}
}

func (this *Storehouse) Set(id uint32, value float64) float64 {
	return this.Oper(id, OS_SET, value)
}

func (this *Storehouse) Oper(id uint32, operSymbol OperSymbol, value float64) (ret float64) {
	if id == 0 {
		panic("[dmp]Storehouse.Oper => 数据ID为0")
	}

	doOper := func(cfg *nameCfg, oldValue float64) (newValue float64) {
		switch operSymbol {
		case OS_INC:
			newValue = oldValue + value
		case OS_DEC:
			newValue = oldValue - value
		case OS_MUL:
			newValue = oldValue * value
		case OS_DIV:
			if value != 0 {
				newValue = oldValue / value
			}
		case OS_SET:
			newValue = value
		}

		if (cfg.max != 0) && (newValue > cfg.max) {
			newValue = cfg.max
		} else if newValue < cfg.min {
			newValue = cfg.min
		}
		return
	}

	defer func() {
		if this.allowTriggerChgEvt {
			defer func() {
				if err := recover(); err != nil {
					writeLog("[dmp]Storehouse.Oper => 数据监听回调异常：%v", err)
				}
			}()
			Lister.trigger(this, id, operSymbol, value)
			this.Lister.trigger(this, id, operSymbol, value)
		}
	}()

	if id < uint32(len(this.datasOfOrderId)) {
		old := &this.datasOfOrderId[id]
		*old = doOper(this.cfgsOfOrderId[id], *old)
		return *old
	}

	data := this.datasOfHashId[id]
	if data == nil {
		data = &Data{}
		data.cfg = this.names.GetCfgById(id)
		if data.cfg == nil {
			writeLog("[dmp]Storehouse.Oper => 无效数据ID: %d", id)
		}
		this.datasOfHashId[id] = data
		this.datasOfCycle[data.cfg.rsc] = append(this.datasOfCycle[data.cfg.rsc], data)
	}

	if data.cfg.rsc == RSC_EVENT {

	} else if data.cfg.setFunc != nil {
		defer func() {
			if err := recover(); err != nil {
				writeLog("[dmp]Storehouse.Set => 数据设置异常：" + fmt.Sprintf("%v", err))
			}
		}()
		data.cfg.setFunc(this, id, operSymbol, value)
	} else if data.cfg.getFunc != nil {
		// 如果设置了get函数但不设置set函数则不操作
	} else {
		data.Value = doOper(data.cfg, data.Value)
	}

	return data.Value
}

func (this *Storehouse) GetByName(name string) float64 {
	return this.Get(this.names.GetIdByName(name))
}

func (this *Storehouse) Get(id uint32) float64 {
	if id == 0 {
		return 0
	}
	if id < uint32(len(this.datasOfOrderId)) {
		return this.datasOfOrderId[id]
	}

	defer func() {
		if err := recover(); err != nil {
			writeLog("[dmp]Storehouse.Get => 数据获取异常：%v", err)
		}
	}()

	data := this.datasOfHashId[id]
	if data == nil {
		cfg := this.names.GetCfgById(id)
		if (cfg == nil) || (cfg.getFunc == nil) {
			return 0
		}
		return cfg.getFunc(this, id)
	}

	if data.cfg.getFunc != nil {
		return data.cfg.getFunc(this, id)
	}
	return data.Value
}

func (this *Storehouse) ResetById(id uint32) float64 {
	data := this.datasOfHashId[id]
	if data == nil {
		return 0
	}

	data.Value = data.cfg.init
	return data.Value
}

func (this *Storehouse) Reset() {
	for _, data := range this.datasOfHashId {
		data.Value = data.cfg.init
	}
}

func (this *Storehouse) ResetByCycle(cycle retsetCycle) {
	datas := this.datasOfCycle[cycle]
	if datas == nil {
		return
	}
	for _, data := range datas {
		data.Value = data.cfg.init
	}
}

func (this *Storehouse) Check(conds CondExp) bool {
	return conds.Check(this)
}

type execExp interface {
	Exec(store *Storehouse) float64
}

func (this *Storehouse) Exec(expr execExp) float64 {
	if expr == nil {
		return 0
	}
	return expr.Exec(this)
}
