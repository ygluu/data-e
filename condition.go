package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 条件系统(Condition System)

import (
	"errors"
	"fmt"
)

type CondExp interface {
	Check(store *Storehouse) bool
	NameExp() string
	ValueExp(store *Storehouse) string
	EachId(fn func(id uint32))
}

type condExp struct {
	oper  string
	names INames
	left  FrmlExp
	right FrmlExp
}

func (this *condExp) EachId(fn func(id uint32)) {
	this.left.EachId(fn)
	this.right.EachId(fn)
}

func (this *condExp) NameExp() string {
	return "(" + this.left.NameExp() + this.oper + this.right.NameExp() + ")"
}

func (this *condExp) ValueExp(store *Storehouse) string {
	return "(" + this.left.ValueExp(store) + this.oper + this.right.ValueExp(store) + ")"
}

type compExpG struct {
	condExp
}

func (this *compExpG) Check(store *Storehouse) bool {
	return this.left.Float64(store) > this.right.Float64(store)
}

type compExpNG struct {
	condExp
}

func (this *compExpNG) Check(store *Storehouse) bool {
	return this.left.Float64(store) <= this.right.Float64(store)
}

type compExpL struct {
	condExp
}

func (this *compExpL) Check(store *Storehouse) bool {
	return this.left.Float64(store) < this.right.Float64(store)
}

type compExpNL struct {
	condExp
}

func (this *compExpNL) Check(store *Storehouse) bool {
	return this.left.Float64(store) >= this.right.Float64(store)
}

type compExpE struct {
	condExp
}

func (this *compExpE) Check(store *Storehouse) bool {
	return this.left.Float64(store) == this.right.Float64(store)
}

type compExpNE struct {
	condExp
}

func (this *compExpNE) Check(store *Storehouse) bool {
	return this.left.Float64(store) != this.right.Float64(store)
}

type logicExp struct {
	logic string
	names INames
	left  CondExp
	right CondExp
}

func (this *logicExp) EachId(fn func(id uint32)) {
	this.left.EachId(fn)
	this.right.EachId(fn)
}

func (this *logicExp) NameExp() string {
	return "(" + this.left.NameExp() + this.logic + this.right.NameExp() + ")"
}

func (this *logicExp) ValueExp(store *Storehouse) string {
	return "(" + this.left.ValueExp(store) + this.logic + this.right.ValueExp(store) + ")"
}

type logicExpAnd struct {
	logicExp
}

func (this *logicExpAnd) Check(store *Storehouse) bool {
	return this.left.Check(store) && this.right.Check(store)
}

type logicExpOr struct {
	logicExp
}

func (this *logicExpOr) Check(store *Storehouse) bool {
	return this.left.Check(store) || this.right.Check(store)
}

type condParser struct {
	parser
}

func (this *condParser) getCompSymbol() string {
	if this.index >= this.end {
		return ""
	}

	this.pass()

	start := this.index

	for {
		if this.index >= this.end {
			break
		}
		char := this.exp[this.index]
		switch char {
		case '>', '<', '=', '!':
			{
				this.index++
				continue
			}
		}
		break
	}

	if start == this.index {
		return ""
	}

	return string(this.exp[start:this.index])
}

func (this *condParser) getCompLeftFloat64() (ret string) {
	ret = ""
	leftParentheses := 0
	for {
		this.pass()

		if this.index >= this.end {
			break
		}

		char := this.exp[this.index]

		if char == '(' {
			leftParentheses++
		} else if char == ')' {
			leftParentheses--
		}

		if leftParentheses == 0 {
			switch char {
			case '>', '<', '=', '!':
				{
					return
				}
			}
		}

		ret = ret + string(char)
		this.index++
	}

	return
}

func (this *condParser) getCompRightFloat64() (ret string) {
	ret = ""
	leftParentheses := 0
	for {
		this.pass()

		if this.index >= this.end {
			break
		}
		char := this.exp[this.index]

		if char == '(' {
			leftParentheses++
		} else if char == ')' {
			if leftParentheses == 0 {
				return
			}
			leftParentheses--
		}

		if leftParentheses == 0 {
			switch char {
			case '&', '|':
				{
					return
				}
			}
		}

		ret = ret + string(char)
		this.index++
	}

	return
}

func (this *condParser) parseCompExp() CondExp {
	leftStr := this.getCompLeftFloat64()
	if leftStr == "" {
		this.doError("parseCompExp 缺失比较左值")
	}

	symbol := this.getCompSymbol()
	if symbol == "" {
		this.doError(fmt.Sprintf("左值“%s”后面缺失比较符", leftStr))
	}

	rightStr := this.getCompRightFloat64()
	if rightStr == "" {
		this.doError(fmt.Sprintf("比较符“%s”缺失右值", symbol))
	}

	defer func() {
		if err := recover(); err != nil {
			this.doError(fmt.Sprintf("%v", err))
		}
	}()

	left := parseFrmlExpByNames(leftStr, this.names)
	right := parseFrmlExpByNames(rightStr, this.names)

	switch symbol {
	case ">":
		{
			ret := &compExpG{}
			ret.oper = symbol
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case "<=":
		{
			ret := &compExpNG{}
			ret.oper = symbol
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case "<":
		{
			ret := &compExpL{}
			ret.oper = symbol
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case ">=":
		{
			ret := &compExpNL{}
			ret.oper = symbol
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case "=":
		{
			ret := &compExpE{}
			ret.oper = symbol
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	case "!=":
		{
			ret := &compExpNE{}
			ret.oper = symbol
			ret.names = this.names
			ret.left = left
			ret.right = right
			return ret
		}
	default:
		{
			this.doError("无效比较符：" + symbol)
		}
	}

	return nil
}

func (this *condParser) parseLogicExp(left CondExp, char rune) CondExp {
	if left == nil {
		this.doError(fmt.Sprintf("逻辑符号“%s”缺失左项", string(char)))
	}

	if this.index >= this.end {
		this.doError("缺失比较符")
	}

	i := this.index
	ret1 := this.exp[i+1]
	if ret1 != char {
		this.doError("缺失符号：" + string(char))
	}

	this.index += 2
	logic := string([]rune{char, char})

	right := this.doParse()
	if right == nil {
		this.doError(fmt.Sprintf("逻辑符号“%s”缺失右项", string(char)))
	}

	if logic == "&&" {
		ret := &logicExpAnd{}
		ret.logic = logic
		ret.names = this.names
		ret.left = left
		ret.right = right
		return ret
	} else {
		ret := &logicExpOr{}
		ret.logic = logic
		ret.names = this.names
		ret.left = left
		ret.right = right
		return ret
	}
}

func (this *condParser) doParse() (ret CondExp) {

	ret = nil

	for {
		this.pass()
		if this.index >= this.end {
			break
		}

		char := this.exp[this.index]

		switch char {
		case '(':
			{
				if ret != nil {
					this.doError(fmt.Sprintf("表达式“%s”后面有多余符号", ret.NameExp()))
				}

				this.index++
				ret = this.doParse()
				this.pass()
				if (this.index >= this.end) || (this.exp[this.index] != ')') {
					this.doError("缺失右括号")
				}
				this.index++
				continue
			}
		case '&', '|':
			{
				ret = this.parseLogicExp(ret, char)
				continue
			}
		case ')':
			{
				return
			}
		default:
			{
				if ret != nil {
					this.doError(fmt.Sprintf("表达式“%s”后面有多余符号", ret.NameExp()))
				}

				ret = this.parseCompExp()
				continue
			}
		}
	}

	return
}

func parseCondExpByNames(exp string, names INames) CondExp {
	if exp == "" {
		return nil
	}

	perser := &condParser{}
	perser.init(exp, names, "Condition")
	ret := perser.doParse()
	perser.checkEnd()

	return ret
}

func ParseCondExpByNames(exp string, names INames) (CondExp, error) {
	if err := recover(); err != nil {
		return nil, errors.New(fmt.Sprintf("%v", err))
	}
	ret := parseCondExpByNames(exp, names)
	return ret, nil
}

func ParseCondExp(exp string) (CondExp, error) {
	return ParseCondExpByNames(exp, Names)
}
