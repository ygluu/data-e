package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

import (
	"log"
)

var stepSeparator = ' '
var paramSeparator = ','

var outStepLog = false
var writeLog = log.Printf

func SetLogFunc(value func(format string, v ...any)) {
	writeLog = value
}

func SetOutStepLog(value bool) {
	outStepLog = value
}

func SetStepSeparator(sep rune) {
	stepSeparator = sep
}

func SetParamSeparator(sep rune) {
	paramSeparator = sep
}
