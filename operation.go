package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 四则运算集合(Operation Set)

import (
	"errors"
	"fmt"
	"strconv"
)

type OperExp interface {
	Exec(store *Storehouse) bool
	NameExp() string
	ValueExp(store *Storehouse) string
	DestId() uint32
}

type operExp struct {
	names  INames
	nameId uint32
	value  FrmlExp
}

func (this *operExp) DestId() uint32 {
	return this.nameId
}

type incOperExp struct {
	operExp
}

func (this *incOperExp) NameExp() string {
	return this.names.GetNameById(this.nameId) + "+=" + this.value.NameExp()
}

func (this *incOperExp) ValueExp(store *Storehouse) string {
	return strconv.FormatFloat(this.value.Float64(store), 'f', -1, 64) + "+=" + this.value.ValueExp(store)
}

func (this *incOperExp) Exec(store *Storehouse) bool {
	store.Oper(this.nameId, OS_INC, this.value.Float64(store))
	return false
}

type decOperExp struct {
	operExp
}

func (this *decOperExp) NameExp() string {
	return this.names.GetNameById(this.nameId) + "-=" + this.value.NameExp()
}

func (this *decOperExp) ValueExp(store *Storehouse) string {
	return strconv.FormatFloat(this.value.Float64(store), 'f', -1, 64) + "-=" + this.value.ValueExp(store)
}

func (this *decOperExp) Exec(store *Storehouse) bool {
	store.Oper(this.nameId, OS_DEC, this.value.Float64(store))
	return false
}

type mulOperExp struct {
	operExp
}

func (this *mulOperExp) NameExp() string {
	return this.names.GetNameById(this.nameId) + "*=" + this.value.NameExp()
}

func (this *mulOperExp) ValueExp(store *Storehouse) string {
	return strconv.FormatFloat(this.value.Float64(store), 'f', -1, 64) + "*=" + this.value.ValueExp(store)
}

func (this *mulOperExp) Exec(store *Storehouse) bool {
	store.Oper(this.nameId, OS_MUL, this.value.Float64(store))
	return false
}

type divOperExp struct {
	operExp
}

func (this *divOperExp) NameExp() string {
	return this.names.GetNameById(this.nameId) + "/=" + this.value.NameExp()
}

func (this *divOperExp) ValueExp(store *Storehouse) string {
	return strconv.FormatFloat(this.value.Float64(store), 'f', -1, 64) + "/=" + this.value.ValueExp(store)
}

func (this *divOperExp) Exec(store *Storehouse) bool {
	store.Oper(this.nameId, OS_DIV, this.value.Float64(store))
	return false
}

type setOperExp struct {
	operExp
}

func (this *setOperExp) NameExp() string {
	return this.names.GetNameById(this.nameId) + "=" + this.value.NameExp()
}

func (this *setOperExp) ValueExp(store *Storehouse) string {
	return strconv.FormatFloat(this.value.Float64(store), 'f', -1, 64) + "=" + this.value.ValueExp(store)
}

func (this *setOperExp) Exec(store *Storehouse) bool {
	store.Oper(this.nameId, OS_SET, this.value.Float64(store))
	return false
}

type OperSet interface {
	NameExp() string
	ValueExp(store *Storehouse) string
	Opers() []OperExp
}

type operSet []OperExp

func (this operSet) Opers() []OperExp {
	return this
}

func (this operSet) NameExp() (ret string) {
	ret = ""
	for _, step := range this {
		if ret == "" {
			ret = step.NameExp()
		} else {
			ret = ret + string(stepSeparator) + step.NameExp()
		}
	}
	return
}

func (this operSet) ValueExp(store *Storehouse) (ret string) {
	ret = ""
	for _, step := range this {
		if ret == "" {
			ret = step.ValueExp(store)
		} else {
			ret = ret + string(stepSeparator) + step.ValueExp(store)
		}
	}
	return
}

type operParser struct {
	parser
	buildExp func(symbol rune, nameStart, nameEnd int) (ret OperExp)
}

func (this *operParser) extractNameValue(nameStart, nameEnd int) (nameId uint32, value FrmlExp) {

	name := string(this.exp[nameStart:nameEnd])
	if name == "" {
		this.doError(fmt.Sprintf("符号“%s”前缺失数据名", string(this.exp[this.index])))
	}

	nameId = this.names.GetIdByName(name)
	if nameId == 0 {
		this.doError("无效数据名：" + name)
	}

	this.index++
	this.pass()
	valueExp := ""
	isValueEnd := false
	hasOper := false
	hasSpace := false
	leftParentheses := 0

	defer func() {
		if valueExp == "" {
			this.doError(fmt.Sprintf("数据“%s”缺失值表达式", name))
		}
		defer func() {
			if err := recover(); err != nil {
				this.doError(fmt.Sprintf("%v", err))
			}
		}()
		value = parseFrmlExpByNames(valueExp, this.names)
	}()

	for {
		if this.index >= this.end {
			//valueExp = string(this.exp[valueStart:this.index])
			break
		}

		char := this.exp[this.index]

		if (char == '/') && (this.index < this.end-1) && (this.exp[this.index+1] == '/') {
			this.toLineEnd()
			this.index++
			continue
		}

		switch char {
		case '(':
			{
				leftParentheses++
				valueExp = valueExp + string(char)
			}
		case ')':
			{
				if leftParentheses == 0 {
					this.doError("多余右括号")
				}
				leftParentheses--
				valueExp = valueExp + string(char)
			}
		case stepSeparator:
			{
				isValueEnd = true
				hasSpace = false
			}
		case ' ':
			{
				isValueEnd = true
				hasSpace = true
			}
		case '\n':
			{
				hasSpace = false
				isValueEnd = true
				hasOper = false
				this.line++
				this.lineStart = this.index + 1
			}
		case '+', '-', '*', '/', '&', '<', '>', '=', '!', ':', ',':
			{
				hasOper = true
				hasSpace = false
				valueExp = valueExp + string(char)
			}
		default:
			{
				if (isValueEnd) && (!hasOper) && (leftParentheses == 0) {
					this.index--
					return
				}
				if (char > 32) && (char != 127) {
					if hasSpace && (!hasOper) {
						this.doError("此处出现多余空格")
					}
					valueExp = valueExp + string(char)
				}
				isValueEnd = false
				hasOper = false
				hasSpace = false
			}
		}

		this.index++
	}

	return
}

func (this *operParser) doBuildExp(symbol rune, nameStart, nameEnd int) (ret OperExp) {

	checkEqu := func() {
		if (this.index >= this.end-1) || (this.exp[this.index+1] != '=') {
			this.doError(fmt.Sprintf("符号“%s”后面缺失等号“=”", string(symbol)))
		}
		this.index++
	}

	if nameEnd == -1 {
		nameEnd = this.index
	}

	switch symbol {
	case '+':
		{
			checkEqu()
			nameId, value := this.extractNameValue(nameStart, nameEnd)
			ret := &incOperExp{}
			ret.names = this.names
			ret.nameId = nameId
			ret.value = value
			return ret
		}
	case '-':
		{
			checkEqu()
			nameId, value := this.extractNameValue(nameStart, nameEnd)
			ret := &decOperExp{}
			ret.names = this.names
			ret.nameId = nameId
			ret.value = value
			return ret
		}
	case '*':
		{
			checkEqu()
			nameId, value := this.extractNameValue(nameStart, nameEnd)
			ret := &mulOperExp{}
			ret.names = this.names
			ret.nameId = nameId
			ret.value = value
			return ret
		}
	case '/':
		{
			checkEqu()
			nameId, value := this.extractNameValue(nameStart, nameEnd)
			ret := &divOperExp{}
			ret.names = this.names
			ret.nameId = nameId
			ret.value = value
			return ret
		}
	case '=', ':':
		{
			nameId, value := this.extractNameValue(nameStart, nameEnd)
			ret := &setOperExp{}
			ret.names = this.names
			ret.nameId = nameId
			ret.value = value
			return ret
		}
	}
	return
}

func (this *operParser) doParse(onGetOperExp func(exp OperExp)) {

	nameStart := -1
	nameEnd := -1

	for {
		firstSpace, hasSpace := this.pass()
		if (hasSpace) && (nameStart != -1) {
			if nameEnd == -1 {
				nameEnd = firstSpace
			} else {
				this.doError(fmt.Sprintf("名字“%s”后面出现多余字符", string(this.exp[nameStart:nameEnd])))
			}
		}

		if this.index >= this.end {
			return
		}

		char := this.exp[this.index]

		if (char == stepSeparator) && (char == ' ') {
			this.index++
			continue
		}

		if (char == '/') && (this.index < this.end-1) && (this.exp[this.index+1] == '/') {
			this.toLineEnd()
			continue
		}

		if nameStart == -1 {
			nameStart = this.index
			nameEnd = -1
		}

		if nameStart != -1 {
			exp := this.buildExp(char, nameStart, nameEnd)
			if exp != nil {
				nameStart = -1
				nameEnd = -1
				onGetOperExp(exp)
			}
		}

		this.index++
	}
}

func parseOperExp(exp string) OperSet {
	ret := operSet{}

	if exp == "" {
		return ret
	}

	onGetOperExp := func(exp OperExp) {
		ret = append(ret, exp)
	}

	perser := &operParser{}
	perser.buildExp = perser.doBuildExp
	perser.init(exp, Names, "OperSet")
	perser.doParse(onGetOperExp)
	perser.checkEnd()

	if len(ret) == 0 {
		perser.doError("无效运算表达式")
	}

	return &ret
}

func ParseOperExp(exp string) (OperSet, error) {
	if err := recover(); err != nil {
		return nil, errors.New(fmt.Sprintf("%v", err))
	}
	ret := parseOperExp(exp)
	return ret, nil
}
