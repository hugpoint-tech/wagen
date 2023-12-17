package main

import (
	"embed"
	"encoding/xml"
)

//go:embed protocols/*
var ProtFS embed.FS

func ReadProtocol(path string) Protocol {
	var bytes []byte
	var err error
	var prot Protocol

	bytes, err = ProtFS.ReadFile(path)
	AssertOk("Failed to read XML file", err)

	err = xml.Unmarshal(bytes, &prot)
	AssertOk("Error unmarshalling XML:", err)

	return prot
}
