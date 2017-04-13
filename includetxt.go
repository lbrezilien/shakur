// package main

// import (
// 	"io"
// 	"io/ioutil"
// 	"os"
// 	"strings"
// )

// func main() {
// 	convertBashToGo()
// }

// func convertBashToGo() {
// 	fs, _ := ioutil.ReadDir(".")
// 	out, _ := os.Create("bashfiles.go")
// 	out.Write([]byte("package main \n\nconst (\n"))
// 	for _, f := range fs {
// 		if strings.HasSuffix(f.Name(), ".sh") {
// 			out.Write([]byte(strings.TrimSuffix(f.Name(), ".sh") + " = `"))
// 			f, _ := os.Open(f.Name())
// 			io.Copy(out, f)
// 			out.Write([]byte("`\n"))
// 		}
// 	}
// 	out.Write([]byte(")\n"))
// }
