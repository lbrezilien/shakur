package main

import (
	"fmt"
	"os/exec"
)

const (
	bsht = `echo "I be bash"`
)

//go:generate go run includetxt.go
func main() {
	cmd := exec.Command("/bin/sh", "bashex.sh")
	fmt.Println(cmd.Output())
	//compile shell files into go-binary using go generate
	//function that can be called from command line that runs the bash files
	//test by writing function into main that reads file and returns bash script and add an eval script in bash_profile
	//write a config function that adds above function as export in users bash_profile
	// addBashFiles()
	// fmt.Println("it worked")
}

