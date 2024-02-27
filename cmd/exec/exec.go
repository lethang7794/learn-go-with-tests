package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Payload struct {
	Message string `xml:"message"`
}

func GetXmlFromCmd(filePath string) io.Reader {
	cmd := exec.Command("cat", filePath)
	stdout, _ := cmd.StdoutPipe()
	_ = cmd.Start()
	all, _ := io.ReadAll(stdout)
	_ = cmd.Wait()
	return bytes.NewReader(all)
}

func GetXmlWithGo(filePath string) io.Reader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func GetMessageFromXml(reader io.Reader) string {
	var payload Payload
	if err := xml.NewDecoder(reader).Decode(&payload); err != nil {
		log.Fatalf("could not decode xml payload: %v:", err)
	}
	return payload.Message
}

func main() {
	fmt.Println(GetMessageFromXml(GetXmlFromCmd("data.xml")))
	fmt.Println(GetMessageFromXml(GetXmlWithGo("data.xml")))
}
