package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {
	line := strings.Split("openssl x509 -in ./cert.pem -noout -enddate", " ")
	cmd := exec.Command(line[0], line[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	out = bytes.Trim(out, "notAfter=")
	out = bytes.TrimSpace(out)

	const layout = "Jan 2 15:04:05 2006 MST"
	notAfter, err := time.Parse(layout, string(out))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("notAfter is", notAfter)
	if notAfter.After(time.Now()) {
		fmt.Println("Certification is valid")
	} else {
		fmt.Println("Certification is invalid")
	}
}
