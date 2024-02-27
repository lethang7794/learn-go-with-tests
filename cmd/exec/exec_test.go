package main

import (
	"log"
	"strings"
	"testing"
)

func TestGetMessageFromXML(t *testing.T) {
	t.Run("returns Happy New Year!", func(t *testing.T) {
		in := strings.NewReader(`<payload>
    <message>Happy New Year!</message>
</payload>`)
		data := GetMessageFromXml(in)

		got := data
		want := "Happy New Year!"
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("returns Happy Birthday!", func(t *testing.T) {
		in := strings.NewReader(`<payload>
    <message>Happy Birthday!</message>
</payload>`)
		data := GetMessageFromXml(in)

		got := data
		want := "Happy Birthday!"
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestGetMessageFromXml_XmlReadWithCmd(t *testing.T) {
	t.Run("returns Happy New Year!", func(t *testing.T) {
		xml, err := GetXmlFromCmd("testdata/msg.xml")
		if err != nil {
			log.Fatal(err)
		}
		data := GetMessageFromXml(xml)

		got := data
		want := "Happy New Year!"
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestGetMessageFromXml_XmlReadWithGo(t *testing.T) {
	t.Run("returns Happy New Year!", func(t *testing.T) {
		xml, err := GetXmlWithGo("testdata/msg.xml")
		if err != nil {
			log.Fatal(err)
		}
		data := GetMessageFromXml(xml)

		got := data
		want := "Happy New Year!"
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}
