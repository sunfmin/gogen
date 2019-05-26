package gogen_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/theplant/testingutils"

	. "github.com/sunfmin/gogen"
)

/*
File means a go source file, set it's package name, and continue it's source code blocks,
like Imports, Structs, Funcs

Block can take a template and will replace it with passed in variables.
*/
func ExampleFile_01Simple() {
	var js = "js"
	f := File("api.go").Package("simple").Body(
		Imports(
			`. "github.com/theplant/htmlgo"`,
			"fmt",
			"strings",
		).Body(
			ImportAs(js, "encoding/json"),
		),

		Snippet(`
				var global int
				const name = "1231"`),

		Struct("Hello").
			FieldsSnippet(`
				Name $pkg.Marshaler
				Person *Person`, "$pkg", js).
			Fields(
				Field(
					"Age",
					"int",
					Tag("gorm", "type:varchar(100);unique_index"),
					Tag("json", "-"),
				),
			).
			Funcs(
				Func("NameLength(name string) (r int, err error)").BodySnippet(`
						return this`),
			).
			ReceiverVar("this").
			Pointer(false),

		Func("func Hello$Type(name string, age *int) (r int, err error)", "$Type", "Golang").BodySnippet(`
				if len(a) > 0 {
					fmt.Println("yes")
				} else if len(a) > 10 {
					fmt.Println("yes!")
				}`,
		),

		If(true,
			Func("func nice()").Body(
				If(true,
					Snippet(`
				ctx = graphql.WithResolverContext(ctx, &graphql.ResolverContext{
					Object: $objectName,
				})`, "$objectName", Quote("MyObject")),
				),
			),
		),
	).BodySnippet(
		`const age = 2`,
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

/*
Write a loop to add SwitchBlock Cases to write more cases
*/
func ExampleFile_02Switch() {

	var strs = []string{"one", "tw\"o", "three"}

	sw := SwitchBlock("switch x")

	for _, s := range strs {
		sw.CasesSnippet(`
			case $v:
				fmt.Println($v)`, "$v", Quote(s))
	}

	sw.Default(`
	default:
		fmt.Println(x, "default")
`)

	f := File("").Package("main").Body(
		Imports("fmt"),
		Func("main()").Body(
			Snippet(`var x = "hello"`),
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

/*
Define a interface type with either one big Block, or add FuncDecl one by one
*/
func ExampleFile_03Interface() {

	f := File("hello.go").Package("main").Body(
		Imports("fmt"),
		Interface("Writer").BodySnippet(`
			Name() string
		`).Body(
			Snippet(`Write() error`),
		),
		Func("main()").BodySnippet(
			`var x Writer
			fmt.Println(x)`,
		),
		Func("").Sig(
			FuncSig("Hello").
				Parameters("name", "string", "count", "int").
				Results("r", "bool"),
		).BodySnippet(`
			return true
		`),
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

func Hello(name string, count int) (r bool) {

	return true

}
`
	diff := testingutils.PrettyJsonDiff(expected, f.MustString(context.Background()))

	fmt.Println(diff)
	//Output:
	//

}

/*
Generate If else blocks
*/
func ExampleFile_04IfBlock() {

	f := File("hello.go").Package("main").Body(
		Imports("fmt"),
		Func("main()").Body(
			Snippet(`var x = 100
			fmt.Println(x)`),
			IfBlock("$var > 0", "$var", "x").ThenSnippet(
				`fmt.Println("x > 0")`,
			).ElseIf("x > 10 && x < 20").ThenSnippet(
				`fmt.Println("x > 10 and x < 20")`,
			).ElseIf("x > 20").Then(
				IfBlock("x == 5"),
				Snippet(`fmt.Println("x > 20")`),
			).ElseSnippet(
				`fmt.Println("else")`,
			),
		),
	)
	expected := `package main

import (
	"fmt"
)

func main() {
	var x = 100
	fmt.Println(x)
	if x > 0 {
		fmt.Println("x > 0")
	} else if x > 10 && x < 20 {
		fmt.Println("x > 10 and x < 20")
	} else if x > 20 {
		if x == 5 {
		}
		fmt.Println("x > 20")
	} else {
		fmt.Println("else")
	}

}
`
	diff := testingutils.PrettyJsonDiff(expected, f.MustString(context.Background()))

	fmt.Println(diff)
	//Output:
	//

}

/*
Generate For blocks
*/
func ExampleFile_05ForBlock() {

	f := File("hello.go").Package("main").Body(
		Imports("fmt"),
		Func("main()").Body(
			ForBlock("").BodySnippet(
				`fmt.Println("hello")`,
			),
			Snippet(`var strs = []string{"1", "2", "3"}`),
			ForBlock("_, x := range strs").Body(
				Snippet(`fmt.Println("hello", x)`),
			),
		),
	)
	expected := `package main

import (
	"fmt"
)

func main() {
	for {
		fmt.Println("hello")
	}
	var strs = []string{"1", "2", "3"}
	for _, x := range strs {
		fmt.Println("hello", x)
	}
}
`
	diff := testingutils.PrettyJsonDiff(expected, f.MustString(context.Background()))

	fmt.Println(diff)
	//Output:
	//

}

/*
Nil Code ignored
*/
func ExampleFile_06Nil() {

	f := File("hello.go").Package("main").Body(
		Struct("Hello").Fields(
			LineComment("hello"),
			Field("Name", "string"),
			LineComment(""),
			Field("Age", "int"),
		),
	)
	expected := `package main

type Hello struct {
	// hello
	Name	string
	Age	int
}
`
	diff := testingutils.PrettyJsonDiff(expected, f.MustString(context.Background()))

	fmt.Println(diff)
	//Output:
	//

}

/*
Consts example
*/
func ExampleFile_07Consts() {

	f := File("hello.go").Package("main").Body(
		ConstBlock().Type("Status", "string").Consts(
			Snippet(`StatusError Status = "Error"`),
			Const("OK", "OK"),
		),

		ConstBlock().Type("HTTPStatus", "int").Consts(
			Snippet(`HTTPStatusOK HTTPStatus = 200`),
			Const("Created", 201),
			Const("NotFound", 404),
		),
	)
	expected := `package main

type Status string

const (
	StatusError	Status	= "Error"
	StatusOK	Status	= "OK"
)

type HTTPStatus int

const (
	HTTPStatusOK		HTTPStatus	= 200
	HTTPStatusCreated	HTTPStatus	= 201
	HTTPStatusNotFound	HTTPStatus	= 404
)
`
	diff := testingutils.PrettyJsonDiff(expected, f.MustString(context.Background()))

	fmt.Println(diff)
	//Output:
	//

}

/*
More about snippet
*/
func ExampleFile_08MoreSnippet() {

	vals := []string{"Newhope", "Empire", "Jedi"}

	valsBlock := Snippets().Separator(",\n", true)
	for _, v := range vals {
		valsBlock.AppendSnippet("$Type$Val", "$Type", "Episode", "$Val", v)
	}

	f := File("hello.go").Package("main").Body(
		Snippet(`
		var All$Type = []$Type {
			$Vals
		}
		`).Var("$Type", "Episode").
			VarCode("$Vals", valsBlock),
	)
	expected := `package main

var AllEpisode = []Episode{
	EpisodeNewhope,
	EpisodeEmpire,
	EpisodeJedi,
}
`
	diff := testingutils.PrettyJsonDiff(expected, f.MustString(context.Background()))

	fmt.Println(diff)
	//Output:
	//

}
