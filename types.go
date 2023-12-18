package main

type ProtocolTypes int

const (
	/// 32-bit signed int
	Int ProtocolTypes = iota

	/// 32-bit unsigned int
	Uint

	/// Signed 23.8 decimal numbers. It is a signed decimal type which offers a sign bit, 23 bits of
	/// integer precision and 8 bits of decimal precision. This is exposed as an opaque struct with
	/// conversion helpers to and from double and int on the C API side.
	Fixed

	/// Starts with an unsigned 32-bit length (including null terminator), followed by the string
	/// contents, including terminating null byte, then padding to a 32-bit boundary. A null value
	/// is represented with a length of 0.
	String

	/// 32-bit object ID. A null value is represented with an ID of 0.
	Object

	/// The 32-bit object ID. Generally, the interface used for the new object is inferred from the
	/// xml, but in the case where it's not specified, a new_id is preceded by a string specifying
	/// the interface name, and a uint specifying the version.
	NewId

	/// Starts with 32-bit array size in bytes, followed by the array contents verbatim, and
	/// finally padding to a 32-bit boundary.
	Array

	/// The file descriptor is not stored in the message buffer, but in the ancillary data of the
	/// UNIX domain socket message (msg_control).
	Fd
)

type ProtocolMaturity int

const (
	Core ProtocolMaturity = iota
	Stable
	Staging
	Unstable
)

// Protocol
// <!ELEMENT protocol (copyright?, description?, interface+)>
// <!ATTLIST protocol name CDATA #REQUIRED>
type Protocol struct {
	Name        string      `xml:"name,attr"`
	Description string      `xml:"description,omitempty"`
	Copyright   string      `xml:"copyright,omitempty"`
	Interfaces  []Interface `xml:"interface"`
	Maturity    ProtocolMaturity
}

type Protocols []Protocol

// Interface
// / <!ELEMENT interface (description?,(request|event|enum)+)>
// / <!ATTLIST interface name CDATA #REQUIRED>
// / <!ATTLIST interface version CDATA #REQUIRED>
type Interface struct {
	Description string    `xml:"description,omitempty"`
	Name        string    `xml:"name,attr"`
	Version     string    `xml:"version,attr"`
	Requests    []Request `xml:"request"`
	Events      []Event   `xml:"event"`
	Enums       []Enum    `xml:"enum"`
}

// Request
// <!ELEMENT request (description?,arg*)>
// <!ATTLIST request name CDATA #REQUIRED>
// <!ATTLIST request type CDATA #IMPLIED>
// <!ATTLIST request since CDATA #IMPLIED>
type Request struct {
	Description string `xml:"description,omitempty"`
	Name        string `xml:"name,attr"`
	Type        string `xml:"type,attr,omitempty"`
	Since       string `xml:"since,attr,omitempty"`
	Args        []Arg  `xml:"arg"`
	Opcode      int
}

type Event struct {
	Description string `xml:"description,omitempty"`
	Name        string `xml:"name,attr"`
	Type        string `xml:"type,attr,omitempty"`
	Since       string `xml:"since,attr,omitempty"`
	Args        []Arg  `xml:"arg"`
	Opcode      int
}

// Enum
// <!ELEMENT enum (description?,entry*)>
// <!ATTLIST enum name CDATA #REQUIRED>
// <!ATTLIST enum since CDATA #IMPLIED>
// <!ATTLIST enum bitfield CDATA #IMPLIED>
type Enum struct {
	Description string  `xml:"description,omitempty"`
	Name        string  `xml:"name,attr"`
	Since       string  `xml:"since,attr,omitempty"`
	Bitfield    string  `xml:"bitfield,attr,omitempty"`
	Entries     []Entry `xml:"entry"`
}

// Entry
// <!ELEMENT entry (description?)>
// <!ATTLIST entry name CDATA #REQUIRED>
// <!ATTLIST entry value CDATA #REQUIRED>
// <!ATTLIST entry summary CDATA #IMPLIED>
// <!ATTLIST entry since CDATA #IMPLIED>
type Entry struct {
	Description string `xml:"description,omitempty"`
	Name        string `xml:"name,attr"`
	Value       string `xml:"value,attr"`
	Summary     string `xml:"summary,attr,omitempty"`
	Since       string `xml:"since,attr,omitempty"`
}

// Arg
// <!ELEMENT arg (description?)>
// <!ATTLIST arg name CDATA #REQUIRED>
// <!ATTLIST arg type CDATA #REQUIRED>
// <!ATTLIST arg summary CDATA #IMPLIED>
// <!ATTLIST arg interface CDATA #IMPLIED>
// <!ATTLIST arg allow-null CDATA #IMPLIED>
// <!ATTLIST arg enum CDATA #IMPLIED>
type Arg struct {
	Description string `xml:"description,omitempty"`
	Name        string `xml:"name,attr"`
	Type        string `xml:"type,attr"`
	Summary     string `xml:"summary,attr,omitempty"`
	Interface   string `xml:"interface,attr,omitempty"`
	AllowNull   bool   `xml:"allow-null,attr,omitempty"`
	Enum        string `xml:"enum,attr,omitempty"`
}
