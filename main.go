package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	addBashFiles()
}

func startUp() {
	//scan bash_profile for necessary lines
	fmt.Println("Hello and welcome to Shakur, please give me a second to set up")
	// function to add sxport script to bash_profile or bashrc
	// conditional to make sure setup completed, if it completes
	fmt.Println("Alright we're ready to go, which commands do you want to watch? Rememeber you can watch more later by running shakur watch 'my command'")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// uCommand := scanner.Text()
		fmt.Println(scanner.Text())
		//call function to append to the scary phrase array
		fmt.Println("uC has been added to your watch list, would you like to add another?")
		//read and then do some action based on yes / no
	}
}

func addBashFiles() {
	err := os.MkdirAll(".bash_shakur", 0744)
	if err != nil {
		fmt.Println("I can't create the folder I want, please check permissions")
		fmt.Println(err)
	}
	ioutil.WriteFile(".bash_shakur/bash_preexec.sh", copy("bash-preexec.sh"), 0744)
	ioutil.WriteFile(".bash_shakur/bash_shakur.sh", copy("bash-exe.sh"), 0744)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copy(src string) []byte {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	checkErr(err)
	return data
}
