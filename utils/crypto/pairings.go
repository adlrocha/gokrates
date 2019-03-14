package crypto

// This file is MIT Licensed.
//
// Copyright 2017 Christian Reitwiessner
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// pragma solidity ^0.4.14;
// library Pairing {
//     struct G1Point {
//         uint X;
//         uint Y;
//     }
//     // Encoding of field elements is: X[0] * z + X[1]
//     struct G2Point {
//         uint[2] X;
//         uint[2] Y;
//     }
//     /// @return the generator of G1
//     function P1() pure internal returns (G1Point) {
//         return G1Point(1, 2);
//     }
//     /// @return the generator of G2
//     function P2() pure internal returns (G2Point) {
//         return G2Point(
//             [11559732032986387107991004021392285783925812861821192530917403151452391805634,
//              10857046999023057135944570762232829481370756359578518086990519993285655852781],
//             [4082367875863433681332203403145435568316851327593401208105741076214120093531,
//              8495653923123431417604973247489272438418190587263600148770280649306958101930]
//         );
//     }
//     /// @return the negation of p, i.e. p.addition(p.negate()) should be zero.
//     function negate(G1Point p) pure internal returns (G1Point) {
//         // The prime q in the base field F_q for G1
//         uint q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;
//         if (p.X == 0 && p.Y == 0)
//             return G1Point(0, 0);
//         return G1Point(p.X, q - (p.Y % q));
//     }
//     /// @return the sum of two points of G1
//     function addition(G1Point p1, G1Point p2) internal returns (G1Point r) {
//         uint[4] memory input;
//         input[0] = p1.X;
//         input[1] = p1.Y;
//         input[2] = p2.X;
//         input[3] = p2.Y;
//         bool success;
//         assembly {
//             success := call(sub(gas, 2000), 6, 0, input, 0xc0, r, 0x60)
//             // Use "invalid" to make gas estimation work
//             switch success case 0 { invalid() }
//         }
//         require(success);
//     }
//     /// @return the sum of two points of G2
//     function addition(G2Point p1, G2Point p2) internal pure returns (G2Point r) {
//         (r.X[1], r.X[0], r.Y[1], r.Y[0]) = BN256G2.ECTwistAdd(p1.X[1],p1.X[0],p1.Y[1],p1.Y[0],p2.X[1],p2.X[0],p2.Y[1],p2.Y[0]);
//     }
//     /// @return the product of a point on G1 and a scalar, i.e.
//     /// p == p.scalar_mul(1) and p.addition(p) == p.scalar_mul(2) for all points p.
//     function scalar_mul(G1Point p, uint s) internal returns (G1Point r) {
//         uint[3] memory input;
//         input[0] = p.X;
//         input[1] = p.Y;
//         input[2] = s;
//         bool success;
//         assembly {
//             success := call(sub(gas, 2000), 7, 0, input, 0x80, r, 0x60)
//             // Use "invalid" to make gas estimation work
//             switch success case 0 { invalid() }
//         }
//         require (success);
//     }
//     /// @return the result of computing the pairing check
//     /// e(p1[0], p2[0]) *  .... * e(p1[n], p2[n]) == 1
//     /// For example pairing([P1(), P1().negate()], [P2(), P2()]) should
//     /// return true.
//     function pairing(G1Point[] p1, G2Point[] p2) internal returns (bool) {
//         require(p1.length == p2.length);
//         uint elements = p1.length;
//         uint inputSize = elements * 6;
//         uint[] memory input = new uint[](inputSize);
//         for (uint i = 0; i < elements; i++)
//         {
//             input[i * 6 + 0] = p1[i].X;
//             input[i * 6 + 1] = p1[i].Y;
//             input[i * 6 + 2] = p2[i].X[0];
//             input[i * 6 + 3] = p2[i].X[1];
//             input[i * 6 + 4] = p2[i].Y[0];
//             input[i * 6 + 5] = p2[i].Y[1];
//         }
//         uint[1] memory out;
//         bool success;
//         assembly {
//             success := call(sub(gas, 2000), 8, 0, add(input, 0x20), mul(inputSize, 0x20), out, 0x20)
//             // Use "invalid" to make gas estimation work
//             switch success case 0 { invalid() }
//         }
//         require(success);
//         return out[0] != 0;
//     }
//     /// Convenience method for a pairing check for two pairs.
//     function pairingProd2(G1Point a1, G2Point a2, G1Point b1, G2Point b2) internal returns (bool) {
//         G1Point[] memory p1 = new G1Point[](2);
//         G2Point[] memory p2 = new G2Point[](2);
//         p1[0] = a1;
//         p1[1] = b1;
//         p2[0] = a2;
//         p2[1] = b2;
//         return pairing(p1, p2);
//     }
//     /// Convenience method for a pairing check for three pairs.
//     function pairingProd3(
//             G1Point a1, G2Point a2,
//             G1Point b1, G2Point b2,
//             G1Point c1, G2Point c2
//     ) internal returns (bool) {
//         G1Point[] memory p1 = new G1Point[](3);
//         G2Point[] memory p2 = new G2Point[](3);
//         p1[0] = a1;
//         p1[1] = b1;
//         p1[2] = c1;
//         p2[0] = a2;
//         p2[1] = b2;
//         p2[2] = c2;
//         return pairing(p1, p2);
//     }
//     /// Convenience method for a pairing check for four pairs.
//     function pairingProd4(
//             G1Point a1, G2Point a2,
//             G1Point b1, G2Point b2,
//             G1Point c1, G2Point c2,
//             G1Point d1, G2Point d2
//     ) internal returns (bool) {
//         G1Point[] memory p1 = new G1Point[](4);
//         G2Point[] memory p2 = new G2Point[](4);
//         p1[0] = a1;
//         p1[1] = b1;
//         p1[2] = c1;
//         p1[3] = d1;
//         p2[0] = a2;
//         p2[1] = b2;
//         p2[2] = c2;
//         p2[3] = d2;
//         return pairing(p1, p2);
//     }
// }
// contract Verifier {
//     using Pairing for *;
//     struct VerifyingKey {
//         Pairing.G2Point A;
//         Pairing.G1Point B;
//         Pairing.G2Point C;
//         Pairing.G2Point gamma;
//         Pairing.G1Point gammaBeta1;
//         Pairing.G2Point gammaBeta2;
//         Pairing.G2Point Z;
//         Pairing.G1Point[] IC;
//     }
//     struct Proof {
//         Pairing.G1Point A;
//         Pairing.G1Point A_p;
//         Pairing.G2Point B;
//         Pairing.G1Point B_p;
//         Pairing.G1Point C;
//         Pairing.G1Point C_p;
//         Pairing.G1Point K;
//         Pairing.G1Point H;
//     }
//     function verifyingKey() pure internal returns (VerifyingKey vk) {
//         vk.A = Pairing.G2Point([0xb8f19d8bd0c4e99b1a5505a2323b05590862d3703c163de79d888065936a887, 0x3b8e028a169bce1a1a24ba7756e4552aa0d6fa89e2918337bd703d83eaf2b7b], [0x2d5a983cd061ab980675598ea6271b508196b41120babd10ceed70626667315a, 0x2084e16c4664daa72e32cbde78b2edbdceadcb53046ffa690ccb62375d9a2295]);
//         vk.B = Pairing.G1Point(0x1833b9c7d23703424a152393971f9e54e39ff2f45bfb625fcbc6665ee37ef050, 0x2630625696e88227fbd4b6b62c9d84076e0dc7ea260497221848fefbf9a5a0fd);
//         vk.C = Pairing.G2Point([0x5bb8d155d0eb90f18c8b0d4b51dfb2bc28f02590c9bf57675414801a9edf01d, 0x7d19de94769d1ba6e2759db9764fc7d28df9f43d32973dc52de89975ef728ff], [0x2cdc5afdebbdf04d581ed88d3b2df68e90e02e88ae0c4cdccef5e98e21d2efe0, 0x27172b122523320d8500ed4a8c1e206fbe085a08863334d09c1f5f35044e91e8]);
//         vk.gamma = Pairing.G2Point([0xe72e0d5ee5ba6ddeaf8b5b26dbfc3252a4ab38869608a218bb02df20f45d689, 0x103863ad1c684088f364ef0298c34c9592021b8338d76c445b8ce447fb8d970], [0x191de42acc14b028b3a4be55e45d46a1a5d5ae830f03b4c8eff3b8761488dfa6, 0x1af9518a72f4462f3aaf492486193cdb6fbb59422cfe6db6211620952842c40c]);
//         vk.gammaBeta1 = Pairing.G1Point(0xe9e603440d912746765bcca0ceb86a5e0b6f6c57414f96c47bf1c7d12eb2632, 0x1ba635bf4ef02457497830ab2ec75cfe78d2ee2ca07721a49926f1d9ce1db101);
//         vk.gammaBeta2 = Pairing.G2Point([0x25d7e22214429dad86f375be622b0c816296d0de89c70434fb964ac1d147840e, 0x131d4721e408a677a6b6065217151254db81acb1bb2b91c046d282b6f9c3e573], [0xe17f9016c3d43ca2c1a51deeb9748e209371256e7fc99ea689f0542878e44b2, 0xdd005f46c343ef314d7569bc54db078c6935e4e43b8c18f9dd67e30a6a9a6b5]);
//         vk.Z = Pairing.G2Point([0xb9a5be62d9b49b80dd791e466035ad21c18a97855b134f93f334c992cf21a4d, 0x28d101e2f84c701dc9aeabe7c8892c28801055944047e79e2579172a9ffacc9a], [0x81501da2b5f4499aecfa215cb650ddcb001ffe62996ce2158c795ba256e86b, 0x134ebacd0991d4ff16eaffcc7f36234898829b8f12b47cbcde42e3c2ecf3cd42]);
//         vk.IC = new Pairing.G1Point[](3);
//         vk.IC[0] = Pairing.G1Point(0x176e26add417bbe292e2ee31604b7a1e817a8f44d861e49e4fd939f1f89eca53, 0x14fa2f17775dfc1d6addcf0c884fd6396a356a0edc5830d784d8a2f0c3fc77be);
//         vk.IC[1] = Pairing.G1Point(0x2f09faab4ac67dcad75ca222c6340e4f6eecedbeda24ead9244948cfefde7121, 0x23e88a720eb6b4a70e21725d36b8ad837d0cda42d0e5c537d10a45702ef4b4ea);
//         vk.IC[2] = Pairing.G1Point(0x8e5e5da475f6eb40e0e9d02e9805b83870d5949e012dffed269ac02d1132bbb, 0x22fd0dbd086c28b0ded52042cf54a095568bb52e7daca8fc6bb7abc8494e3750);
//     }
//     function verify(uint[] input, Proof proof) internal returns (uint) {
//         VerifyingKey memory vk = verifyingKey();
//         require(input.length + 1 == vk.IC.length);
//         // Compute the linear combination vk_x
//         Pairing.G1Point memory vk_x = Pairing.G1Point(0, 0);
//         for (uint i = 0; i < input.length; i++)
//             vk_x = Pairing.addition(vk_x, Pairing.scalar_mul(vk.IC[i + 1], input[i]));
//         vk_x = Pairing.addition(vk_x, vk.IC[0]);
//         if (!Pairing.pairingProd2(proof.A, vk.A, Pairing.negate(proof.A_p), Pairing.P2())) return 1;
//         if (!Pairing.pairingProd2(vk.B, proof.B, Pairing.negate(proof.B_p), Pairing.P2())) return 2;
//         if (!Pairing.pairingProd2(proof.C, vk.C, Pairing.negate(proof.C_p), Pairing.P2())) return 3;
//         if (!Pairing.pairingProd3(
//             proof.K, vk.gamma,
//             Pairing.negate(Pairing.addition(vk_x, Pairing.addition(proof.A, proof.C))), vk.gammaBeta2,
//             Pairing.negate(vk.gammaBeta1), proof.B
//         )) return 4;
//         if (!Pairing.pairingProd3(
//                 Pairing.addition(vk_x, proof.A), proof.B,
//                 Pairing.negate(proof.H), vk.Z,
//                 Pairing.negate(proof.C), Pairing.P2()
//         )) return 5;
//         return 0;
//     }
//     event Verified(string s);
//     function verifyTx(
//             uint[2] a,
//             uint[2] a_p,
//             uint[2][2] b,
//             uint[2] b_p,
//             uint[2] c,
//             uint[2] c_p,
//             uint[2] h,
//             uint[2] k,
//             uint[2] input
//         ) public returns (bool r) {
//         Proof memory proof;
//         proof.A = Pairing.G1Point(a[0], a[1]);
//         proof.A_p = Pairing.G1Point(a_p[0], a_p[1]);
//         proof.B = Pairing.G2Point([b[0][0], b[0][1]], [b[1][0], b[1][1]]);
//         proof.B_p = Pairing.G1Point(b_p[0], b_p[1]);
//         proof.C = Pairing.G1Point(c[0], c[1]);
//         proof.C_p = Pairing.G1Point(c_p[0], c_p[1]);
//         proof.H = Pairing.G1Point(h[0], h[1]);
//         proof.K = Pairing.G1Point(k[0], k[1]);
//         uint[] memory inputValues = new uint[](input.length);
//         for(uint i = 0; i < input.length; i++){
//             inputValues[i] = input[i];
//         }
//         if (verify(inputValues, proof) == 0) {
//             emit Verified("Transaction successfully verified.");
//             return true;
//         } else {
//             return false;
//         }
//     }
// }
