package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/ernesto-jimenez/gogen/automock"
	"github.com/ernesto-jimenez/gogen/importer"
	"github.com/ernesto-jimenez/gogen/strconv"
)

var (
	out      = flag.String("o", "", "specify the name of the generated code. Default value is by generated based on the name of the variable, e.g.: DefaultClient -> default_client_funcs.go (use \"-\" to print to stdout)")
	mockName = flag.String("mock-name", "", "name for the mock")
	mockPkg  = flag.String("mock-pkg", "", "package name for the mock")
	pkg      = flag.String("pkg", ".", "what package to get the interface from")
	inPkg    = flag.Bool("in-pkg", false, "whether the mock is internal to the package")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	iface := flag.Arg(0)

	if iface == "" {
		log.Fatal("need to specify an interface name")
	}

	gen, err := automock.NewGenerator(*pkg, iface)
	if err != nil {
		log.Fatal(err)
	}

	if *mockName != "" {
		gen.SetName(*mockName)
	}
	if *pkg == "." && path.Dir(*out) == "." {
		*inPkg = true
	}
	gen.SetInternal(*inPkg)
	if *mockPkg == "" && !*inPkg {
		p, err := importer.Default().Import(".")
		if err != nil {
			log.Fatal(err)
		}
		*mockPkg = p.Name()
	}
	if *mockPkg != "" {
		gen.SetPackage(*mockPkg)
	}

	w := os.Stdout
	if *out == "" {
		*out = fmt.Sprintf("%s_test.go", gen.Name())
		if p := regexp.MustCompile(".*/").ReplaceAllString(*pkg, ""); !*inPkg && p != "" && p != "." {
			*out = p + "_" + *out
		}
	}
	if *out != "-" {
		*out = strconv.SnakeCase(*out)
		log.Printf("Generating mock for %s in %s", iface, *out)
		w, err = os.OpenFile(*out, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = gen.Write(w)
	if err != nil {
		log.Fatal(err)
	}
}
