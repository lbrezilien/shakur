package main

import (
	"os"
	"github.com/progrium/go-basher"
		"log"

)

const (
	bsht = `echo "I be bash"`
)

//go:generate go run includetxt.go
func main() {
	bash, _ := basher.NewContext("/bin/bash", false)
	bash.Source("bashex.sh", nil)
	status2, err := bash.Run("loadpreex", os.Args[1:])
	status, err := bash.Run("preexec", os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(status)
	os.Exit(status2)
	//compile shell files into go-binary using go generate
	//function that can be called from command line that runs the bash files
	//test by writing function into main that reads file and returns bash script and add an eval script in bash_profile
	//write a config function that adds above function as export in users bash_profile
	// addBashFiles()
	// fmt.Println("it worked")
}

