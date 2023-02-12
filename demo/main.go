package main

import (
	"fmt"
	"log"

	"lib/data-e"
)

func main() {
	log.Println("可配置化数学模型库演示样例.....")

	// 实际应用中需视场景采用运算集合还是流程步骤（或者组合）
	// 1、注册所有数据名（可从配置加载）
	de.Names.RegisterName(de.StrTypeName, "卢益贵", 0)
	de.Names.RegisterName("", "名字", 0)
	de.Names.RegisterName("", "钱包", 0)
	de.Names.RegisterName("", "年龄", 0)
	de.Names.RegisterName("", "生年", 0)
	de.Names.RegisterName("", "年份", 0)

	// 2.解析表达式

	// 运算集合表达式解析（可从配置加载）
	strOper := "生年=2002 年份=2022 年龄=年份-生年 名字=卢益贵"
	oper, err := de.ParseOperExp(strOper)

	if err != nil {
		panic(err.Error())
	}

	// 步骤流程表达式解析（可从配置加载），在流程中出现的名称没有在de.Names.RegisterName注册的视为临时变量
	strProc := `
		return(名字!=卢益贵,0)
		
		钱包  = 年龄*10
		
		// “当前值”是临时变量
		当前值 = If(年龄>=10, 年龄 * 0.2, 0)
		钱包 += 当前值
		
		// 可分行表达
		钱包 = 钱包 + If(年龄>=15, 
		  	年龄 * (Random(3,6) / 10), 0)
				
		// 小于18岁就这点钱了
		Return(年龄<18, 钱包)
		
		// 18岁以上能不能发财全靠运气，临时变量：狗屎运、财神运
		狗屎运=Random(5000)
		财神运=Random(10000, 20000)
		钱包 = 钱包 + (财神运+狗屎运)*(钱包*钱包)/(年龄*年龄) * Sin(Mod(年龄, 3)) * PI
		Return(钱包)
	`
	proc, err := de.ParseProcExp(strProc)
	if err != nil {
		panic(err.Error())
	}
	log.Println("解析后的流程步骤表达式：", "\r\n"+proc.NameExp(), "\r\n")

	// 3.运行时执行运算和步骤流程

	cond, err := de.ParseCondExp("钱包>1000000&&年龄>10")
	if err != nil {
		panic(err.Error())
	}
	onCond := func(store *de.Storehouse, extParam uintptr) {
		log.Println("**************************************************************")
		log.Println("当" + cond.NameExp() + "时打印：【钱包极度膨胀】")
		log.Println("**************************************************************")
	}

	store := de.NewStorehouse(nil)

	workstat := de.NewWorkstat()
	// 设置监听条件：当条件满足时，执行OnCond
	workstat.ListenCond(store, cond, onCond, 0)

	// 执行运算集合（初始化基础数据：生年=2002 年份=2022 年龄=年份-生年 名字=卢益贵）
	workstat.ExecOper(store, oper, false)
	log.Println("运算集合名字字符表达式：", "\r\n"+oper.NameExp(), "\r\n")

	// 设置输出流程步骤日志可清晰输出每步数值情况(日志执行Random和步骤执行Random打印结果会不一致)
	de.SetOutStepLog(true)
	log.Println("执行流程步骤并输出单步日志：")

	// 执行流程步骤并返回结果
	money := workstat.ExecProc(store, proc, false)
	log.Println("")
	log.Println("流程步骤返回值1：钱包 =", fmt.Sprintf("%f", money), "\r\n")

	log.Println("将数据名和各类表达式从配置文件加载即可实现可配置化")
	log.Println("测试结束，若有Bug待解")
}
