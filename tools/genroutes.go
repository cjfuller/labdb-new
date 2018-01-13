package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

type TemplParams struct {
	MType string
}

type AllRoutesParams struct {
	Code string
}

// TODO(colin): infer this from the models directory.
var allTypes = []string{
	"plasmid",
	"oligo",
	"line",
	"sample",
	"bacterium",
	"yeaststrain",
	"user",
	"antibody",
	"rnaiclone",
	"seqlib",
}

var plurals = map[string]string{
	"bacterium": "bacteria",
	"antibody":  "antibodies",
}

func pluralize(s string) string {
	if pl, found := plurals[s]; found {
		return pl
	}
	return s + "s"
}

var outputFile = "routes/routes.go"

func genAllRoutesCode() string {
	code := &bytes.Buffer{}
	for _, t := range allTypes {
		fmt.Fprintln(code, "install"+t+"(r)")
		fmt.Fprintln(code, "install"+pluralize(t)+"(r)")
	}
	return string(code.Bytes())
}

func generateCode() {
	allroutes, err := template.ParseFiles("tools/templates/allroutes.go.template")
	if err != nil {
		panic(err)
	}
	route, err := template.ParseFiles("tools/templates/routes.go.template")
	if err != nil {
		panic(err)
	}
	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = allroutes.Execute(f, AllRoutesParams{Code: genAllRoutesCode()})
	if err != nil {
		panic(err)
	}
	for _, t := range allTypes {
		err = route.Execute(f, TemplParams{MType: t})
		// TODO(colin): be consistent about plural / singular everywhere.
		err = route.Execute(f, TemplParams{MType: pluralize(t)})
		if err != nil {
			panic(err)
		}
	}
}

func formatCode() {
	err := exec.Command("go", "fmt", outputFile).Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	generateCode()
	formatCode()
}
