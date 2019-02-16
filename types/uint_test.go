package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUintPanics(t *testing.T) {
	// Max Uint = 1.15e+77
	// Min Uint = 0
	u1 := NewUint(0)
	u2 := OneUint()

	require.Equal(t, uint64(0), u1.Uint64())
	require.Equal(t, uint64(1), u2.Uint64())

	require.Panics(t, func() { NewUintFromBigInt(big.NewInt(-5)) })
	require.Panics(t, func() { NewUintFromString("-1") })
	require.NotPanics(t, func() {
		require.True(t, NewUintFromString("0").Equal(ZeroUint()))
		require.True(t, NewUintFromString("5").Equal(NewUint(5)))
	})

	// Overflow check
	require.True(t, u1.Add(u1).Equal(ZeroUint()))
	require.True(t, u1.Add(OneUint()).Equal(OneUint()))
	require.Equal(t, uint64(0), u1.Uint64())
	require.Equal(t, uint64(1), OneUint().Uint64())
	require.Panics(t, func() { u1.SubUint64(2) })
	require.True(t, u1.SubUint64(0).Equal(ZeroUint()))
	require.True(t, u2.Add(OneUint()).Sub(OneUint()).Equal(OneUint()))    // i2 == 1
	require.True(t, u2.Add(OneUint()).Mul(NewUint(5)).Equal(NewUint(10))) // i2 == 10
	require.True(t, NewUint(7).Div(NewUint(2)).Equal(NewUint(3)))
	require.True(t, NewUint(0).Div(NewUint(2)).Equal(ZeroUint()))
	require.True(t, NewUint(5).MulUint64(4).Equal(NewUint(20)))
	require.True(t, NewUint(5).MulUint64(0).Equal(ZeroUint()))

	// divs by zero
	require.Panics(t, func() { OneUint().Mul(ZeroUint().SubUint64(uint64(1))) })
	require.Panics(t, func() { OneUint().DivUint64(0) })
	require.Panics(t, func() { OneUint().Div(ZeroUint()) })
	require.Panics(t, func() { ZeroUint().DivUint64(0) })
	require.Panics(t, func() { OneUint().Div(ZeroUint().Sub(OneUint())) })

	require.Equal(t, uint64(0), MinUint(ZeroUint(), OneUint()).Uint64())
	require.Equal(t, uint64(1), MaxUint(ZeroUint(), OneUint()).Uint64())

	// comparison ops
	require.True(t,
		OneUint().GT(ZeroUint()),
	)
	require.False(t,
		OneUint().LT(ZeroUint()),
	)
	require.True(t,
		OneUint().GTE(ZeroUint()),
	)
	require.False(t,
		OneUint().LTE(ZeroUint()),
	)

	require.False(t, ZeroUint().GT(OneUint()))
	require.True(t, ZeroUint().LT(OneUint()))
	require.False(t, ZeroUint().GTE(OneUint()))
	require.True(t, ZeroUint().LTE(OneUint()))

	// require.Panics(t, func() { i2.Add(i2) })
	// require.Panics(t, func() { i3.Add(i3) })

	// require.Panics(t, func() { i1.Mul(i1) })
	// require.Panics(t, func() { i2.Mul(i2) })
	// require.Panics(t, func() { i3.Mul(i3) })

	// // Underflow check
	// require.NotPanics(t, func() { i2.Sub(i1) })
	// require.NotPanics(t, func() { i2.Sub(i2) })
	// require.Panics(t, func() { i2.Sub(i3) })

	// // Bound check
	// uintmax := NewUintFromBigInt(new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1)))
	// uintmin := NewUint(0)
	// require.NotPanics(t, func() { uintmax.Add(ZeroUint()) })
	// require.NotPanics(t, func() { uintmin.Sub(ZeroUint()) })
	// require.Panics(t, func() { uintmax.Add(OneUint()) })
	// require.Panics(t, func() { uintmin.Sub(OneUint()) })

	// // Division-by-zero check
	// require.Panics(t, func() { i1.Div(uintmin) })
}

/*
func TestIdentUint(t *testing.T) {
	for d := 0; d < 1000; d++ {
		n := rand.Uint64()
		i := NewUint(n)

		ifromstr, ok := NewUintFromString(strconv.FormatUint(n, 10))
		require.True(t, ok)

		cases := []uint64{
			i.Uint64(),
			i.BigInt().Uint64(),
			ifromstr.Uint64(),
			NewUintFromBigInt(new(big.Int).SetUint64(n)).Uint64(),
			NewUintWithDecimal(n, 0).Uint64(),
		}

		for tcnum, tc := range cases {
			require.Equal(t, n, tc, "Uint is modified during conversion. tc #%d", tcnum)
		}
	}
}

func TestArithUint(t *testing.T) {
	for d := 0; d < 1000; d++ {
		n1 := uint64(rand.Uint32())
		i1 := NewUint(n1)
		n2 := uint64(rand.Uint32())
		i2 := NewUint(n2)

		cases := []struct {
			ires Uint
			nres uint64
		}{
			{i1.Add(i2), n1 + n2},
			{i1.Mul(i2), n1 * n2},
			{i1.Div(i2), n1 / n2},
			{i1.AddRaw(n2), n1 + n2},
			{i1.MulRaw(n2), n1 * n2},
			{i1.DivRaw(n2), n1 / n2},
			{MinUint(i1, i2), minuint(n1, n2)},
			{MaxUint(i1, i2), maxuint(n1, n2)},
		}

		for tcnum, tc := range cases {
			require.Equal(t, tc.nres, tc.ires.Uint64(), "Uint arithmetic operation does not match with uint64 operation. tc #%d", tcnum)
		}

		if n2 > n1 {
			continue
		}

		subs := []struct {
			ires Uint
			nres uint64
		}{
			{i1.Sub(i2), n1 - n2},
			{i1.SubRaw(n2), n1 - n2},
		}

		for tcnum, tc := range subs {
			require.Equal(t, tc.nres, tc.ires.Uint64(), "Uint subtraction does not match with uint64 operation. tc #%d", tcnum)
		}
	}
}

func TestCompUint(t *testing.T) {
	for d := 0; d < 1000; d++ {
		n1 := rand.Uint64()
		i1 := NewUint(n1)
		n2 := rand.Uint64()
		i2 := NewUint(n2)

		cases := []struct {
			ires bool
			nres bool
		}{
			{i1.Equal(i2), n1 == n2},
			{i1.GT(i2), n1 > n2},
			{i1.LT(i2), n1 < n2},
		}

		for tcnum, tc := range cases {
			require.Equal(t, tc.nres, tc.ires, "Uint comparison operation does not match with uint64 operation. tc #%d", tcnum)
		}
	}
}
*/

// func TestImmutabilityAllUint(t *testing.T) {
// 	ops := []func(*Uint){
// 		func(i *Uint) { _ = i.Add(NewUint(rand.Uint64())) },
// 		func(i *Uint) { _ = i.Sub(NewUint(rand.Uint64() % i.Uint64())) },
// 		func(i *Uint) { _ = i.Mul(randuint()) },
// 		func(i *Uint) { _ = i.Div(randuint()) },
// 		func(i *Uint) { _ = i.AddRaw(rand.Uint64()) },
// 		func(i *Uint) { _ = i.SubRaw(rand.Uint64() % i.Uint64()) },
// 		func(i *Uint) { _ = i.MulRaw(rand.Uint64()) },
// 		func(i *Uint) { _ = i.DivRaw(rand.Uint64()) },
// 		func(i *Uint) { _ = i.IsZero() },
// 		func(i *Uint) { _ = i.Sign() },
// 		func(i *Uint) { _ = i.Equal(randuint()) },
// 		func(i *Uint) { _ = i.GT(randuint()) },
// 		func(i *Uint) { _ = i.LT(randuint()) },
// 		func(i *Uint) { _ = i.String() },
// 	}

// 	for i := 0; i < 1000; i++ {
// 		n := rand.Uint64()
// 		ni := NewUint(n)

// 		for opnum, op := range ops {
// 			op(&ni)

// 			require.Equal(t, n, ni.Uint64(), "Uint is modified by operation. #%d", opnum)
// 			require.Equal(t, NewUint(n), ni, "Uint is modified by operation. #%d", opnum)
// 		}
// 	}
// }

// type uintop func(Uint, *big.Int) (Uint, *big.Int)

// func uintarith(uifn func(Uint, Uint) Uint, bifn func(*big.Int, *big.Int, *big.Int) *big.Int, sub bool) uintop {
// 	return func(ui Uint, bi *big.Int) (Uint, *big.Int) {
// 		r := rand.Uint64()
// 		if sub && ui.IsUint64() {
// 			if ui.IsZero() {
// 				return ui, bi
// 			}
// 			r = r % ui.Uint64()
// 		}
// 		ur := NewUint(r)
// 		br := new(big.Int).SetUint64(r)
// 		return uifn(ui, ur), bifn(new(big.Int), bi, br)
// 	}
// }

// func uintarithraw(uifn func(Uint, uint64) Uint, bifn func(*big.Int, *big.Int, *big.Int) *big.Int, sub bool) uintop {
// 	return func(ui Uint, bi *big.Int) (Uint, *big.Int) {
// 		r := rand.Uint64()
// 		if sub && ui.IsUint64() {
// 			if ui.IsZero() {
// 				return ui, bi
// 			}
// 			r = r % ui.Uint64()
// 		}
// 		br := new(big.Int).SetUint64(r)
// 		mui := ui.ModRaw(math.MaxUint64)
// 		mbi := new(big.Int).Mod(bi, new(big.Int).SetUint64(math.MaxUint64))
// 		return uifn(mui, r), bifn(new(big.Int), mbi, br)
// 	}
// }

// func TestImmutabilityArithUint(t *testing.T) {
// 	size := 500

// 	ops := []uintop{
// 		uintarith(Uint.Add, (*big.Int).Add, false),
// 		uintarith(Uint.Sub, (*big.Int).Sub, true),
// 		uintarith(Uint.Mul, (*big.Int).Mul, false),
// 		uintarith(Uint.Div, (*big.Int).Div, false),
// 		uintarithraw(Uint.AddRaw, (*big.Int).Add, false),
// 		uintarithraw(Uint.SubRaw, (*big.Int).Sub, true),
// 		uintarithraw(Uint.MulRaw, (*big.Int).Mul, false),
// 		uintarithraw(Uint.DivRaw, (*big.Int).Div, false),
// 	}

// 	for i := 0; i < 100; i++ {
// 		uis := make([]Uint, size)
// 		bis := make([]*big.Int, size)

// 		n := rand.Uint64()
// 		ui := NewUint(n)
// 		bi := new(big.Int).SetUint64(n)

// 		for j := 0; j < size; j++ {
// 			op := ops[rand.Intn(len(ops))]
// 			uis[j], bis[j] = op(ui, bi)
// 		}

// 		for j := 0; j < size; j++ {
// 			require.Equal(t, 0, bis[j].Cmp(uis[j].BigInt()), "Int is different from *big.Int. tc #%d, Int %s, *big.Int %s", j, uis[j].String(), bis[j].String())
// 			require.Equal(t, NewUintFromBigInt(bis[j]), uis[j], "Int is different from *big.Int. tc #%d, Int %s, *big.Int %s", j, uis[j].String(), bis[j].String())
// 			require.True(t, uis[j].i != bis[j], "Pointer addresses are equal. tc #%d, Int %s, *big.Int %s", j, uis[j].String(), bis[j].String())
// 		}
// 	}
// }

// func randuint() Uint {
// 	return NewUint(rand.Uint64())
// }

// func TestSafeSub(t *testing.T) {
// 	testCases := []struct {
// 		x, y     Uint
// 		expected uint64
// 		overflow bool
// 	}{
// 		{NewUint(0), NewUint(0), 0, false},
// 		{NewUint(10), NewUint(5), 5, false},
// 		{NewUint(5), NewUint(10), 5, true},
// 		{NewUint(math.MaxUint64), NewUint(0), math.MaxUint64, false},
// 	}

// 	for i, tc := range testCases {
// 		res, overflow := tc.x.SafeSub(tc.y)
// 		require.Equal(
// 			t, tc.overflow, overflow,
// 			"invalid overflow result; x: %s, y: %s, tc: #%d", tc.x, tc.y, i,
// 		)
// 		require.Equal(
// 			t, tc.expected, res.BigInt().Uint64(),
// 			"invalid subtraction result; x: %s, y: %s, tc: #%d", tc.x, tc.y, i,
// 		)
// 	}
// }
