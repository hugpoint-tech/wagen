package main

import (
	"embed"
	"encoding/xml"
	"io/fs"
	"path/filepath"
	"strings"
)

var (
	//go:embed protocols/*
	protfs    embed.FS
	PROTOCOLS = make(map[string]Protocol)
)

func init() {
	var err error
	var bytes []byte

	err = fs.WalkDir(protfs, ".", func(path string, d fs.DirEntry, err error) error {
		parts := strings.Split(path, "/")
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".xml" {
			bytes, err = protfs.ReadFile(path)
			AssertOk("Failed to read XML file", err)
			var prot Protocol
			err = xml.Unmarshal(bytes, &prot)
			AssertOk("Error unmarshalling XML:", err)

			var maturity ProtocolMaturity
			switch parts[1] {
			case "wayland.xml":
				maturity = Core
			case "stable":
				maturity = Stable
			case "staging":
				maturity = Staging
			case "unstable":
				maturity = Unstable
			default:
				panic("unexpected protocol maturity")
			}
			prot.Maturity = maturity
			fixupProtocol(&prot)

			PROTOCOLS[prot.Name] = prot
		}
		return nil
	})
	AssertOk("fatal: failed to list embedded xml protocol definitions", err)
}

func fixupProtocol(prot *Protocol) {
	prot.Copyright = cleanText(prot.Copyright)

	for j := range prot.Interfaces {
		iface := &prot.Interfaces[j]
		iface.Description = cleanText(iface.Description)

		for reqopcode := range iface.Requests {
			req := &iface.Requests[reqopcode]
			req.Description = cleanText(req.Description)
			req.Opcode = reqopcode

			for k := range req.Args {
				arg := &req.Args[k]
				arg.Summary = cleanText(arg.Summary)
				arg.Description = cleanText(arg.Description)
			}
		}

		for eventopcode := range iface.Events {
			event := &iface.Events[eventopcode]
			event.Description = cleanText(event.Description)
			event.Opcode = eventopcode

			for l := range event.Args {
				arg := &event.Args[l]
				arg.Summary = cleanText(arg.Summary)
				arg.Description = cleanText(arg.Description)
			}
		}

		for m := range iface.Enums {
			enum := &iface.Enums[m]
			enum.Description = cleanText(enum.Description)

			for n := range enum.Entries {
				entry := &enum.Entries[n]
				entry.Description = cleanText(entry.Description)
				entry.Summary = cleanText(entry.Summary)
			}
		}
	}
}

func cleanText(input string) string {

	lines := strings.Split(input, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	cleanedText := strings.Join(lines, "\n")
	cleanedText = strings.TrimSpace(cleanedText)

	return cleanedText
}
