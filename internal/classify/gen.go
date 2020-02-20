// +build ignore

// This generates stopwords.go by running "go generate"
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"unicode"

	"github.com/photoprism/photoprism/pkg/fs"
	"gopkg.in/yaml.v2"
)

// LabelRule defines the rule for a given Label
type LabelRule struct {
	Label      string
	See        string
	Threshold  float32
	Categories []string
	Priority   int
}

type LabelRules map[string]LabelRule

// This function generates the rules.go file containing rule extracted from rules.yml file
func main() {
	rules := make(LabelRules)

	fileName := "rules.yml"

	if !fs.FileExists(fileName) {
		log.Panicf("tensorflow: label rules file not found in \"%s\"", filepath.Base(fileName))
	}

	yamlConfig, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlConfig, rules)

	for label, rule := range rules {
		for _, char := range label {
			if unicode.IsUpper(char) {
				log.Panicf("label must be lowercase: %s", label)
			}
		}

		if rule.See != "" {
			rule, ok := rules[rule.See]

			if !ok {
				log.Panicf("missing label: %s", rule.See)
			}

			rules[label] = rule
		}
	}

	f, err := os.Create("rules.go")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	packageTemplate.Execute(f, struct {
		Rules LabelRules
	}{
		Rules: rules,
	})
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package classify

var rules = LabelRules{
{{- range $key, $value := .Rules }}
	{{ printf "%q" $key }}:  {
		Label:      {{ printf "%q" $value.Label }},
		Threshold:  {{ printf "%f" $value.Threshold }},
		Priority:   {{ $value.Priority }},
		Categories: []string{ {{- range $value.Categories }} {{ printf "%q" . }}, {{- end }} },
	},
{{- end }}
}`))
