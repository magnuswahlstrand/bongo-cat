package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	content, err := ioutil.ReadFile("resources/cat.png")
	if err != nil {
		log.Fatal(err)
	}
	content2, err := ioutil.ReadFile("resources/bongo.png")
	if err != nil {
		log.Fatal(err)
	}
	content3, err := ioutil.ReadFile("resources/Bongo1.wav")
	if err != nil {
		log.Fatal(err)
	}
	content4, err := ioutil.ReadFile("resources/Bongo4.wav")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("resources/resources.go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(f, `
package resources
var Cat_png    = []byte(%q)
var Bongo_png  = []byte(%q)
var Bongo1_wav = []byte(%q)
var Bongo4_wav = []byte(%q)
	`, string(content), string(content2), string(content3), string(content4))
}
