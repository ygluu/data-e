package de

/*******************************************************************************

//     Data Engine 数据引擎 (data-e)

//        Author: Yigui Lu (卢益贵)
// Contact WX/QQ: 48092788
//          Blog: https://blog.csdn.net/guestcode
//   Creation by: 2018-2020

*******************************************************************************/

// 数学函数库(Math Function Library)

import (
	"math"
)

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Trunc(params[0].Float64(store))
	}
	RegisterFunc("Trunc", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.NaN()
	}
	RegisterFunc("NaN", funcExec, 0)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Inf(int(params[0].Float64(store)))
	}
	RegisterFunc("Inf", funcExec, 0)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Cbrt(params[0].Float64(store))
	}
	RegisterFunc("Cbrt", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Sqrt(params[0].Float64(store))
	}
	RegisterFunc("Sqrt", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Hypot(params[0].Float64(store), params[1].Float64(store))
	}
	RegisterFunc("Hypot", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Sin(params[0].Float64(store))
	}
	RegisterFunc("Sin", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Cos(params[0].Float64(store))
	}
	RegisterFunc("Cos", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Tan(params[0].Float64(store))
	}
	RegisterFunc("Tan", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Log(params[0].Float64(store))
	}
	RegisterFunc("Log", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Log2(params[0].Float64(store))
	}
	RegisterFunc("Log2", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Log10(params[0].Float64(store))
	}
	RegisterFunc("Log10", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Log1p(params[0].Float64(store))
	}
	RegisterFunc("Log1p", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Logb(params[0].Float64(store))
	}
	RegisterFunc("Logb", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return float64(math.Ilogb(params[0].Float64(store)))
	}
	RegisterFunc("Ilogb", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Abs(params[0].Float64(store))
	}
	RegisterFunc("Abs", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Floor(params[0].Float64(store))
	}
	RegisterFunc("Floor", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Ceil(params[0].Float64(store))
	}
	RegisterFunc("Ceil", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Mod(params[0].Float64(store), params[1].Float64(store))
	}
	RegisterFunc("Mod", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Pow(params[0].Float64(store), params[1].Float64(store))
	}
	RegisterFunc("Pow", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Copysign(params[0].Float64(store), params[1].Float64(store))
	}
	RegisterFunc("Copysign", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Nextafter(params[0].Float64(store), params[1].Float64(store))
	}
	RegisterFunc("Nextafter", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Remainder(params[0].Float64(store), params[1].Float64(store))
	}
	RegisterFunc("Remainder", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Dim(params[0].Float64(store), params[1].Float64(store))
	}
	RegisterFunc("Dim", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Asin(params[0].Float64(store))
	}
	RegisterFunc("Asin", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Acos(params[0].Float64(store))
	}
	RegisterFunc("Acos", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Atan2(params[0].Float64(store), params[1].Float64(store))
	}
	RegisterFunc("Atan2", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Atan(params[0].Float64(store))
	}
	RegisterFunc("Atan", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Sinh(params[0].Float64(store))
	}
	RegisterFunc("Sinh", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Cosh(params[0].Float64(store))
	}
	RegisterFunc("Cosh", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Tanh(params[0].Float64(store))
	}
	RegisterFunc("Tanh", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Asinh(params[0].Float64(store))
	}
	RegisterFunc("Asinh", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Acosh(params[0].Float64(store))
	}
	RegisterFunc("Acosh", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Atanh(params[0].Float64(store))
	}
	RegisterFunc("Atanh", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Ldexp(params[0].Float64(store), int(params[1].Float64(store)))
	}
	RegisterFunc("Ldexp", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Exp(params[0].Float64(store))
	}
	RegisterFunc("Exp", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Exp2(params[0].Float64(store))
	}
	RegisterFunc("Exp2", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Expm1(params[0].Float64(store))
	}
	RegisterFunc("Expm1", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Pow10(int(params[0].Float64(store)))
	}
	RegisterFunc("Pow10", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Gamma(params[0].Float64(store))
	}
	RegisterFunc("Gamma", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Erf(params[0].Float64(store))
	}
	RegisterFunc("Erf", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Erfc(params[0].Float64(store))
	}
	RegisterFunc("Erfc", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.J0(params[0].Float64(store))
	}
	RegisterFunc("J0", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.J1(params[0].Float64(store))
	}
	RegisterFunc("J1", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Jn(int(params[0].Float64(store)), params[1].Float64(store))
	}
	RegisterFunc("Jn", funcExec, 2)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Y0(params[0].Float64(store))
	}
	RegisterFunc("Y0", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Y1(params[0].Float64(store))
	}
	RegisterFunc("Y1", funcExec, 1)
}

func init() {
	funcExec := func(store *Storehouse, params []FrmlExp) float64 {
		return math.Yn(int(params[0].Float64(store)), params[1].Float64(store))
	}
	RegisterFunc("Yn", funcExec, 2)
}

func init() {
	typ := "mathconst"
	Names.RegisterName(typ, "E", 1)
	Names.RegisterName(typ, "PI", 2)
	Names.RegisterName(typ, "PHI", 3)
	Names.RegisterName(typ, "SQRT2", 4)
	Names.RegisterName(typ, "SQRTE", 5)
	Names.RegisterName(typ, "SQRTPI", 6)
	Names.RegisterName(typ, "SQRTPHI", 7)
	Names.RegisterName(typ, "LN2", 8)
	Names.RegisterName(typ, "LOG2E", 9)
	Names.RegisterName(typ, "LOG10", 10)
	Names.RegisterName(typ, "LOG10E", 11)

	getMathConst := func(store *Storehouse, id uint32) float64 {
		values := []float64{0, math.E, math.Pi, math.Phi, math.Sqrt2, math.SqrtE, math.SqrtPi, math.SqrtPhi, math.Ln2, math.Log2E, math.Ln10, math.Log10E}
		rid := id & RawIDMark
		if rid >= uint32(len(values)) {
			return 0
		}
		return values[rid]
	}
	Names.RegisterGetFuncByType(typ, getMathConst)
}
