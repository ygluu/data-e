package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 解析器基类(base parser)

import (
	"fmt"
)

type parser struct {
	exp       []rune
	end       int
	index     int
	line      int
	lineStart int
	errFlag   string
	names     INames
}

func (this *parser) init(exp string, names INames, errFlag string) {
	this.exp = []rune(exp)
	this.end = len(this.exp)
	this.index = 0
	this.line = 1
	this.lineStart = 0
	this.names = names
	this.errFlag = errFlag
}

func (this *parser) Name() INames {
	return this.names
}

func (this *parser) toLineEnd() {
	for {
		if this.index >= this.end {
			return
		}
		char := this.exp[this.index]
		this.index++
		if char == '\n' {
			this.line++
			this.lineStart = this.index
			return
		}
	}
}

func (this *parser) pass() (first int, ret bool) {
	first = -1
	for {
		if this.index >= this.end {
			break
		}
		char := this.exp[this.index]
		if char == '\n' {
			this.line++
			this.lineStart = this.index
		}
		if (char <= 32) || (char == 127) {
			ret = true
			if first == -1 {
				first = this.index
			}
			this.index++
		} else {
			break
		}
	}
	return
}

func (this *parser) doError(msg string) {
	if (this.index) < this.end {
		panic(fmt.Sprintf("[dmp]%s => %s，在第%d行第%d列字符[%s]附近，源串：%s",
			this.errFlag, msg, this.line, this.index-this.lineStart, string(this.exp[this.index:]), string(this.exp)))
	} else {
		panic(fmt.Sprintf("[dmp]%s => %s，在第%d行第%d列附近，源串：%s",
			this.errFlag, msg, this.line, this.index-this.lineStart, string(this.exp)))
	}

}

func (this *parser) checkEnd() {
	this.pass()
	if this.index < this.end {
		this.doError("表达式后面有多余字符")
	}
}

func (this *parser) getNameId(name string) uint32 {
	ret := this.names.GetIdByName(name)
	if ret == 0 {
		this.doError("无效数据名：" + name)
	}
	return ret
}

func (this *parser) readFuncParams() (ret []string) {
	this.index++
	start := -1
	leftParentheses := 0

	getParam := func() {
		if start == -1 {
			return
		}
		param := string(this.exp[start:this.index])
		start = -1
		if param != "" {
			ret = append(ret, param)
		}
	}

	for {
		this.pass()
		if this.index >= this.end {
			break
		}

		if start == -1 {
			start = this.index
		}

		char := this.exp[this.index]
		switch char {
		case '(':
			{
				leftParentheses++
			}
		case ')':
			{
				leftParentheses--
				if leftParentheses == 0 {
					this.index++
					getParam()
					continue
				} else if leftParentheses < 0 {
					getParam()
					return
				}
			}
		case paramSeparator:
			{
				if start == -1 {
					this.doError("此处出现多余函数参数分隔符：" + string(paramSeparator))
				}
				if leftParentheses == 0 {
					getParam()
				}
			}
		}

		this.index++
	}

	return
}
