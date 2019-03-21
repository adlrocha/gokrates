package lib

import (
	"fmt"

	"github.com/adlrocha/gokrates/utils/docker"
)

const dockerfilesPath = "./utils/dockerfiles/"
const zkMaterialPath = "./zk-material/"

// CompileZkCode flattens computational code.
func compileZkCode(codeFile string) error {
	fmt.Println("[*] Building image for compiler...")
	_, err := docker.BuildImage("gokrates-compiler", dockerfilesPath+"compiler")
	if err != nil {
		return err
	}

	fmt.Println("[*] Compiling program...")
	compilerCode := fmt.Sprintf("./zokrates compile -i ./code/%v.code", codeFile)
	_, err = docker.RunContainer("gokrates-compiler", "gokrates-compiler", compilerCode)
	if err != nil {
		fmt.Println("[!] Error compiling program")
		return err
	}

	fmt.Println("[*] Commiting for compiled container for setup...")
	_, err = docker.CommitContainer("gokrates-compiler", "gokrates-setup")
	if err != nil {
		fmt.Println("[!] Error commiting setup image")
		return err
	}
	fmt.Println("[*] Removing intermediate images...")
	_, err = docker.RemoveContainer("gokrates-compiler")
	_, err = docker.RemoveImage("gokrates-compiler")

	return nil
}

// Setup phase for ZKSnarks proof and verifier keys.
func setupKeys() error {
	fmt.Println("[*] Setting up ZK keys for computation...")
	_, err := docker.RunContainer("gokrates-setup", "gokrates-setup", "./zokrates setup")
	if err != nil {
		fmt.Println("[!] Error setting up ZK keys")
		return err
	}

	fmt.Println("[*] Storing generated keys and crypto material")
	_, err = docker.StoreFiles("gokrates-setup", "/home/zokrates/out", zkMaterialPath)
	_, err = docker.StoreFiles("gokrates-setup", "/home/zokrates/proof.json", zkMaterialPath)
	_, err = docker.StoreFiles("gokrates-setup", "/home/zokrates/proving.key", zkMaterialPath)
	_, err = docker.StoreFiles("gokrates-setup", "/home/zokrates/variables.inf", zkMaterialPath)
	_, err = docker.StoreFiles("gokrates-setup", "/home/zokrates/verification.key", zkMaterialPath)
	if err != nil {
		fmt.Println("[!] Error storing crypto material")
		return err
	}

	fmt.Println("[*] Removing all gokrates images...")
	_, err = docker.RemoveContainer("gokrates-setup")
	_, err = docker.RemoveImage("gokrates-setup")
	// _, err = docker.CleanContainers()
	// _, err = docker.CleanImages()

	return nil
}

// Setup Phase for ZKSnarks
func Setup(codeFile string) error {
	// Compile selected code for flattening.
	err := compileZkCode(codeFile)
	if err != nil {
		return err
	}
	// Generate cryptomaterial
	err = setupKeys()
	if err != nil {
		return err
	}
	return nil
}
