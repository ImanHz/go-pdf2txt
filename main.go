package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"code.sajari.com/docconv/client"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Not enough parameters. see -h for help")
		return
	}

	command := os.Args[1]

	if command == "-h" {

		fmt.Println("use pdf2txt -c filename.pdf to convert pdf to text")
		fmt.Println("use pdf2txt -a to convert all pdf files in the current folder to text")
		return

	}

	fmt.Println("Starting docd service!")

	ch := make(chan bool)

	go runDocd(ch)

	switch command {

	case "-c":

		if len(os.Args) <= 2 {
			fmt.Println("Enter the pdf file name")
			return
		}

		convertPDF(os.Args[2])

		fmt.Println("finished!")

	case "-a":
		pwd, _ := os.Getwd()
		files, err := os.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if (strings.Contains(file.Name(), ".pdf")) && !file.IsDir() {
				fmt.Println("converting ", file.Name())
				convertPDF(file.Name())
			}

		}

	}

	fmt.Println("Closing docd process")

	ch <- true

	fmt.Println("Finished converting!")

}

func convertPDF(filename string) {

	c := client.New()

	res, err := client.ConvertPath(c, filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	str := res.Body

	str = strings.ReplaceAll(str, ")", "~")
	str = strings.ReplaceAll(str, "(", ")")
	str = strings.ReplaceAll(str, "~", "(")

	f, err := os.Create(strings.ReplaceAll(filename, ".pdf", ".txt"))

	if err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()

	_, err2 := f.WriteString(str)

	if err2 != nil {
		log.Fatal(err2)
		return
	}

}

func runDocd(ch chan bool) {
	// s := []string{"cmd.exe", "/C", "start", "/b", `docd.exe`}

	cmd := exec.Command("docd.exe")
	if err := cmd.Start(); err != nil {
		log.Println("Error:", err)

	}
	<-ch
	if err := cmd.Process.Kill(); err != nil {
		fmt.Println("failed to kill process: ", err)
		return
	}

}
