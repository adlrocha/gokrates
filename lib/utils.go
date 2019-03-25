package lib

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"regexp"
	"strings"

	"github.com/adlrocha/gokrates/utils/pairings"
)

// ReadFile reads any text file
func ReadFile(file string) string {
	b, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	return string(b)
}

// verifierG2Point extracts a G2 point from verifying key
func verifierG2Point(name string, data string) pairings.G2Point {
	// Get G2Point line
	regExp := fmt.Sprintf(`(%v = Pairing.G2Point\(\[)(\w+)(\, )(\w+)(], \[)(\w+)(\, )(\w+)(\]\)\;)`, name)
	r, _ := regexp.Compile(regExp)
	b := r.FindAllString(data, -1)[0]
	regExp = fmt.Sprintf(`%v = Pairing.G2Point`, name)
	b = strings.Replace(b, regExp, "", -1)
	// Split different parts of the string
	s := strings.Split(b, ",")

	// Get coordinates by string
	sb0 := strings.Replace(removeCharacters(s[0]), "0x", "", -1)
	sb1 := strings.Replace(removeCharacters(s[1]), "0x", "", -1)
	sb2 := strings.Replace(removeCharacters(s[2]), "0x", "", -1)
	sb3 := strings.Replace(removeCharacters(s[3]), "0x", "", -1)

	// Create G2Point
	b0, b1, b2, b3 := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)
	b0.SetString(sb0, 16)
	b1.SetString(sb1, 16)
	b2.SetString(sb2, 16)
	b3.SetString(sb3, 16)

	return pairings.G2Point{X: [2]*big.Int{b0, b1}, Y: [2]*big.Int{b2, b3}}
}

// verifierG1Point extracts a G1 point from verifying key
func verifierG1Point(name string, data string) pairings.G1Point {
	// Get G2Point line
	regExp := fmt.Sprintf(`(%v = Pairing.G1Point\()(\w+)(\,\ )(\w+)(\)\;)`, name)
	r, _ := regexp.Compile(regExp)
	b := r.FindAllString(data, -1)[0]
	regExp = fmt.Sprintf(`%v = Pairing.G1Point`, name)
	b = strings.Replace(b, regExp, "", -1)
	// Split different parts of the string
	s := strings.Split(b, ",")

	// Get coordinates by string
	sbtmp0 := strings.Replace(removeCharacters(s[0]), removeCharacters(regExp), "", -1)
	sb0 := strings.Replace(sbtmp0, "0x", "", -1)
	sb1 := strings.Replace(removeCharacters(s[1]), "0x", "", -1)

	// Create G2Point
	b0, b1 := big.NewInt(0), big.NewInt(0)
	b0.SetString(sb0, 16)
	b1.SetString(sb1, 16)

	return pairings.G1Point{X: b0, Y: b1}
}

// Removes non-alphanumeric characters
func removeCharacters(a string) string {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		fmt.Println("Error")
	}
	processedString := reg.ReplaceAllString(a, "")

	return processedString
}
