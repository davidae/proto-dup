package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outDir := flag.String("out", "", "output of duplicate .proto file")
	packageName := flag.String("package", "", "set a new package name")
	goPackageName := flag.String("go_package", "", "set a new go_package name")
	postfix := flag.String("postfix", "", "add a postfix to all message and enum types")
	prefix := flag.String("prefix", "", "add a prefix to all message and enum types")

	addOptional := flag.Bool("add-optional", false, "add optional to all applicable fields")
	removeOptional := flag.Bool("remove-optional", false, "remove optional from all applicable fields")
	flag.Parse()

	if *addOptional && *removeOptional {
		log.Fatal("invalid arguments, cannot add and remove optional fields at the same time")
	}

	if len(flag.Args()) != 1 {
		log.Fatal("expected one argument, a .proto file, after any flags")
	}

	filename := flag.Args()[0]
	if filepath.Ext(filename) != ".proto" {
		log.Fatal("expected a .proto file as argument")
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("unable to open file %s: %s", filename, err)
	}
	modifier := &Modifier{
		fieldToAddPostPreFix: make(map[string]bool),
		packageName:          *packageName,
		goPackageName:        *goPackageName,
		postfix:              *postfix,
		prefix:               *prefix,
		addOptional:          *addOptional,
		removeOptional:       *removeOptional,
	}

	out := duplicateFile(modifier, f)
	if *outDir == "" {
		fmt.Println(out)
		return
	}

	if err := ioutil.WriteFile(*outDir, []byte(out), 0644); err != nil {
		log.Fatalf("could not write file to %s: %s", *outDir, err)
	}
}

func duplicateFile(m *Modifier, f io.Reader) string {
	scanner := bufio.NewScanner(f)
	var out string

	for scanner.Scan() {
		out += m.modifyLine(scanner.Text()) + "\n"
	}

	out = out[:len(out)-1] // remove last \n

	return m.addPostOrPrefixToFields(out)
}

type Modifier struct {
	skipUntilCloseBracket bool
	fieldToAddPostPreFix  map[string]bool

	packageName   string
	goPackageName string
	postfix       string
	prefix        string

	removeOptional bool
	addOptional    bool
}

func (m Modifier) ChangePackageName() bool {
	return m.packageName != ""
}

func (m Modifier) ChangeGoPackageName() bool {
	return m.goPackageName != ""
}

func (m Modifier) AddPostfix() bool {
	return m.postfix != ""
}

func (m Modifier) AddPrefix() bool {
	return m.prefix != ""
}

func (m *Modifier) addPostOrPrefixToFields(proto string) string {
	builder := strings.Builder{}
	lines := strings.Split(proto, "\n")
	for i, line := range lines {
		var (
			fields         = strings.Fields(line)
			equalIndex     = 2
			fieldTypeIndex = 0
			minFields      = 3
		)

		if len(fields) > 0 && fields[0] == "optional" {
			equalIndex = 3
			fieldTypeIndex = 1
			minFields = 4
		}

		if len(fields) >= minFields && m.fieldToAddPostPreFix[fields[fieldTypeIndex]] && fields[equalIndex] == "=" {
			fieldType := fields[fieldTypeIndex]
			index := strings.Index(line, fieldType)
			builder.WriteString(line[:index] + m.prefix + fieldType + m.postfix + line[len(fieldType)+index:] + "\n")
		} else {
			if i == len(lines)-1 {
				builder.WriteString(line)
			} else {
				builder.WriteString(line + "\n")
			}
		}
	}

	return builder.String()
}

func (m *Modifier) modifyLine(line string) string {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return line
	}

	if m.skipUntilCloseBracket {
		if fields[len(fields)-1] == "}" {
			m.skipUntilCloseBracket = false
		}
		return line
	}

	if fields[0] == "package" && m.ChangePackageName() {
		return "package " + m.packageName + ";"
	}

	if fields[0] == "option" && fields[1] == "go_package" && m.ChangeGoPackageName() {
		return "option go_package =\"" + m.goPackageName + "\";"
	}

	if fields[0] == "option" || fields[0] == "import" {
		return line
	}

	if fields[0] == "optional" && m.removeOptional {
		return strings.Replace(line, "optional ", "", -1)
	}

	if fields[0] == "oneof" {
		m.skipUntilCloseBracket = true
		return line
	}

	if fields[0] == "enum" {
		m.skipUntilCloseBracket = true
		m.fieldToAddPostPreFix[fields[1]] = true
		index := strings.Index(line, fields[1])

		oldType := fields[1]
		newType := m.prefix + oldType + m.postfix

		return line[:index] + newType + line[len(oldType)+index:]
	}

	if fields[0] == "message" {
		m.fieldToAddPostPreFix[fields[1]] = true

		oldType := fields[1]
		newType := m.prefix + oldType + m.postfix
		index := strings.Index(line, oldType)

		return line[:index] + newType + line[len(oldType)+index:]
	}

	if len(fields) >= 3 && fields[2] == "=" && fields[0] != "optional" && m.addOptional {
		index := strings.Index(line, fields[0])
		line = line[:index] + "optional " + line[index:]
	}

	return line
}
