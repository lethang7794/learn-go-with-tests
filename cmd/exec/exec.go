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

func GetXmlFromCmd(filePath string) (io.Reader, error) {
	cmd := exec.Command("cat", filePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not connect to cmd stdout: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start cmd: %v", err)
	}
	all, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("could not read cmd stdout: %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("could not wait cmd: %v", err)
	}
	return bytes.NewReader(all), nil
}

func GetXmlWithGo(filePath string) (io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read xml file with Go: %v", err)
	}
	return file, nil
}

func GetMessageFromXml(reader io.Reader) string {
	var payload Payload
	if err := xml.NewDecoder(reader).Decode(&payload); err != nil {
		log.Fatalf("could not decode xml payload: %v:", err)
	}
	return payload.Message
}

func main() {
	xmlReader, err := GetXmlFromCmd("data.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(GetMessageFromXml(xmlReader))

	xmlReader, err = GetXmlWithGo("data.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(GetMessageFromXml(xmlReader))
}
