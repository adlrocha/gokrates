package docker

import (
	"fmt"
	"os/exec"
)

// BuildImage builds docker image
func BuildImage(name string, dockerfile string) (string, error) {
	cmd := fmt.Sprintf("docker image build -t %v . --file %v", name, dockerfile)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(out), nil
}

//RemoveImage removes docker image
func RemoveImage(name string) (string, error) {
	cmd := fmt.Sprintf("docker rmi %v", name)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//RunContainer runs docker container
func RunContainer(name string, image string, entrypoint string) (string, error) {
	cmd := fmt.Sprintf("docker run --name %v %v %v", name, image, entrypoint)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//RemoveContainer removes docker container
func RemoveContainer(name string) (string, error) {
	cmd := fmt.Sprintf("docker rm -f %v", name)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//CommitContainer commits docker container
func CommitContainer(srcContainer string, dstImage string) (string, error) {
	cmd := fmt.Sprintf("docker commit %v %v", srcContainer, dstImage)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//Pwd Auxiliary function for relative paths
func Pwd() (string, error) {
	cmd := fmt.Sprintf("pwd")
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//CleanImages clean Gokrates images
func CleanImages() (string, error) {
	cmd := fmt.Sprintf("docker rmi $(docker images %v | awk '{print $3}' | awk '{if(NR>1)print}')", "gokrates")
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//CleanContainers clean Gokrates containers
func CleanContainers() (string, error) {
	cmd := fmt.Sprintf("docker rm -f $(docker ps -a | grep %v* | awk '{print $1}' | awk '{if(NR>1)print}')", "gokrates")
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// StoreFiles stores the files from a docker container in your local environment
func StoreFiles(container string, source string, destination string) (string, error) {
	cmd := fmt.Sprintf("docker cp %v:%v %v", container, source, destination)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
