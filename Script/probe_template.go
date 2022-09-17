package main

import (
	"fmt"
	"html/template"
	"os"
)

type Person struct {
	Name string
	Age    int
}

func main() {
	p := Person{"longshuai", 23}
	tmpl, err := template.New("test").Parse("Name: {{.Name}}, Age: {{.Age}}")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", tmpl)
}