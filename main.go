package main

import (
	"fmt"
	"os"
	exe "os/exec"
	str "strings"

	gk "github.com/jwalton/gchalk"
)

var inputFile string = "."
var targetOs string = os.Getenv("GOOS")
var targetArch string = os.Getenv("GOARCH")
var out string = ""

var knownOS = map[string]bool{
	"aix":       true,
	"android":   true,
	"darwin":    true,
	"dragonfly": true,
	"freebsd":   true,
	"hurd":      true,
	"illumos":   true,
	"ios":       true,
	"js":        true,
	"linux":     true,
	"nacl":      true,
	"netbsd":    true,
	"openbsd":   true,
	"plan9":     true,
	"solaris":   true,
	"wasip1":    true,
	"windows":   true,
	"zos":       true,
}

var knownArch = map[string]bool{
	"386":         true,
	"amd64":       true,
	"amd64p32":    true,
	"arm":         true,
	"armbe":       true,
	"arm64":       true,
	"arm64be":     true,
	"loong64":     true,
	"mips":        true,
	"mipsle":      true,
	"mips64":      true,
	"mips64le":    true,
	"mips64p32":   true,
	"mips64p32le": true,
	"ppc":         true,
	"ppc64":       true,
	"ppc64le":     true,
	"riscv":       true,
	"riscv64":     true,
	"s390":        true,
	"s390x":       true,
	"sparc":       true,
	"sparc64":     true,
	"wasm":        true,
}

func PrintUsage() {
	fmt.Println(gk.Yellow("Usage:") + " gob --option=<option> <file (optional)>")
	fmt.Println("Use gob --help for a list of commands")
}

func ParseArgs() {
	for _, arg := range os.Args[1:] {
		if str.HasPrefix(arg, "--") {

			if !str.Contains(arg, "=") {
				switch arg {
				case "--help":
					fmt.Println(gk.WithBold().Red("HELP: ") + `
				--help: show this help
				--os: set target platform OS
				--arch: set target platform CPU architecture
				--out: set output file
				--list-os: list available OS
				--list-arch: list available CPU architectures
				Example:
				gob --os=linux --arch=amd64 --out=linux_amd64 main.go`)
					os.Exit(0)
				case "--list-os":
					fmt.Println(gk.Green("Available OS:"))
					for k := range knownOS {
						fmt.Println(k)
					}
					os.Exit(0)
				case "--list-arch":
					fmt.Println(gk.Green("Available architectures:"))
					for k := range knownArch {
						fmt.Println(k)
					}
					os.Exit(0)
				}

				PrintUsage()
				os.Exit(1)
			}

			option := arg[:str.Index(arg, "=")][2:]

			value := arg[str.Index(arg, "=")+1:]
			switch option {

			case "os":
				if !knownOS[value] {
					fmt.Println(gk.Red("Error: ") + "OS " + value + " does not exist")
					os.Exit(1)
				}
				targetOs = value
			case "arch":
				if !knownArch[value] {
					fmt.Println(gk.Red("Error: ") + "Arch " + value + " does not exist")
					os.Exit(1)
				}
				targetArch = value
			case "out":
				out = value
			default:
				PrintUsage()
			}
			continue
		} else {
			if inputFile == "." {
				inputFile = arg
			} else {
				fmt.Println(gk.Red("Error: ") + "Too many arguments")
				PrintUsage()
				os.Exit(1)
				return
			}
		}
	}
}

func main() {
	fmt.Println(gk.WithGreen().Bold("GOB:") + " go cross-platform builder tool (--help for options)")

	ParseArgs()

	if _, err := os.Open(inputFile); err != nil && inputFile != "." {
		fmt.Println(gk.Red("Error: ") + err.Error())
		os.Exit(1)
	}

	fmt.Println(gk.Green("Target OS: ") + targetOs)
	fmt.Println(gk.Green("Target Arch: ") + targetArch)
	fmt.Println(gk.Green("Output file: ") + out)
	fmt.Println(gk.Green("Input File: ") + inputFile)

	os.Setenv("GOOS", targetOs)
	os.Setenv("GOARCH", targetArch)

	fmt.Println(gk.Yellow("Building... "))
	cmd := exe.Command("go", "build", "-o", out, inputFile)
	if bytes, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(gk.Red("Building Error: ") + string(bytes) + "\n" + err.Error())
		os.Exit(1)
	}
}
