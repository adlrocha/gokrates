package lib

import (
	"fmt"
	"os/exec"

	"github.com/adlrocha/gokrates/utils/docker"
)

// GenerateProof generates the witness proof
func GenerateProof(witnessName string) error {
	fmt.Println("[*] Building proof image ...")
	// Preparing files
	cmd := fmt.Sprintf("cp %v/%v %v/witness", zkMaterialPath, witnessName+".witness", zkMaterialPath)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println(err)
	}
	// Building the image
	_, err = docker.BuildImage("gokrates-proof", dockerfilesPath+"witness")
	if err != nil {
		return err
	}

	// Creating witness
	fmt.Println("[*] Generating proof for witness " + witnessName)
	compilerCode := fmt.Sprintf("./zokrates generate-proof")
	_, err = docker.RunContainer("gokrates-proof", "gokrates-proof", compilerCode)
	if err != nil {
		fmt.Println("[!] Error generating proof")
		return err
	}

	fmt.Println("[*] Storing generated proof")
	_, err = docker.StoreFiles("gokrates-proof", "/home/zokrates/proof.json", zkMaterialPath)

	fmt.Println("[*] Removing intermediate images...")
	_, err = docker.RemoveContainer("gokrates-proof")
	_, err = docker.RemoveImage("gokrates-proof")

	return nil
}
