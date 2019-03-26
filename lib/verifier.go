package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/adlrocha/gokrates/utils/docker"
	"github.com/adlrocha/gokrates/utils/pairings"
)

// const zkMaterialPath = "./zk-material/"

// VerifyingKey ZK verifier structure
type VerifyingKey struct {
	A          pairings.G2Point
	B          pairings.G1Point
	C          pairings.G2Point
	gamma      pairings.G2Point
	gammaBeta1 pairings.G1Point
	gammaBeta2 pairings.G2Point
	Z          pairings.G2Point
	IC         [3]pairings.G1Point
}

// Proof ZK prover structure
type Proof struct {
	A     pairings.G1Point
	Ap    pairings.G1Point
	B     pairings.G2Point
	Bp    pairings.G1Point
	C     pairings.G1Point
	Cp    pairings.G1Point
	K     pairings.G1Point
	H     pairings.G1Point
	Input [2]*big.Int
}

// ProofString auxiliary struct to get Proof points
type ProofString struct {
	A  []string   `json:"A"`
	Ap []string   `json:"A_p"`
	B  [][]string `json:"B"`
	Bp []string   `json:"B_p"`
	C  []string   `json:"C"`
	Cp []string   `json:"C_p"`
	K  []string   `json:"K"`
	H  []string   `json:"H"`
}

// GenerateVk generate verifying keys
func GenerateVk(verifierName string) error {
	fmt.Println("[*] Building verifying image ...")
	_, err := docker.BuildImage("gokrates-verify", dockerfilesPath+"witness")
	if err != nil {
		return err
	}

	// Create verifying keys
	fmt.Println("[*] Creating verifying key...")
	compilerCode := fmt.Sprintf("./zokrates export-verifier")
	_, err = docker.RunContainer("gokrates-verify", "gokrates-verify", compilerCode)
	if err != nil {
		fmt.Println("[!] Error creating verifying key")
		return err
	}

	fmt.Println("[*] Storing generated witness")
	_, err = docker.StoreFiles("gokrates-verify", "/home/zokrates/verifier.sol", zkMaterialPath+verifierName+".verifier")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("[*] Removing intermediate images...")
	_, err = docker.RemoveContainer("gokrates-verify")
	_, err = docker.RemoveImage("gokrates-verify")

	return nil
}

// GetVerifyingKey gets verifying key from zk-material
func GetVerifyingKey(verifierName string) (vk *VerifyingKey) {
	data := ReadFile(zkMaterialPath + verifierName + ".verifier")

	return &VerifyingKey{
		A:          verifierG2Point("vk.A", data),
		B:          verifierG1Point("vk.B", data),
		C:          verifierG2Point("vk.C", data),
		gamma:      verifierG2Point("vk.gamma", data),
		gammaBeta1: verifierG1Point("vk.gammaBeta1", data),
		gammaBeta2: verifierG2Point("vk.gammaBeta2", data),
		Z:          verifierG2Point("vk.Z", data),
		IC: [3]pairings.G1Point{verifierG1Point(`vk.IC\[0\]`, data),
			verifierG1Point(`vk.IC\[1\]`, data),
			verifierG1Point(`vk.IC\[2\]`, data)},
	}
}

// GetProof gets proof from zk-material
func GetProof(proofName string) (proof *Proof) {
	// Get the data
	var proofString map[string]ProofString
	data := ReadFile(zkMaterialPath + proofName + ".proof")
	_ = json.Unmarshal([]byte(data), &proofString)
	proofPoints := proofString["proof"]

	var prInput map[string][]int64
	_ = json.Unmarshal([]byte(data), &prInput)
	proofInput := prInput["input"]

	// Create GPoints
	a0, a1 := big.NewInt(0), big.NewInt(0)
	a0.SetString(strings.Replace(proofPoints.A[0], "0x", "", -1), 16)
	a1.SetString(strings.Replace(proofPoints.A[1], "0x", "", -1), 16)
	ap0, ap1 := big.NewInt(0), big.NewInt(0)
	ap0.SetString(strings.Replace(proofPoints.Ap[0], "0x", "", -1), 16)
	ap1.SetString(strings.Replace(proofPoints.Ap[1], "0x", "", -1), 16)
	b00, b01, b10, b11 := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	b00.SetString(strings.Replace(proofPoints.B[0][0], "0x", "", -1), 16)
	b01.SetString(strings.Replace(proofPoints.B[0][1], "0x", "", -1), 16)
	b10.SetString(strings.Replace(proofPoints.B[1][0], "0x", "", -1), 16)
	b11.SetString(strings.Replace(proofPoints.B[1][1], "0x", "", -1), 16)
	bp0, bp1 := big.NewInt(0), big.NewInt(0)
	bp0.SetString(strings.Replace(proofPoints.Bp[0], "0x", "", -1), 16)
	bp1.SetString(strings.Replace(proofPoints.Bp[1], "0x", "", -1), 16)
	c0, c1 := big.NewInt(0), big.NewInt(0)
	c0.SetString(strings.Replace(proofPoints.C[0], "0x", "", -1), 16)
	c1.SetString(strings.Replace(proofPoints.C[1], "0x", "", -1), 16)
	cp0, cp1 := big.NewInt(0), big.NewInt(0)
	cp0.SetString(strings.Replace(proofPoints.Cp[0], "0x", "", -1), 16)
	cp1.SetString(strings.Replace(proofPoints.Cp[1], "0x", "", -1), 16)
	k0, k1 := big.NewInt(0), big.NewInt(0)
	k0.SetString(strings.Replace(proofPoints.K[0], "0x", "", -1), 16)
	k1.SetString(strings.Replace(proofPoints.K[1], "0x", "", -1), 16)
	h0, h1 := big.NewInt(0), big.NewInt(0)
	h0.SetString(strings.Replace(proofPoints.H[0], "0x", "", -1), 16)
	h1.SetString(strings.Replace(proofPoints.H[1], "0x", "", -1), 16)

	return &Proof{
		A:  pairings.G1Point{X: a0, Y: a1},
		Ap: pairings.G1Point{X: ap0, Y: ap1},
		B: pairings.G2Point{
			X: [2]*big.Int{b00, b01},
			Y: [2]*big.Int{b10, b11}},
		Bp:    pairings.G1Point{X: bp0, Y: bp1},
		C:     pairings.G1Point{X: c0, Y: c1},
		Cp:    pairings.G1Point{X: cp0, Y: cp1},
		K:     pairings.G1Point{X: k0, Y: k1},
		H:     pairings.G1Point{X: h0, Y: h1},
		Input: [2]*big.Int{big.NewInt(proofInput[0]), big.NewInt(proofInput[1])},
	}
}

// Verify Zk Proof
func Verify(proof *Proof, vk *VerifyingKey) (uint, error) {
	// Verify lengths
	if len(proof.Input)+1 != len(vk.IC) {
		return 0, errors.New("Proof.Input + 1 and vk.IC lengths not equal")
	}
	// Compute linear combination of vk_x
	vkx := pairings.G1Point{X: big.NewInt(0), Y: big.NewInt(0)}
	for i := 0; i < len(proof.Input); i++ {
		vkx = pairings.AdditionG1(vkx, pairings.ScalarMul(vk.IC[i+1], proof.Input[i]))
	}
	vkx = pairings.AdditionG1(vkx, vk.IC[0])
	if()
	return 0, nil
}

/*

    function verify(uint[] input, Proof proof) internal returns (uint) {
        VerifyingKey memory vk = verifyingKey();
        require(input.length + 1 == vk.IC.length);
        // Compute the linear combination vk_x
        Pairing.G1Point memory vk_x = Pairing.G1Point(0, 0);
        for (uint i = 0; i < input.length; i++)
            vk_x = Pairing.addition(vk_x, Pairing.scalar_mul(vk.IC[i + 1], input[i]));
        vk_x = Pairing.addition(vk_x, vk.IC[0]);
        if (!Pairing.pairingProd2(proof.A, vk.A, Pairing.negate(proof.A_p), Pairing.P2())) return 1;
        if (!Pairing.pairingProd2(vk.B, proof.B, Pairing.negate(proof.B_p), Pairing.P2())) return 2;
        if (!Pairing.pairingProd2(proof.C, vk.C, Pairing.negate(proof.C_p), Pairing.P2())) return 3;
        if (!Pairing.pairingProd3(
            proof.K, vk.gamma,
            Pairing.negate(Pairing.addition(vk_x, Pairing.addition(proof.A, proof.C))), vk.gammaBeta2,
            Pairing.negate(vk.gammaBeta1), proof.B
        )) return 4;
        if (!Pairing.pairingProd3(
                Pairing.addition(vk_x, proof.A), proof.B,
                Pairing.negate(proof.H), vk.Z,
                Pairing.negate(proof.C), Pairing.P2()
        )) return 5;
        return 0;
    }
    event Verified(string s);
    function verifyTx(
            uint[2] a,
            uint[2] a_p,
            uint[2][2] b,
            uint[2] b_p,
            uint[2] c,
            uint[2] c_p,
            uint[2] h,
            uint[2] k,
            uint[2] input
        ) public returns (bool r) {
        Proof memory proof;
        proof.A = Pairing.G1Point(a[0], a[1]);
        proof.A_p = Pairing.G1Point(a_p[0], a_p[1]);
        proof.B = Pairing.G2Point([b[0][0], b[0][1]], [b[1][0], b[1][1]]);
        proof.B_p = Pairing.G1Point(b_p[0], b_p[1]);
        proof.C = Pairing.G1Point(c[0], c[1]);
        proof.C_p = Pairing.G1Point(c_p[0], c_p[1]);
        proof.H = Pairing.G1Point(h[0], h[1]);
        proof.K = Pairing.G1Point(k[0], k[1]);
        uint[] memory inputValues = new uint[](input.length);
        for(uint i = 0; i < input.length; i++){
            inputValues[i] = input[i];
        }
        if (verify(inputValues, proof) == 0) {
            emit Verified("Transaction successfully verified.");
            return true;
        } else {
            return false;
        }
    }
}
*/
