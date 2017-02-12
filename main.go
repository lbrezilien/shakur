package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

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
