package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var indexTemplate = `{{range $v := .Files}}import { {{$v}}, I{{$v}} } from "./{{$v}}";
{{end}}

export {
    {{range $v := .Files}}{{$v}},
    {{end}}
    {{range $v := .Files}}I{{$v}},
    {{end}}
}
`

func main() {
	// get args first
	if len(os.Args) != 2 {
		fmt.Println("error: require src root path")
		os.Exit(1)
	}

	curDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	srcDir := filepath.Join(curDir, os.Args[1])
	outFile := filepath.Join(curDir, os.Args[1], "index.ts")

	fmt.Printf("using sources from: %s\n", srcDir)

	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	fileNames := []string{}
	for _, file := range files {
		if file.Name() != "index.ts" && strings.HasSuffix(file.Name(), ".ts") {
			fileNames = append(fileNames, strings.TrimSuffix(file.Name(), ".ts"))
		}
	}

	t := template.Must(template.New("index").Parse(indexTemplate))
	f, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
	err = t.Execute(f, map[string]interface{}{
		"Files": fileNames,
	})
	if err != nil {
		fmt.Printf("error: generating index.ts, %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("index.ts generated successfully at: %s\n", outFile)
}
