package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 四则运算公式系统(Formula System)

import (
	"errors"
	"fmt"
	"strconv"
)

type FrmlExp interface {
	NameExp() string
	ValueExp(store *Storehouse) string
	Float64(store *Storehouse) float64
	EachId(fn func(id uint32))
}

const ERR_COMPENSATION = float64(0.000000000000000000001)

type constExp struct {
	value float64
}

func (this *constExp) EachId(fn func(id uint32)) {
}

func (this *constExp) NameExp() string {
	return strconv.FormatFloat(this.value, 'f', -1, 64)
}

func (this *constExp) ValueExp(store *Storehouse) string {
	return strconv.FormatFloat(this.value, 'f', -1, 64)
}

func (this *constExp) Float64(store *Storehouse) float64 {
	return this.value
}

type idenExp struct {
	names INames
	id    uint32
	name  string
}

func (this *idenExp) EachId(fn func(id uint32)) {
	fn(this.id)
}

func (this *idenExp) NameExp() string {
	return this.names.GetNameById(this.id)
}

func (this *idenExp) ValueExp(store *Storehouse) string {
	return strconv.FormatFloat(store.Get(this.id), 'f', -1, 64)
}

func (this *idenExp) Float64(store *Storehouse) float64 {
	return store.Get(this.id)
}

type frmlExp struct {
	names INames
	oper  string
	left  FrmlExp
	right FrmlExp
	value float64
}

func (this *frmlExp) NameExp() string {
	return "(" + this.left.NameExp() + this.oper + this.right.NameExp() + ")"
}

func (this *frmlExp) ValueExp(store *Storehouse) string {
	return "(" + this.left.ValueExp(store) + this.oper + this.right.ValueExp(store) + ")"
}

func (this *frmlExp) EachId(fn func(id uint32)) {
	this.left.EachId(fn)
	this.right.EachId(fn)
}

type incExp struct {
	frmlExp
}

func (this *incExp) Float64(store *Storehouse) float64 {
	return this.left.Float64(store) + this.right.Float64(store)
}

type decExp struct {
	frmlExp
}

func (this *decExp) Float64(store *Storehouse) float64 {
	return this.left.Float64(store) - this.right.Float64(store)
}

type mulExp struct {
	frmlExp
}

func (this *mulExp) Float64(store *Storehouse) float64 {
	return this.left.Float64(store) * this.right.Float64(store)
}

type divExp struct {
	frmlExp
}

func (this *divExp) Float64(store *Storehouse) float64 {
	rv := this.right.Float64(store)
	if rv == 0 {
		return this.left.Float64(store)
	}
	return this.left.Float64(store) / rv
}

type modExp struct {
	frmlExp
}

func (this *modExp) Float64(store *Storehouse) float64 {
	if this.right.Float64(store) == 0 {
		return this.left.Float64(store)
	}
	return float64(int64(this.left.Float64(store)) % int64(this.right.Float64(store)))
}

type frmlParser struct {
	parser
}

func (this *frmlParser) parseValueExp(value string, isNum bool) FrmlExp {
	if isNum {
		v, e := strconv.ParseFloat(value, 64)
		if e != nil {
			this.doError("无效数值：" + value)
		}
		return &constExp{value: v}
	}

	id := this.names.GetIdByName(value)
	if id == 0 {
		this.doError("无效值名：" + value)
	}
	return &idenExp{id: id, name: value, names: this.names}
}

func (this *frmlParser) parseFunc(fp FuncParser) FrmlExp {
	defer func() {
		if err := recover(); err != nil {
			this.doError(fmt.Sprintf("%v", err))
		}
	}()
	return fp.doParse(this.names, this.readFuncParams())
}

func (this *frmlParser) getValueExp() (ret FrmlExp) {
	ret = nil
	str := ""
	leftParentheses := 0
	isNum := true
	count := 0

	defer func() {
		if ret == nil {
			ret = this.parseValueExp(str, isNum)
		}
	}()

	for {
		this.pass()

		if this.index >= this.end {
			break
		}

		char := this.exp[this.index]

		if char == '(' {
			if str != "" {
				fp := getFuncParser(str)
				if fp == nil {
					this.doError("未定义的函数：" + str)
				}
				ret = this.parseFunc(fp)
				this.index++
				return
			}

			if count == 0 {
				this.index++
				ret = this.doParse()
				this.pass()
				if (this.index >= this.end) || (this.exp[this.index] != ')') {
					this.doError("缺失右括号")
				}
				this.index++
				return
			}
			leftParentheses++
		} else if char == ')' {
			if leftParentheses == 0 {
				return
			}
			leftParentheses--
		}

		if leftParentheses == 0 {
			switch char {
			case '+', '*', '/', '%':
				{
					return
				}
			case '-':
				{
					if str != "" {
						return
					}
				}
			case '<', '>', '=', '!', ':':
				{
					this.doError(fmt.Sprintf("名称不能含有字符“%s”", string(char)))
				}
			}
		}

		if (char < '0' || char > '9') && (char != '.') && (char != '-') {
			isNum = false
		}

		count++
		str = str + string(char)
		this.index++
	}

	return
}

func (this *frmlParser) buildExp(left FrmlExp, symbol rune, right FrmlExp) FrmlExp {
	switch symbol {
	case '+':
		{
			ret := &incExp{}
			ret.oper = string(symbol)
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case '-':
		{
			ret := &decExp{}
			ret.oper = string(symbol)
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case '*':
		{
			ret := &mulExp{}
			ret.oper = string(symbol)
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case '/':
		{
			ret := &divExp{}
			ret.oper = string(symbol)
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case '%':
		{
			ret := &modExp{}
			ret.oper = string(symbol)
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	default:
		{
			return nil
		}
	}
}

func (this *frmlParser) doParse() (ret FrmlExp) {

	left := FrmlExp(nil)
	right := FrmlExp(nil)
	symbol := rune(0)

	defer func() {
		if left == nil {
			this.doError("无效表达式")
		}
		if symbol == 0 {
			ret = left
		} else if right == nil {
			this.doError(fmt.Sprintf("运算符“%s”缺失右项", string(symbol)))
		} else {
			ret = this.buildExp(left, symbol, right)
		}
	}()

	operPriority := func(symbol rune) int {
		switch symbol {
		case '+', '-':
			{
				return 0
			}
		default:
			{
				return 1
			}
		}
	}

	setRight := func(char rune, exp FrmlExp) {
		if right == nil {
			symbol = char
			right = exp
		} else {
			if operPriority(char) <= operPriority(symbol) {
				left = this.buildExp(left, symbol, right)
				symbol = char
				right = exp
			} else {
				right = this.buildExp(right, char, exp)
			}
		}
	}

	for {
		this.pass()
		if this.index >= this.end {
			break
		}

		char := this.exp[this.index]

		switch char {
		case '(':
			{
				this.index++
				cur := this.doParse()
				if cur == nil {
					this.doError(fmt.Sprintf("运算符“%s”缺失右值", string(char)))
				}

				if left == nil {
					left = cur
				} else {
					setRight(char, cur)
				}

				this.pass()
				if (this.index >= this.end) || (this.exp[this.index] != ')') {
					this.doError("缺失右括号")
				}
				this.index++
				continue
			}
		case '+', '-', '*', '/', '%':
			{
				if left == nil {
					if char != '-' {
						this.doError(fmt.Sprintf("运算符“%s”缺失左值", string(char)))
					}
					// 负数
					left = this.getValueExp()
					continue
				}

				this.index++
				cur := this.getValueExp()
				if cur == nil {
					this.doError(fmt.Sprintf("运算符“%s”缺失右值", string(char)))
				}
				setRight(char, cur)
				continue
			}
		case ')':
			{
				return
			}
		default:
			{
				if left != nil {
					this.doError("此处有多余符号")
				}
				left = this.getValueExp()
				continue
			}
		}
	}

	return
}

func parseFrmlExpByNames(exp string, names INames) FrmlExp {
	if exp == "" {
		return nil
	}

	perser := &frmlParser{}
	perser.init(exp, names, "Formula")
	ret := perser.doParse()
	perser.checkEnd()

	return ret
}

func ParseFrmlExpByNames(exp string, names INames) (FrmlExp, error) {
	if err := recover(); err != nil {
		return nil, errors.New(fmt.Sprintf("%v", err))
	}
	ret := parseFrmlExpByNames(exp, names)
	return ret, nil
}

func ParseFrmlExp(exp string) (FrmlExp, error) {
	return ParseFrmlExpByNames(exp, Names)
}
