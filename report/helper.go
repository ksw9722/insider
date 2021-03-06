package report

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

func cleanDRA(dir string, dra []DRA) ([]DRA, error) {
	for i := 0; i < len(dra); i++ {
		d := &dra[i]
		rel, err := filepath.Rel(dir, d.FilePath)
		if err != nil {
			return nil, err
		}
		d.FilePath = rel
	}
	return dra, nil
}

func unique(slice []DRA) []DRA {
	keys := make(map[string]bool)
	list := []DRA{}
	for _, entry := range slice {
		if _, value := keys[entry.Data]; !value {
			keys[entry.Data] = true
			list = append(list, entry)
		}
	}
	return list
}

func reportJson(r interface{}, out io.Writer) error {
	b, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return err
	}

	if _, err := out.Write(b); err != nil {
		return fmt.Errorf("write bytes: %w", err)
	}

	return nil
}

func reportHTML(t string, r interface{}, out io.Writer) error {
	tmpl, err := template.New("report").Parse(t)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(out, r); err != nil {
		return err
	}

	return downloadCss()
}

func resumeReport(score float64, dra, vulnerabilities, none, low, medium, high, critical, total int, out io.Writer) {
	fmt.Fprintf(out, "\n-----------------------------------------------\n")
	fmt.Fprintf(out, "Score Security %v/100\n", score)
	fmt.Fprintf(out, "Vulnerabilities\t%3v \n", vulnerabilities)
	fmt.Fprintf(out, "DRA\t\t%3v \n", dra)
	fmt.Fprintf(out, "None\t\t%3v \n", none)
	fmt.Fprintf(out, "Low\t\t%3v \n", low)
	fmt.Fprintf(out, "Medium\t\t%3v \n", medium)
	fmt.Fprintf(out, "High\t\t%3v \n", high)
	fmt.Fprintf(out, "Critical\t%3v \n", critical)
	fmt.Fprintf(out, "Total\t\t%3v \n", total)
	fmt.Fprintln(out, "-----------------------------------------------------------------------------------------------------------------------")
	fmt.Fprintln(out, "You are using the Insider open source version. If you like the product and want more features,")
	fmt.Fprintln(out, "visit http://insidersec.io and get to know our enterprise version.")
	fmt.Fprintln(out, "If you are a developer, then you can contribute to the improvement of the software while using an open source version")
	fmt.Fprintln(out, "-----------------------------------------------------------------------------------------------------------------------")
}

func consoleReport(score float64, dra []DRA, libraries []Library, vulnerabilities []Vulnerability, out io.Writer) {
	fmt.Fprintf(out, "\n---------------------------------------------------------------------\n")
	fmt.Fprintf(out, "Score Security %v\n\n", score)

	for i, k := range dra {
		if i == 0 {
			fmt.Fprintf(out, "DRA - Data Risk Analytics\n")
		}
		fmt.Fprintf(out, "File %s\n", k.FilePath)
		fmt.Fprintf(out, "Dra %s\n", k.Data)
		fmt.Fprintf(out, "Type %s\n", k.Type)
	}
	fmt.Fprintln(out, " ")

	for i, k := range libraries {
		if i == 0 {
			fmt.Fprintf(out, "%-20v %-10v \n", "Library", "Version")
		}
		fmt.Fprintf(out, "%-20v %-10v \n", k.Name, k.Version)
	}
	fmt.Fprintln(out, " ")

	for _, k := range vulnerabilities {
		fmt.Fprintln(out, "CVSS", k.CVSS)
		fmt.Fprintln(out, "Severity", k.Severity)
		fmt.Fprintln(out, "Class", k.Class)
		fmt.Fprintln(out, "VulnerabilityID", k.VulnerabilityID)
		fmt.Fprintln(out, "Description", k.Description)
		fmt.Fprintln(out, "ClassMessage", k.ClassMessage)
		fmt.Fprintln(out, "Recomendation", k.Recomendation)
		fmt.Fprintln(out, "")
	}

	fmt.Fprintln(out, "---------------------------------------------------------------------")

}
