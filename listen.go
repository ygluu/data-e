package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 数据变化监听系统(Data Status Monitoring System)

type ListenFuncById = func(store *Storehouse, id uint32, operSymbol OperSymbol, value float64)

type lister struct {
	funcsById map[uint32]map[*ListenFuncById]bool
}

var Lister = newLister()

func newLister() *lister {
	return &lister{
		funcsById: make(map[uint32]map[*ListenFuncById]bool),
	}
}

func (this *lister) trigger(store *Storehouse, id uint32, operSymbol OperSymbol, value float64) {
	funcsByIdOfId := this.funcsById[id]
	if funcsByIdOfId != nil {
		for fc, _ := range funcsByIdOfId {
			(*fc)(store, id, operSymbol, value)
		}
	}
}

func (this *lister) AddById(id uint32, listerFunc ListenFuncById) *ListenFuncById {
	if id == 0 {
		panic("[dmp]lister.AddById => ID不能为0")
	}

	funcsByIdOfId := this.funcsById[id]
	if funcsByIdOfId == nil {
		funcsByIdOfId = make(map[*ListenFuncById]bool)
		this.funcsById[id] = funcsByIdOfId
	}
	if funcsByIdOfId[&listerFunc] {
		panic("[dmp]lister.AddById => 同一个ID不能重复相同的监听器")
	}

	funcsByIdOfId[&listerFunc] = true
	return &listerFunc
}

func (this *lister) AddByName(name string, listerFunc ListenFuncById) *ListenFuncById {
	return this.AddById(Names.GetIdByName(name), listerFunc)
}

func (this *lister) Clear() {
	this.funcsById = make(map[uint32]map[*ListenFuncById]bool)
}

func (this *lister) DelById(id uint32, listerFunc *ListenFuncById) {
	funcsByIdOfId := this.funcsById[id]
	if funcsByIdOfId != nil {
		delete(funcsByIdOfId, listerFunc)
	}
}
