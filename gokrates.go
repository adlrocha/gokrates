package main

import (
	"fmt"

	"github.com/adlrocha/gokrates/lib"
	"github.com/adlrocha/gokrates/utils/docker"
)

func main() {
	// if os.Args[1] ...
	err := lib.CompileZkCode("sample.code")
	if err != nil {
		_, err = docker.CleanImages()
		_, err = docker.CleanContainers()
		fmt.Println(err)
	}

	err = lib.Setup()
	if err != nil {
		_, err = docker.CleanImages()
		_, err = docker.CleanContainers()
		fmt.Println(err)
	}

	err = lib.CreateWitness([]string{"337", "113569"})
	if err != nil {
		_, err = docker.CleanImages()
		_, err = docker.CleanContainers()
		fmt.Println(err)
	}

	err = lib.GenerateProof("sample")
	if err != nil {
		_, err = docker.CleanImages()
		_, err = docker.CleanContainers()
		fmt.Println(err)
	}

}
