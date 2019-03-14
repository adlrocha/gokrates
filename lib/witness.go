package lib

import (
	"fmt"
	"strings"

	"github.com/adlrocha/gokrates/utils/docker"
)

// CreateWitness creates a new witness for certain computation.
func CreateWitness(attributes []string) error {
	fmt.Println("[*] Building witness image ...")
	_, err := docker.BuildImage("gokrates-witness", dockerfilesPath+"witness")
	if err != nil {
		return err
	}

	// Extracting and formatting attributes
	fmt.Println("[*] Creating new witness...")
	attr := fmt.Sprintf("%v", attributes)
	attr = strings.Replace(attr, "[", "", -1)
	attr = strings.Replace(attr, "]", "", -1)

	// Creating witness
	compilerCode := fmt.Sprintf("./zokrates compute-witness -a %v", attr)
	_, err = docker.RunContainer("gokrates-witness", "gokrates-witness", compilerCode)
	if err != nil {
		fmt.Println("[!] Error compiling program")
		return err
	}

	fmt.Println("[*] Storing generated witness")
	_, err = docker.StoreFiles("gokrates-witness", "/home/zokrates/witness", zkMaterialPath)

	fmt.Println("[*] Removing intermediate images...")
	_, err = docker.RemoveContainer("gokrates-witness")
	_, err = docker.RemoveImage("gokrates-witness")

	return nil
}
