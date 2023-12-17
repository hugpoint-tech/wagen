package main

import (
	"fmt"
	"os"
	"strings"
)

func cleanText(input string) string {

	lines := strings.Split(input, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	cleanedText := strings.Join(lines, "\n")
	cleanedText = strings.TrimSpace(cleanedText)

	return cleanedText
}

func AssertOk(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	protocols := make([]Protocol, 0, 10)

	protocols = append(protocols, ReadProtocol("protocols/wayland.xml"))
	protocols = append(protocols, ReadProtocol("protocols/stable/xdg-shell/xdg-shell.xml"))
	protocols = append(protocols, ReadProtocol("protocols/unstable/linux-dmabuf/linux-dmabuf-unstable-v1.xml"))

	// Fixups
	for i := range protocols {
		protocols[i].Copyright = cleanText(protocols[i].Copyright)

		for j := range protocols[i].Interfaces {
			iface := &protocols[i].Interfaces[j]
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

	//tmplFiles, err := filepath.Match("./*")
	//AssertOk("failed to locate template files", err)

	templates := ReadTemplates()
	ExecuteTemplates(templates, protocols)

}
