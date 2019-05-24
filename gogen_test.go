package gogen_test

import (
	"context"
	"fmt"

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
				const name = "1231"
			`),
		Struct("Hello").Block(`
				Name $pkg.Marshaler
				Person *Person
			`, "$pkg", js).AppendField("Age", "int",
			Tag("gorm", "type:varchar(100);unique_index"),
			Tag("json", "-"),
		).Funcs(
			Func("NameLength(name string) (r int, err error)").Block(`
						return this
					`),
		).ReceiverVar("this").Pointer(false),
		Func("func Hello$Type(name string, age *int) (r int, err error)", "$Type", "Golang").Blocks(
			Block(`
				if len(a) > 0 {
					fmt.Println("yes")
				} else if len(a) > 10 {
					fmt.Println("yes!")
				}			
				`),
		),

		If(true,
			Func("func nice()").Blocks(
				If(true,
					Block(`
				ctx = graphql.WithResolverContext(ctx, &graphql.ResolverContext{
					Object: $objectName,
				})
				`, "$objectName", Quote("MyObject")),
				),
			),
		),
	).Blocks(
		Block(`
				const age = 2
			`),
	)

	fmt.Println(MustString(f, context.Background()))
	//Output:
	//123
}
