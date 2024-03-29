package main

import (
	"fmt"

	"github.com/adlrocha/gokrates/lib"
)

const zkMaterialPath = "./zk-material/"

func main() {
	// if os.Args[1] ...
	// err := lib.CompileZkCode("sample.code")
	// if err != nil {
	// 	_, err = docker.CleanImages()
	// 	_, err = docker.CleanContainers()
	// 	fmt.Println(err)
	// }

	// err := lib.Setup("sample")
	// if err != nil {
	// 	_, err = docker.CleanImages()
	// 	_, err = docker.CleanContainers()
	// 	fmt.Println(err)
	// }

	// err = lib.CreateWitness("sample", []string{"337", "113569"})
	// if err != nil {
	// 	_, err = docker.CleanImages()
	// 	_, err = docker.CleanContainers()
	// 	fmt.Println(err)
	// }

	// err = lib.GenerateProof("sample")
	// if err != nil {
	// 	_, err = docker.CleanImages()
	// 	_, err = docker.CleanContainers()
	// 	fmt.Println(err)
	// }

	// err := lib.GenerateVk("sample")
	// if err != nil {
	// 	_, err = docker.CleanImages()
	// 	_, err = docker.CleanContainers()
	// 	fmt.Println(err)
	// }

	// fmt.Println(lib.GetVerifyingKey("sample"))
	a := lib.GetProof("sample")
	fmt.Println(a.Input)
}
