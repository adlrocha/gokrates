package pairings

import (
	"fmt"
	"math/big"
)

//BN256G2CURVE structure
type BN256G2CURVE struct {
	FieldModulus *big.Int
	TWISTBX      *big.Int
	TWISTBY      *big.Int
	PTXX         uint
	PTXY         uint
	PTYX         uint
	PTYY         uint
	PTZX         uint
	PTZY         uint
}

// Init Initialized the required curve
func Init() *BN256G2CURVE {
	FieldModulus, TWISTBX, TWISTBY := big.NewInt(0), big.NewInt(0), big.NewInt(0)
	FieldModulus.SetString("30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47", 16)
	TWISTBX.SetString("2b149d40ceb8aaae81be18991be06ac3b5b4c5e559dbefa33267e6dc24a138e5", 16)
	TWISTBY.SetString("9713b03af0fed4cd2cafadeed8fdf4a74fa084e52d1852e4a2bd0685c315d2", 16)

	return &BN256G2CURVE{
		FieldModulus: FieldModulus,
		TWISTBX:      TWISTBX,
		TWISTBY:      TWISTBY,
		PTXX:         0,
		PTXY:         1,
		PTYX:         2,
		PTYY:         3,
		PTZX:         4,
		PTZY:         5,
	}
}

// ECTwistAdd Addition of twisted points in EC
func (crv *BN256G2CURVE) ECTwistAdd(pt1xx *big.Int, pt1xy *big.Int,
	pt1yx *big.Int, pt1yy *big.Int,
	pt2xx *big.Int, pt2xy *big.Int,
	pt2yx *big.Int, pt2yy *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	if pt1xx == zero && pt1xy == zero && pt2yx == zero && pt2yy == zero {
		if !(pt2xx == zero && pt2xy == zero && pt2yx == zero && pt2yy == zero) {
			// Assert if point in curve
			if !crv.isOnCurve(pt2xx, pt2xy, pt2yx, pt2yy) {
				return nil, nil, nil, nil, fmt.Errorf("The point is not in the curve")
			}
		}
		return pt2xx, pt2xy, pt2yx, pt2yy, nil

	} else if pt2xx == zero && pt2xy == zero && pt2yx == zero && pt2yy == zero {
		// Assert if point in curve
		if !crv.isOnCurve(pt1xx, pt1xy, pt1yx, pt1yy) {
			return nil, nil, nil, nil, fmt.Errorf("The point is not in the curve")
		}
		return pt1xx, pt1xy, pt1yx, pt1yy, nil
	}

	if !crv.isOnCurve(pt2xx, pt2xy, pt2yx, pt2yy) {
		return nil, nil, nil, nil, fmt.Errorf("The point is not in the curve")
	}

	if !crv.isOnCurve(pt1xx, pt1xy, pt1yx, pt1yy) {
		return nil, nil, nil, nil, fmt.Errorf("The point is not in the curve")
	}

	pt3 := crv.ecTwistAddJacobian(pt1xx, pt1xy, pt1yx, pt1yy, big.NewInt(1), big.NewInt(0), pt2xx, pt2xy, pt2yx, pt2yy, one, zero)

	j1, j2, j3, j4 := crv.fromJacobian(pt3[crv.PTXX], pt3[crv.PTXY], pt3[crv.PTYX], pt3[crv.PTYY], pt3[crv.PTZX], pt3[crv.PTZY])
	return j1, j2, j3, j4, nil

}

// isOnCurve verifies if points in the curve
func (crv *BN256G2CURVE) isOnCurve(xx *big.Int, xy *big.Int, yy *big.Int, yx *big.Int) bool {
	var yyx, yyy, xxxx, xxxy *big.Int
	yyx, yyy = crv.fq2mul(yx, yy, yx, yy)
	xxxx, xxxy = crv.fq2mul(xx, xy, xx, xy)
	xxxx, xxxy = crv.fq2mul(xxxx, xxxy, xx, xy)
	yyx, yyy = crv.fq2sub(yyx, yyy, xxxx, xxxy)
	yyx, yyy = crv.fq2sub(yyx, yyy, crv.TWISTBX, crv.TWISTBY)

	return yyx == big.NewInt(0) && yyy == big.NewInt(0)
}

func (crv *BN256G2CURVE) fq2mul(xx *big.Int, xy *big.Int, yx *big.Int, yy *big.Int) (out1 *big.Int, out2 *big.Int) {
	out1 = submod(mulmod(xx, yx, crv.FieldModulus), mulmod(xy, yy, crv.FieldModulus), crv.FieldModulus)
	out2 = submod(mulmod(xx, yy, crv.FieldModulus), mulmod(xy, yx, crv.FieldModulus), crv.FieldModulus)
	return
}

func (crv *BN256G2CURVE) fq2sub(xx *big.Int, xy *big.Int, yx *big.Int, yy *big.Int) (rx *big.Int, ry *big.Int) {
	ry = submod(xx, yx, crv.FieldModulus)
	rx = submod(xy, yy, crv.FieldModulus)
	return
}

func submod(a *big.Int, b *big.Int, n *big.Int) *big.Int {
	return addmod(a, big.NewInt(0).Sub(n, b), n)
}

func addmod(x *big.Int, y *big.Int, k *big.Int) (out *big.Int) {
	out.Add(x, y)
	return big.NewInt(0).Mod(out, k)
}

func mulmod(x *big.Int, y *big.Int, k *big.Int) (out *big.Int) {
	out.Mul(x, y)
	return big.NewInt(0).Mod(out, k)
}

func (crv *BN256G2CURVE) fq2muc(xx *big.Int, xy *big.Int, c *big.Int) (*big.Int, *big.Int) {
	return mulmod(xx, c, crv.FieldModulus), mulmod(xy, c, crv.FieldModulus)
}

func (crv *BN256G2CURVE) ecTwistDoubleJacobian(pt1xx *big.Int, pt1xy *big.Int,
	pt1yx *big.Int, pt1yy *big.Int,
	pt1zx *big.Int, pt1zy *big.Int) (pt2xx, pt2xy,
	pt2yx *big.Int, pt2yy *big.Int,
	pt2zx *big.Int, pt2zy *big.Int) {
	pt2xx, pt2xy = crv.fq2muc(pt1xx, pt1xy, big.NewInt(3)) // 3 * x
	pt2xx, pt2xy = crv.fq2mul(pt2xx, pt2xy, pt1xx, pt1xy)  // W = 3 * x * x
	pt1zx, pt1zy = crv.fq2mul(pt1yx, pt1yy, pt1zx, pt1zy)  // S = y * z
	pt2yx, pt2yy = crv.fq2mul(pt1xx, pt1xy, pt1yx, pt1yy)  // x * y
	pt2yx, pt2yy = crv.fq2mul(pt2yx, pt2yy, pt1zx, pt1zy)  // B = x * y * S
	pt1xx, pt1xy = crv.fq2mul(pt2xx, pt2xy, pt2xx, pt2xy)  // W * W
	pt2zx, pt2zy = crv.fq2muc(pt2yx, pt2yy, big.NewInt(8)) // 8 * B
	pt1xx, pt1xy = crv.fq2sub(pt1xx, pt1xy, pt2zx, pt2zy)  // H = W * W - 8 * B
	pt2zx, pt2zy = crv.fq2mul(pt1zx, pt1zy, pt1zx, pt1zy)  // S_squared = S * S
	pt2yx, pt2yy = crv.fq2muc(pt2yx, pt2yy, big.NewInt(4)) // 4 * B
	pt2yx, pt2yy = crv.fq2sub(pt2yx, pt2yy, pt1xx, pt1xy)  // 4 * B - H
	pt2yx, pt2yy = crv.fq2mul(pt2yx, pt2yy, pt2xx, pt2xy)  // W * (4 * B - H)
	pt2xx, pt2xy = crv.fq2muc(pt1yx, pt1yy, big.NewInt(8)) // 8 * y
	pt2xx, pt2xy = crv.fq2mul(pt2xx, pt2xy, pt1yx, pt1yy)  // 8 * y * y
	pt2xx, pt2xy = crv.fq2mul(pt2xx, pt2xy, pt2zx, pt2zy)  // 8 * y * y * S_squared
	pt2yx, pt2yy = crv.fq2sub(pt2yx, pt2yy, pt2xx, pt2xy)  // newy = W * (4 * B - H) - 8 * y * y * S_squared
	pt2xx, pt2xy = crv.fq2muc(pt1xx, pt1xy, big.NewInt(2)) // 2 * H
	pt2xx, pt2xy = crv.fq2mul(pt2xx, pt2xy, pt1zx, pt1zy)  // newx = 2 * H * S
	pt2zx, pt2zy = crv.fq2mul(pt1zx, pt1zy, pt2zx, pt2zy)  // S * S_squared
	pt2zx, pt2zy = crv.fq2muc(pt2zx, pt2zy, big.NewInt(8)) // newz = 8 * S * S_squared

	return
}

func (crv *BN256G2CURVE) ecTwistAddJacobian(pt1xx *big.Int, pt1xy *big.Int,
	pt1yx *big.Int, pt1yy *big.Int,
	pt1zx *big.Int, pt1zy *big.Int,
	pt2xx *big.Int, pt2xy *big.Int,
	pt2yx *big.Int, pt2yy *big.Int,
	pt2zx *big.Int, pt2zy *big.Int) (pt3 [6]*big.Int) {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	if pt1zx == zero && pt1zy == zero {
		pt3[crv.PTXX], pt3[crv.PTXY], pt3[crv.PTYX], pt3[crv.PTYY], pt3[crv.PTZX], pt3[crv.PTZY] = pt2xx, pt2xy, pt2yx, pt2yy, pt2zx, pt2zy
		return
	} else if pt2zx == zero && pt2zy == zero {
		pt3[crv.PTXX], pt3[crv.PTXY], pt3[crv.PTYX], pt3[crv.PTYY], pt3[crv.PTZX], pt3[crv.PTZY] = pt1xx, pt1xy, pt1yx, pt1yy, pt1zx, pt1zy
		return
	}

	pt2yx, pt2yy = crv.fq2mul(pt2yx, pt2yy, pt1zx, pt1zy)                 // U1 = y2 * z1
	pt3[crv.PTYX], pt3[crv.PTYY] = crv.fq2mul(pt1yx, pt1yy, pt2zx, pt2zy) // U2 = y1 * z2
	pt2xx, pt2xy = crv.fq2mul(pt2xx, pt2xy, pt1zx, pt1zy)                 // V1 = x2 * z1
	pt3[crv.PTZX], pt3[crv.PTZY] = crv.fq2mul(pt1xx, pt1xy, pt2zx, pt2zy) // V2 = x1 * z2

	if pt2xx == pt3[crv.PTZX] && pt2xy == pt3[crv.PTZY] {
		if pt2yx == pt3[crv.PTYX] && pt2yy == pt3[crv.PTYY] {
			pt3[crv.PTXX], pt3[crv.PTXY], pt3[crv.PTYX], pt3[crv.PTYY], pt3[crv.PTZX], pt3[crv.PTZY] = crv.ecTwistDoubleJacobian(pt1xx, pt1xy, pt1yx, pt1yy, pt1zx, pt1zy)
			return
		}
		pt3[crv.PTXX], pt3[crv.PTXY], pt3[crv.PTYX], pt3[crv.PTYY], pt3[crv.PTZX], pt3[crv.PTZY] = one, zero, one, zero, zero, zero
		return
	}

	pt2zx, pt2zy = crv.fq2mul(pt1zx, pt1zy, pt2zx, pt2zy)                 // W = z1 * z2
	pt1xx, pt1xy = crv.fq2sub(pt2yx, pt2yy, pt3[crv.PTYX], pt3[crv.PTYY]) // U = U1 - U2
	pt1yx, pt1yy = crv.fq2sub(pt2xx, pt2xy, pt3[crv.PTZX], pt3[crv.PTZY]) // V = V1 - V2
	pt1zx, pt1zy = crv.fq2mul(pt1yx, pt1yy, pt1yx, pt1yy)                 // V_squared = V * V
	pt2yx, pt2yy = crv.fq2mul(pt1zx, pt1zy, pt3[crv.PTZX], pt3[crv.PTZY]) // V_squared_times_V2 = V_squared * V2
	pt1zx, pt1zy = crv.fq2mul(pt1zx, pt1zy, pt1yx, pt1yy)                 // V_cubed = V * V_squared
	pt3[crv.PTZX], pt3[crv.PTZY] = crv.fq2mul(pt1zx, pt1zy, pt2zx, pt2zy) // newz = V_cubed * W
	pt2xx, pt2xy = crv.fq2mul(pt1xx, pt1xy, pt1xx, pt1xy)                 // U * U
	pt2xx, pt2xy = crv.fq2mul(pt2xx, pt2xy, pt2zx, pt2zy)                 // U * U * W
	pt2xx, pt2xy = crv.fq2sub(pt2xx, pt2xy, pt1zx, pt1zy)                 // U * U * W - V_cubed
	pt2zx, pt2zy = crv.fq2muc(pt2yx, pt2yy, big.NewInt(2))                // 2 * V_squared_times_V2
	pt2xx, pt2xy = crv.fq2sub(pt2xx, pt2xy, pt2zx, pt2zy)                 // A = U * U * W - V_cubed - 2 * V_squared_times_V2
	pt3[crv.PTXX], pt3[crv.PTXY] = crv.fq2mul(pt1yx, pt1yy, pt2xx, pt2xy) // newx = V * A
	pt1yx, pt1yy = crv.fq2sub(pt2yx, pt2yy, pt2xx, pt2xy)                 // V_squared_times_V2 - A
	pt1yx, pt1yy = crv.fq2mul(pt1xx, pt1xy, pt1yx, pt1yy)                 // U * (V_squared_times_V2 - A)
	pt1xx, pt1xy = crv.fq2mul(pt1zx, pt1zy, pt3[crv.PTYX], pt3[crv.PTYY]) // V_cubed * U2
	pt3[crv.PTYX], pt3[crv.PTYY] = crv.fq2sub(pt1yx, pt1yy, pt1xx, pt1xy) // newy = U * (V_squared_times_V2 - A) - V_cubed * U2
	return
}

func (crv *BN256G2CURVE) fq2inv(x *big.Int, y *big.Int) (*big.Int, *big.Int) {
	inv := modInv(addmod(mulmod(y, y, crv.FieldModulus), mulmod(x, x, crv.FieldModulus), crv.FieldModulus), crv.FieldModulus)
	return mulmod(x, inv, crv.FieldModulus), big.NewInt(0).Sub(crv.FieldModulus, mulmod(y, inv, crv.FieldModulus))
}

func modInv(a *big.Int, n *big.Int) (t *big.Int) {
	t = big.NewInt(0)
	newT := big.NewInt(1)
	r := n
	newR := a
	q := big.NewInt(0)
	for newR != big.NewInt(0) {
		q.Div(r, newR)
		t, newT = newT, submod(t, mulmod(q, newT, n), n)
		tmp1 := big.NewInt(0).Mul(q, newR)
		r, newR = newR, big.NewInt(0).Sub(r, tmp1)
	}
	return
}

func (crv *BN256G2CURVE) fromJacobian(
	pt1xx *big.Int, pt1xy *big.Int,
	pt1yx *big.Int, pt1yy *big.Int,
	pt1zx *big.Int, pt1zy *big.Int) (pt2xx *big.Int, pt2xy *big.Int,
	pt2yx *big.Int, pt2yy *big.Int) {

	invzx, invzy := big.NewInt(0), big.NewInt(0)
	invzx, invzy = crv.fq2inv(pt1zx, pt1zy)
	pt2xx, pt2xy = crv.fq2mul(pt1xx, pt1xy, invzx, invzy)
	pt2yx, pt2yy = crv.fq2mul(pt1yx, pt1yy, invzx, invzy)
	return
}
