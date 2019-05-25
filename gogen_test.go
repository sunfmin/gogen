package gogen_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/theplant/testingutils"

	. "github.com/sunfmin/gogen"
)

func ExamplePackage_01Simple() {
	var js = "js"
	f := File("api.go").Package("simple").Blocks(
		Imports(
			`. "github.com/theplant/htmlgo"`,
			"fmt",
			"strings",
		).Blocks(
			ImportAs(js, "encoding/json"),
		),
		Block(`
				var global int
				const name = "1231"`),
		Struct("Hello").Block(`
				Name $pkg.Marshaler
				Person *Person`, "$pkg", js).AppendField("Age", "int",
			Tag("gorm", "type:varchar(100);unique_index"),
			Tag("json", "-"),
		).Funcs(
			Func("NameLength(name string) (r int, err error)").Block(`
						return this`),
		).ReceiverVar("this").Pointer(false),
		Func("func Hello$Type(name string, age *int) (r int, err error)", "$Type", "Golang").Blocks(
			Block(`
				if len(a) > 0 {
					fmt.Println("yes")
				} else if len(a) > 10 {
					fmt.Println("yes!")
				}`),
		),

		If(true,
			Func("func nice()").Blocks(
				If(true,
					Block(`
				ctx = graphql.WithResolverContext(ctx, &graphql.ResolverContext{
					Object: $objectName,
				})`, "$objectName", Quote("MyObject")),
				),
			),
		),
	).Blocks(
		Block(`const age = 2`),
	)
	expected := `package simple

import (
	. "github.com/theplant/htmlgo"
	"fmt"
	"strings"
	js "encoding/json"
)

var global int

const name = "1231"

type Hello struct {
	Name	js.Marshaler
	Person	*Person
	Age	int	$Qgorm "type:varchar(100);unique_index" json "-"$Q
}

func (this Hello) NameLength(name string) (r int, err error) {

	return this
}

func HelloGolang(name string, age *int) (r int, err error) {

	if len(a) > 0 {
		fmt.Println("yes")
	} else if len(a) > 10 {
		fmt.Println("yes!")
	}
}

func nice() {

	ctx = graphql.WithResolverContext(ctx, &graphql.ResolverContext{
		Object: "MyObject",
	})
}

const age = 2
`
	diff := testingutils.PrettyJsonDiff(strings.ReplaceAll(expected, "$Q", "`"), f.MustString(context.Background()))
	fmt.Println(diff)
	//Output:
	//
}

func ExamplePackage_02Switch() {

	var strs = []string{"one", "tw\"o", "three"}

	sw := SwitchBlock("switch x")

	for _, s := range strs {
		sw.Cases(Block(`
			case $v:
				fmt.Println($v)`, "$v", Quote(s)))
	}

	sw.Default(`
	default:
		fmt.Println(x, "default")
`)

	f := File("").Package("main").Blocks(
		Imports("fmt"),
		Func("main()").Blocks(
			Block(`var x = "hello"`),
			sw,
		),
	)
	expected := `package main

import (
	"fmt"
)

func main() {
	var x = "hello"
	switch x {

	case "one":
		fmt.Println("one")

	case "tw\"o":
		fmt.Println("tw\"o")

	case "three":
		fmt.Println("three")

	default:
		fmt.Println(x, "default")

	}

}
`
	diff := testingutils.PrettyJsonDiff(expected, f.MustString(context.Background()))

	fmt.Println(diff)
	//Output:
	//

}

func ExamplePackage_03Interface() {

	f := File("").Package("main").Blocks(
		Imports("fmt"),
		Interface("Writer").Block(`
			Name() string
		`).AppendFuncDecl(`Write() error`),
		Func("main()").Blocks(
			Block(`var x Writer
			fmt.Println(x)`),
		),
	)
	expected := `package main

import (
	"fmt"
)

type Writer interface {
	Name() string

	Write() error
}

func main() {
	var x Writer
	fmt.Println(x)
}
`
	diff := testingutils.PrettyJsonDiff(expected, f.MustString(context.Background()))

	fmt.Println(diff)
	//Output:
	//

}
