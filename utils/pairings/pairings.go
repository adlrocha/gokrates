package pairings

import (
	"math/big"
)

// G1Point Encoding
type G1Point struct {
	X *big.Int
	Y *big.Int
}

// G2Point Encoding of field elements is: X[0] * z + X[1]
type G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// P1 returns the generator for G1
func P1() G1Point {
	return G1Point{big.NewInt(1), big.NewInt(2)}
}

// P2 returns the generator of G2
func P2() G2Point {
	a, b, c, d := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	a.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634", 10)
	b.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781", 10)
	c.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531", 10)
	d.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930", 10)

	p1 := [2]*big.Int{a, b}
	p2 := [2]*big.Int{c, d}
	return G2Point{p1, p2}
}

// Negate the negation of p, i.e. p.addition(p.negate()) should be zero.
func Negate(p G1Point) G1Point {
	// The prime q in the base field F_q for G1
	q, zero, tmp := big.NewInt(0), big.NewInt(0), big.NewInt(0)
	q.SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	if p.X == zero && p.Y == zero {
		return G1Point{zero, zero}
	}
	tmp.Mod(p.Y, q)
	return G1Point{p.X, big.NewInt(0).Sub(q, tmp)}
}

// AdditionG1 The sum of two points of G1
func AdditionG1(p1 G1Point, p2 G1Point) (r G1Point) {
	if p1.X == p2.X && p1.Y == p2.Y {
		Double(p1)
	} else if p1.X == p2.X {
		return G1Point{}
	} else {
		// l = (y2 - y1) / (x2 - x1)
		suby := big.NewInt(0).Sub(p2.Y, p1.Y)
		subx := big.NewInt(0).Sub(p2.X, p1.X)
		l := big.NewInt(0).Div(suby, subx)

		// newx = l**2 - x1 - x2
		lexp := big.NewInt(0).Exp(l, big.NewInt(2), nil)
		subs := big.NewInt(0).Sub(lexp, p1.X)
		x := big.NewInt(0).Sub(subs, p2.X)

		// newy = -l * newx + l * x1 - y1
		ilx := big.NewInt(0).Mul(big.NewInt(0).Neg(l), x)
		lx := big.NewInt(0).Mul(l, p1.X)
		tmp := big.NewInt(0).Add(ilx, lx)
		y := big.NewInt(0).Sub(tmp, p1.Y)

		// assert newy == (-l * newx + l * x2 - y2)
		check1 := big.NewInt(0).Mul(l, p2.X)
		check2 := big.NewInt(0).Add(ilx, check1)
		check := big.NewInt(0).Sub(check2, p2.Y)
		if check != y {
			return G1Point{}
		}

		return G1Point{x, y}
	}
	return G1Point{}
}

// Double of two points of G1
func Double(pt G1Point) (r G1Point) {
	// l = 3 * x**2 / (2 * y)
	xexp3 := big.NewInt(0).Mul(big.NewInt(3), big.NewInt(0).Exp(pt.X, big.NewInt(2), nil))
	doubley := big.NewInt(0).Mul(big.NewInt(2), pt.Y)
	l := big.NewInt(0).Div(xexp3, doubley)

	// newx = l**2 - 2 * x
	xpow2 := big.NewInt(0).Exp(l, big.NewInt(2), nil)
	doublex := big.NewInt(0).Mul(pt.X, big.NewInt(2))
	x := big.NewInt(0).Sub(xpow2, doublex)

	// newy = -l * newx + l * x - y
	ilx := big.NewInt(0).Mul(big.NewInt(0).Neg(l), x)
	lx := big.NewInt(0).Mul(l, pt.X)
	tmp := big.NewInt(0).Add(ilx, lx)
	y := big.NewInt(0).Sub(tmp, pt.Y)

	return G1Point{x, y}
}

// ScalarMul Multiply point by scalar
func ScalarMul(p G1Point, s *big.Int) (r G1Point) {
	if s == big.NewInt(0) {
		return G1Point{}
	} else if s == big.NewInt(1) {
		return p
	} else if big.NewInt(0).Mod(s, big.NewInt(2)) != big.NewInt(0) {
		// return multiply(double(pt), n // 2)
		return ScalarMul(Double(p), big.NewInt(0).Div(s, big.NewInt(2)))
	} else {
		// return add(multiply(double(pt), int(n // 2)), pt)
		return AdditionG1(ScalarMul(Double(p), big.NewInt(0).Div(s, big.NewInt(2))), p)
	}
}

// AdditionG2 The sum of two points of G2
func AdditionG2(p1 G2Point, p2 G2Point) (r G2Point) {
	bn256g2 := Init()
	r.X[1], r.X[0], r.Y[1], r.Y[0], _ = bn256g2.ECTwistAdd(p1.X[1], p1.X[0], p1.Y[1], p1.Y[0], p2.X[1], p2.X[0], p2.Y[1], p2.Y[0])
	//TODO: Catch error
	return
}

// func pairing(p1 []G1Point, p2 []G2Point) (bool, error) {
// 	if len(p1) != len(p2) {
// 		return false, errors.New("Arguments length incorrect")
// 	}

// 	elements := len(p1)
// 	inputSize := elements * 6
// 	input := make([]*big.Int, inputSize)

// 	for i := 0; i < elements; i++ {
// 		input[i*6+0] = p1[i].X
// 		input[i*6+1] = p1[i].Y
// 		input[i*6+2] = p2[i].X[0]
// 		input[i*6+3] = p2[i].X[1]
// 		input[i*6+4] = p2[i].Y[0]
// 		input[i*6+5] = p2[i].Y[1]
// 	}

// 	res := pairing(big.NewInt(0).Add())
// 	call(g, a, v, in, insize, out, outsize)
// 	success := call(sub(gas, 2000), 8, 0, add(input, 0x20), mul(inputSize, 0x20), out, 0x20)
// }

/*
# Main miller loop
def miller_loop(Q, P):
    if Q is None or P is None:
        return FQ12.one()
    R = Q
    f = FQ12.one()
    for i in range(log_ate_loop_count, -1, -1):
        f = f * f * linefunc(R, R, P)
        R = double(R)
        if ate_loop_count & (2**i):
            f = f * linefunc(R, Q, P)
            R = add(R, Q)
    # assert R == multiply(Q, ate_loop_count)
    Q1 = (Q[0] ** field_modulus, Q[1] ** field_modulus)
    # assert is_on_curve(Q1, b12)
    nQ2 = (Q1[0] ** field_modulus, -Q1[1] ** field_modulus)
    # assert is_on_curve(nQ2, b12)
    f = f * linefunc(R, Q1, P)
    R = add(R, Q1)
    f = f * linefunc(R, nQ2, P)
    # R = add(R, nQ2) This line is in many specifications but it technically does nothing
    return f ** ((field_modulus ** 12 - 1) // curve_order)

# Pairing computation
def pairing(Q, P):
    assert is_on_curve(Q, b2)
    assert is_on_curve(P, b)
    return miller_loop(twist(Q), cast_point_to_fq12(P))
*/
/*

/// @return the result of computing the pairing check
/// e(p1[0], p2[0]) *  .... * e(p1[n], p2[n]) == 1
/// For example pairing([P1(), P1().negate()], [P2(), P2()]) should
/// return true.
function pairing(G1Point[] p1, G2Point[] p2) internal returns (bool) {
	require(p1.length == p2.length);
	uint elements = p1.length;
	uint inputSize = elements * 6;
	uint[] memory input = new uint[](inputSize);
	for (uint i = 0; i < elements; i++)
	{
		input[i * 6 + 0] = p1[i].X;
		input[i * 6 + 1] = p1[i].Y;
		input[i * 6 + 2] = p2[i].X[0];
		input[i * 6 + 3] = p2[i].X[1];
		input[i * 6 + 4] = p2[i].Y[0];
		input[i * 6 + 5] = p2[i].Y[1];
	}
	uint[1] memory out;
	bool success;
	assembly {
		success := call(sub(gas, 2000), 8, 0, add(input, 0x20), mul(inputSize, 0x20), out, 0x20)
		// Use "invalid" to make gas estimation work
		switch success case 0 { invalid() }
	}
	require(success);
	return out[0] != 0;
}


/// Convenience method for a pairing check for two pairs.
function pairingProd2(G1Point a1, G2Point a2, G1Point b1, G2Point b2) internal returns (bool) {
	G1Point[] memory p1 = new G1Point[](2);
	G2Point[] memory p2 = new G2Point[](2);
	p1[0] = a1;
	p1[1] = b1;
	p2[0] = a2;
	p2[1] = b2;
	return pairing(p1, p2);
}
/// Convenience method for a pairing check for three pairs.
function pairingProd3(
		G1Point a1, G2Point a2,
		G1Point b1, G2Point b2,
		G1Point c1, G2Point c2
) internal returns (bool) {
	G1Point[] memory p1 = new G1Point[](3);
	G2Point[] memory p2 = new G2Point[](3);
	p1[0] = a1;
	p1[1] = b1;
	p1[2] = c1;
	p2[0] = a2;
	p2[1] = b2;
	p2[2] = c2;
	return pairing(p1, p2);
}
/// Convenience method for a pairing check for four pairs.
function pairingProd4(
		G1Point a1, G2Point a2,
		G1Point b1, G2Point b2,
		G1Point c1, G2Point c2,
		G1Point d1, G2Point d2
) internal returns (bool) {
	G1Point[] memory p1 = new G1Point[](4);
	G2Point[] memory p2 = new G2Point[](4);
	p1[0] = a1;
	p1[1] = b1;
	p1[2] = c1;
	p1[3] = d1;
	p2[0] = a2;
	p2[1] = b2;
	p2[2] = c2;
	p2[3] = d2;
	return pairing(p1, p2);
}
}
*/
