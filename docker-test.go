package main

// import (
// 	"fmt"
// 	"os/exec"
// )

// func main() {

// 	// endpoint := "unix:///var/run/docker.sock"
// 	// client, err := docker.NewClient(endpoint)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// // imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// // Configure PortBinding
// 	// portBinding := docker.PortBinding{HostIP: "0.0.0.0", HostPort: "9000"}
// 	// listBind := make(map[docker.Port][]docker.PortBinding, 0)
// 	// thePort := docker.Port("9000")
// 	// listBind[thePort] = append(listBind[docker.Port("tcp/9000")], portBinding)

// 	out, err := exec.Command("sh", "-c", "docker image build -t zokrates-compile . --file ./Dockerfile").Output()
// 	fmt.Println(string(out))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// Build new image
// 	// var b bytes.Buffer
// 	// foo := bufio.NewWriter(&b)
// 	// buildImageOptions := docker.BuildImageOptions{Name: "zokrates-compile",
// 	// 	Dockerfile:   "./Dockerfile",
// 	// 	OutputStream: foo,
// 	// 	Remote:       "unix:///var/run/docker.sock",
// 	// }
// 	// err = client.BuildImage(buildImageOptions)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	// // Configurations
// 	// config := docker.Config{
// 	// 	Image:      "zokrates",
// 	// 	Entrypoint: []string{"sleep 6000"},
// 	// }
// 	// hostConfig := docker.HostConfig{PortBindings: listBind}
// 	// ctx := context.Background()
// 	// fmt.Println("Creating container options")
// 	// createContainerOptions := docker.CreateContainerOptions{Context: ctx, Name: "test", Config: &config, HostConfig: &hostConfig}

// 	// // Create and start the container
// 	// fmt.Println("Creating container")
// 	// container, err := client.CreateContainer(createContainerOptions)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// client.StartContainer(container.ID, &hostConfig)

// 	// out, err := exec.Command("sh", "-c", "docker ps -a").Output()
// 	// fmt.Println(string(out))

// 	// execCmd := fmt.Sprintf("docker exec %v %v", container.ID, "./zokrates")
// 	// out, err = exec.Command("sh", "-c", execCmd).Output()
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(string(out))
// 	// out, err = exec.Command("sh", "-c", "docker rm -f "+container.ID).Output()
// 	// fmt.Println(string(out))

// 	// createExecOptions := docker.CreateExecOptions{Cmd: []string{"ls"}, AttachStdout: true, Container: container.ID}
// 	// exec, err := docker.CreateExec(&createContainerOptions)
// 	// inspect, err := docker.ExecInspect(exec.ID)
// 	// fmt.Println(inspect.)

// 	// // List images
// 	// for _, img := range imgs {
// 	// 	fmt.Println("ID: ", img.ID)
// 	// 	fmt.Println("RepoTags: ", img.RepoTags)
// 	// 	fmt.Println("Created: ", img.Created)
// 	// 	fmt.Println("Size: ", img.Size)
// 	// 	fmt.Println("VirtualSize: ", img.VirtualSize)
// 	// 	fmt.Println("ParentId: ", img.ParentID)
// 	// }
// }
