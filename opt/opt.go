// Package opt defines option and parameter types.
package opt

import (
	"fmt"
	"strings"
)

// OptionType is a type for our options. These can be provided with optional parameters or can be already set in commandMap
type OptionType int

// OptionList is a slice of OptionTypes
type OptionList []OptionType

// List of OptionTypes
const (
	DeleteSource  OptionType = iota + 1 // Delete source file/object
	IfNotExists                         // Run only if destination does not exist
	Parents                             // Just like cp --parents
	RR                                  // Reduced-redundancy
	IA                                  // Infrequent-access
	Recursive                           // Recursive copy/move (local)
	ListETags                           // Include ETags in listing
	HumanReadable                       // Human Readable file sizes (ls, du)
)

var optionsHelpOrder = [...]OptionType{
	IfNotExists,
	Parents,
	Recursive,
	RR,
	IA,
	ListETags,
	HumanReadable,
}

// Has determines if the opt.OptionList contains this OptionType
func (l OptionList) Has(check OptionType) bool {
	for _, i := range l {
		if i == check {
			return true
		}
	}
	return false
}

// GetParam returns the string/command parameter representation of a specific OptionType
func (o OptionType) GetParam() string {
	switch o {
	case IfNotExists:
		return "-n"
	case Parents:
		return "--parents"
	case RR:
		return "-rr"
	case IA:
		return "-ia"
	case Recursive:
		return "-R"
	case ListETags:
		return "-e"
	case HumanReadable:
		return "-h"
	}
	return ""
}

// HelpMessage returns the help message for a specific OptionType
func (o OptionType) HelpMessage() string {
	switch o {
	case IfNotExists:
		return "Do not overwrite existing files/objects (no-clobber)"
	case Parents:
		return "Create directory structure in destination, starting from the first wildcard"
	case Recursive:
		return "Recursive operation"
	case RR:
		return "Store with Reduced-Redundancy mode"
	case IA:
		return "Store with Infrequent-Access mode"
	case ListETags:
		return "Show ETags in listing"
	case HumanReadable:
		return "Human-readable output for file sizes"
	}
	return ""
}

// GetOptionList returns a text of accepted command options with their help messages
func GetOptionList() string {
	var out []string

	// use the order in optionsHelpOrder
	for _, o := range optionsHelpOrder {
		str := fmt.Sprintf("  %-10s %s", o.GetParam(), o.HelpMessage())
		out = append(out, str)
	}

	return strings.Join(out, "\n") + "\n"
}

// GetParams runs GetParam() on an opt.OptionList and returns a concatenated string
func (l OptionList) GetParams() string {
	r := make([]string, 0)
	for _, o := range l {
		p := o.GetParam()
		if p != "" {
			r = append(r, p)
		}
	}

	j := strings.Join(r, " ")
	if j != "" {
		return " " + j
	}
	return ""
}

// ParamType is the type of our parameter. Determines how we validate the arguments.
type ParamType int

// List of ParamTypes
const (
	Unchecked          ParamType = iota // Arbitrary single parameter
	UncheckedOneOrMore                  // One or more arbitrary parameters (special case)
	S3Obj                               // Bucket or bucket + key
	S3Dir                               // Bucket or bucket + key + "/" (prefix)
	S3ObjOrDir                          // Bucket or bucket + key [+ "/"]
	S3WildObj                           // Bucket + key with wildcard
	FileObj                             // Filename
	Dir                                 // Dir name or non-existing name ("/" appended)
	FileOrDir                           // File or directory (if existing directory, "/" appended)
	Glob                                // String containing a valid glob pattern (non-S3)
)

// String returns the string representation of ParamType
func (p ParamType) String() string {
	switch p {
	case Unchecked:
		return "param"
	case UncheckedOneOrMore:
		return "param..."
	case S3Obj:
		return "s3://bucket[/object]"
	case S3Dir:
		return "s3://bucket[/object]/"
	case S3ObjOrDir:
		return "s3://bucket[/object[/]]"
	case S3WildObj:
		return "s3://bucket/wild/*/obj*"
	case FileObj:
		return "filename"
	case Dir:
		return "directory"
	case FileOrDir:
		return "file-or-directory"
	case Glob:
		return "glob-pattern*"
	default:
		return "unknown"
	}
}
