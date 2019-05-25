

# gogen: Generate Go code with structured blocks and composition


```go
import . "github.com/sunfmin/gogen"
```



File means a go source file, set it's package name, and continue it's source code blocks,
like Imports, Structs, Funcs

Block can take a template and will replace it with passed in variables.
```go
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
```

Write a loop to add SwitchBlock Cases to write more cases
```go
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
```

Define a interface type with either one big Block, or add FuncDecl one by one
```go
	f := File("hello.go").Package("main").Blocks(
	    Imports("fmt"),
	    Interface("Writer").Block(`
	        Name() string
	    `).AppendFuncDecl(`Write() error`),
	    Func("main()").Blocks(
	        Block(`var x Writer
	        fmt.Println(x)`),
	    ),
	    Func("").Sig(
	        FuncSig("Hello").
	            Parameters("name", "string", "count", "int").
	            Results("r", "bool"),
	    ).Block(`
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
```

Generate If else blocks
```go
	f := File("hello.go").Package("main").Blocks(
	    Imports("fmt"),
	    Func("main()").Blocks(
	        Block(`var x = 100
	        fmt.Println(x)`),
	        IfBlock("$var > 0", "$var", "x").Then(
	            Block(`fmt.Println("x > 0")`),
	        ).ElseIf("x > 10 && x < 20").Then(
	            Block(`fmt.Println("x > 10 and x < 20")`),
	        ).ElseIf("x > 20").Then(
	            IfBlock("x == 5"),
	            Block(`fmt.Println("x > 20")`),
	        ).Else(
	            Block(`fmt.Println("else")`),
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
```

Generate For blocks
```go
	f := File("hello.go").Package("main").Blocks(
	    Imports("fmt"),
	    Func("main()").Blocks(
	        ForBlock("").Blocks(
	            Block(`fmt.Println("hello")`),
	        ),
	        Block(`var strs = []string{"1", "2", "3"}`),
	        ForBlock("_, x := range strs").Blocks(
	            Block(`fmt.Println("hello", x)`),
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
```

Nil Code ignored
```go
	f := File("hello.go").Package("main").Blocks(
	    Struct("Hello").
	        AppendFieldComment("hello").
	        AppendField("Name", "string").
	        AppendFieldComment("").
	        AppendField("Age", "int"),
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
```

Consts example
```go
	f := File("hello.go").Package("main").Blocks(
	    ConstBlock().Type("Status", "string").Block(
	        `StatusError Status = "Error"`,
	    ).AppendConst("OK", "OK"),
	
	    ConstBlock().Type("HTTPStatus", "int").Block(
	        `HTTPStatusOK HTTPStatus = 200`,
	    ).AppendConst("Created", 201).
	        AppendConst("NotFound", 404),
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
```



