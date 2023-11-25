package main

import (
	"embed"
	"encoding/xml"
	"fmt"
)

// <!ELEMENT protocol (copyright?, description?, interface+)>
// <!ATTLIST protocol name CDATA #REQUIRED>
type Protocol struct {
	Name        string      `xml:"name,attr"`
	Description *string     `xml:"description,omitempty"`
	Copyright   *string     `xml:"copyright,omitempty"`
	Interfaces  []Interface `xml:"interface"`
}

// <!ELEMENT interface (description?,(request|event|enum)+)>
// <!ATTLIST interface name CDATA #REQUIRED>
// <!ATTLIST interface version CDATA #REQUIRED>
type Interface struct {
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
}

//go:embed protocols/*
var ProtFS embed.FS

func main() {
	var err error
	var prot Protocol
	var bytes []byte

	bytes, err = ProtFS.ReadFile("protocols/wayland.xml")
	if err != nil {
		fmt.Println("Failed to read XML file", err)
		return
	}

	err = xml.Unmarshal(bytes, &prot)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return
	}

	fmt.Println(prot)
}
