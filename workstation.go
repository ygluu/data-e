package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 数据工作站(Data Workstation)

type ListenFuncByCond = func(store *Storehouse, extParam uintptr)

type condLister struct {
	fnOfCond ListenFuncByCond
	fnsOfId  map[uint32]*ListenFuncById
	cond     CondExp
	extParam uintptr
}

func (this *condLister) onDataChg(store *Storehouse, id uint32, operSymbol OperSymbol, value float64) {
	if this.cond.Check(store) {
		defer func() {
			if err := recover(); err != nil {
				writeLog("[dmp]condLister.onDataChg => 条件监听回调异常：%v，条件：%s", err, this.cond.NameExp())
			}
		}()
		this.fnOfCond(store, uintptr(this.extParam))
	}
}

type workstatLog struct {
	produceIdSaveFlagOfOperSet map[OperSet]bool
	produceIdSaveFlagOfProcExp map[ProcExp]bool
	idSaveFlagOfProduceId      map[uint32]bool
	listerOfCond               map[CondExp]*condLister
}

func newWorkstatLog() *workstatLog {
	return &workstatLog{
		produceIdSaveFlagOfOperSet: make(map[OperSet]bool),
		produceIdSaveFlagOfProcExp: make(map[ProcExp]bool),
		idSaveFlagOfProduceId:      make(map[uint32]bool),
		listerOfCond:               make(map[CondExp]*condLister),
	}
}

type Workstat struct {
}

func NewWorkstat() *Workstat {
	return &Workstat{}
}

var WorkStat = NewWorkstat()

func (this *Workstat) ResetMyData(store *Storehouse) {
	wLog := store.workstatLog[this]
	if wLog != nil {
		for id, _ := range wLog.idSaveFlagOfProduceId {
			store.ResetById(id)
		}
	}
}

func (this *Workstat) myLog(store *Storehouse) *workstatLog {
	ret := store.workstatLog[this]
	if ret == nil {
		ret = newWorkstatLog()
		store.workstatLog[this] = ret
	}
	return ret
}

func (this *Workstat) CancelCondListen(store *Storehouse, cond CondExp) {
	wLog := this.myLog(store)
	lister := wLog.listerOfCond[cond]
	if lister == nil {
		return
	}
	for id, fn := range lister.fnsOfId {
		store.Lister.DelById(id, fn)
	}
	delete(wLog.listerOfCond, cond)
}

func (this *Workstat) ListenCond(store *Storehouse, cond CondExp, fn ListenFuncByCond, extParam uintptr) {
	ids := make(map[uint32]uint32)
	cond.EachId(func(id uint32) {
		ids[id] = id
	})
	len := len(ids)
	if len == 0 {
		writeLog("[dmp]Workstat.ListenCond => 监听的条件无实际数据名：" + cond.NameExp())
	}

	wLog := this.myLog(store)
	lister := wLog.listerOfCond[cond]

	if lister == nil {
		lister := &condLister{
			fnOfCond: fn,
			cond:     cond,
			fnsOfId:  make(map[uint32]*ListenFuncById),
			extParam: extParam,
		}
		wLog.listerOfCond[cond] = lister
		for _, id := range ids {
			lister.fnsOfId[id] = store.Lister.AddById(id, lister.onDataChg)
		}
	}
}

func (this *Workstat) ExecOper(store *Storehouse, exp OperSet, recordProduce bool) {
	opers := exp.Opers()

	if recordProduce {
		wLog := this.myLog(store)
		if !wLog.produceIdSaveFlagOfOperSet[exp] {
			wLog.produceIdSaveFlagOfOperSet[exp] = true
			for _, oper := range opers {
				id := oper.DestId()
				if id != 0 {
					wLog.idSaveFlagOfProduceId[id] = true
				}
			}
		}
	}

	for _, step := range opers {
		step.Exec(store)
	}
}

func (this *Workstat) ExecProc(store *Storehouse, exp ProcExp, recordProduce bool) float64 {
	exp.LoadFrom(store)

	steps := exp.Steps()
	if recordProduce {
		wLog := this.myLog(store)
		if !wLog.produceIdSaveFlagOfProcExp[exp] {
			wLog.produceIdSaveFlagOfProcExp[exp] = true
			for _, oper := range steps {
				id := oper.DestId()
				if id != 0 {
					wLog.idSaveFlagOfProduceId[id] = true
				}
			}
		}
	}

	procStore := exp.MyStore()
	for i, step := range steps {
		if outStepLog {
			writeLog("[Step%d].名字表达式:%s", i+1, step.NameExp())
			writeLog("[Step%d].数值表达式:%s", i+1, step.ValueExp(procStore))
		}
		if step.Exec(procStore) {
			break
		}
	}

	exp.SaveTo(store)
	return procStore.Get(returnID)
}
