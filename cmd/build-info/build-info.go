package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	var (
		flRoot = flag.String(
			"root",
			"",
			"Package the entire contents of the directory tree at root-path, typically a destination root created by xcodebuild(1).",
		)
		flIdentifier = flag.String(
			"identifier",
			"",
			`Specify a unique identifier for this package. The OS X Installer recognizes a package as being an upgrade to an already-installed package only if the package identifiers match, so it is advisable
	to set a meaningful, consistent identifier when you build the package.  pkgbuild will infer an identifier when building a package from a single component, but will fail otherwise if the identifier
	has not been set.`,
		)
		flPkgVersion = flag.String(
			"version",
			"",
			`Specify a version for the package. Packages with the same identifier are compared using this version, to determine if the package is an upgrade or downgrade. 
	If you don't specify a version, a default of zero is assumed, but this may prevent proper upgrade/downgrade checking.`,
		)
		flOut = flag.String("output", "", "path to output file")
	)
	flag.Parse()

	size, count, err := fileInfo(*flRoot)
	if err != nil {
		log.Fatal(err)
	}

	pkginfo := &info{
		Identifier: *flIdentifier,
		Version:    *flPkgVersion,
		Size:       size,
		NumFiles:   count,
	}

	os.Remove(*flOut)
	os.MkdirAll(filepath.Dir(*flOut), os.ModePerm)
	out, err := os.Create(*flOut)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	funcMap := template.FuncMap{"scripts": pkginfo.Scripts}
	tmpl := template.Must(template.New("").
		Funcs(funcMap).
		Parse(pkginfoTemplate),
	)
	if err := tmpl.Execute(out, pkginfo); err != nil {
		log.Fatal(err)
	}

}

type info struct {
	Identifier string
	Version    string
	Size       int64
	NumFiles   int64
}

func (i *info) Scripts() string {
	if fi, err := os.Stat("scripts"); os.IsNotExist(err) {
		return ""
	} else if err != nil {
		panic(err)
	} else if !fi.IsDir() {
		panic("scripts must be a directory")
	}

	dirFiles, err := ioutil.ReadDir("scripts")
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.WriteString("    <scripts>\n")
	for _, f := range dirFiles {
		if f.Name() == "preinstall" {
			buf.WriteString(`		<preinstall file="./preinstall"/>`)
			buf.WriteString("\n")
		}
		if f.Name() == "postinstall" {
			buf.WriteString(`        <postinstall file="./postinstall"/>`)
			buf.WriteString("\n")
		}
	}
	buf.WriteString("    </scripts>")

	return buf.String()
}

func fileInfo(path string) (int64, int64, error) {
	var size int64
	var count int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info == nil || err != nil {
			return err
		}
		count++
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size / 1024.0, count, err
}

var pkginfoTemplate = `<pkg-info format-version="2" identifier="{{ .Identifier }}" version="{{ .Version }}" install-location="/" auth="root">
  <payload installKBytes="{{ .Size }}" numberOfFiles="{{ .NumFiles }}"/>
{{ scripts }}
    <bundle-version/>
    <upgrade-bundle/>
    <update-bundle/>
    <atomic-update-bundle/>
    <strict-identifier/>
    <relocate/>
</pkg-info>`
