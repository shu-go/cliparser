package cliparser_test

import (
	"testing"

	"github.com/shu-go/cliparser"
	"github.com/shu-go/gotwant"
)

func TestParser(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, (*cliparser.Component)(nil))
	})

	t.Run("Option1", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
	})

	t.Run("Option2", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "-b"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "true"})
	})

	t.Run("OptionWithArg", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "hoge"})
		p.HintWithArg("a")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "hoge"})
	})

	t.Run("OptionWithArgErr", func(t *testing.T) {
		//p := cliparser.New([]string{"-a", "hoge"})
		p := cliparser.New()
		p.Feed([]string{"-a", "-b"})
		p.HintWithArg("a")

		err := p.Parse()
		gotwant.TestError(t, err, nil)
		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "-b"})
	})

	t.Run("OptionLong", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"--abc"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "abc", Arg: "true"})
	})

	t.Run("OptionLongHyphen", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"--a-b-c"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a-b-c", Arg: "true"})
	})

	t.Run("OptionShortLong", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-abc"})
		p.HintLongName("abc")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "abc", Arg: "true"})
	})

	t.Run("OptionShortConcat", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-abc"})
		//p.HintLongName("abc")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "c", Arg: "true"})
	})

	t.Run("OptionShortConcatArg", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-abc", "ccc"})
		p.HintWithArg("c")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "c", Arg: "ccc"})
	})

	t.Run("OptionShort=Long", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-abc"})
		//p.HintLongName("abc")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "c", Arg: "true"})

		//

		p.Reset()
		p.Feed([]string{"-abc", "-def"})
		p.HintNoOptionsGrouped()

		err = p.Parse()
		gotwant.TestError(t, err, nil)

		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "abc", Arg: "true"})
	})

	t.Run("SLSL", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "--bb", "-c", "--dd"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "bb", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "c", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "dd", Arg: "true"})
	})

	t.Run("cliparser.Command", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "sub", "-b"})
		p.HintCommand("sub")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "true"})
	})

	t.Run("Arg", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "sub", "-b"})
		//p.HintCommand("sub")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Arg: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "-b"})
	})

	t.Run("CommandArg", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "sub", "-b", "ccc", "ddd"})
		p.HintCommand("sub")
		p.HintWithArg("b", []string{"sub"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "ccc"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Arg: "ddd"})
	})

	t.Run("SubCommandArg", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "sub", "-b", "ccc", "subsub", "-d", "eee", "fff"})
		p.HintCommand("sub")
		p.HintCommand("subsub", []string{"sub"})
		p.HintWithArg("b", []string{"sub"})
		p.HintWithArg("d", []string{"sub", "subsub"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "ccc"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "subsub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "d", Arg: "eee"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Arg: "fff"})
	})

	t.Run("Command1", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"a"})
		p.HintCommand("a")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "a"})
	})

	t.Run("Arg1", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"a"})
		//p.HintCommand("a")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Arg: "a"})
	})

	t.Run("Namespace1", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "sub", "-b", "ccc", "ddd"})
		p.HintCommand("sub")
		p.HintWithArg("b", []string{"sub"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "ccc"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Arg: "ddd"})
	})

	t.Run("Namespace2", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"-a", "-b", "sub", "-b", "ccc", "ddd"})
		p.HintCommand("sub")
		p.HintWithArg("b", []string{"sub"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "b", Arg: "ccc"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Arg: "ddd"})
	})

	t.Run("NSAlias", func(t *testing.T) {
		p := cliparser.New()
		// Opt1 bool `cli:"go, opt1"`
		p.HintAlias("go", "opt1")
		p.HintLongName("go")
		p.HintLongName("opt1")
		// Cmd Cmd `cli:"c,cmd"`
		p.HintAlias("c", "cmd")
		p.HintCommand("c")
		p.HintCommand("cmd")
		// Cmd.Opt2 bool `cli:"opt2, co"`
		p.HintAlias("co", "opt2", []string{"cmd"})
		p.HintLongName("co", []string{"cmd"})
		p.HintLongName("opt2", []string{"cmd"})

		p.Feed([]string{"-go", "c", "-co"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "opt1", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "cmd", Arg: ""})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "opt2", Arg: "true"})
	})

	t.Run("Omit", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"--bool"})
		err := p.Parse()
		gotwant.TestError(t, err, nil)
		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "bool", Arg: "true"})

		p.Reset()
		p.Feed([]string{"--string"})
		p.HintWithArg("string")
		err = p.Parse()
		gotwant.TestError(t, err, "without argument")
		c = p.GetComponent()
		gotwant.Test(t, c, (*cliparser.Component)(nil))

		p.Reset()
		p.Feed([]string{"--string="})
		p.HintWithArg("string")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "string", Arg: ""})

		p.Reset()
		p.Feed([]string{"--string", "="})
		p.HintWithArg("string")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "string", Arg: ""})

		p.Reset()
		p.Feed([]string{"--string", "--hoge"})
		p.HintWithArg("string")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "string", Arg: "--hoge"})

		p.Reset()
		p.Feed([]string{"--string", "=", "--hoge"})
		p.HintWithArg("string")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "string", Arg: "--hoge"})

		p.Reset()
		p.Feed([]string{"--string", "sub"})
		p.HintWithArg("string")
		p.HintCommand("sub")
		err = p.Parse()
		gotwant.TestError(t, err, "without argument")
		c = p.GetComponent()
		gotwant.Test(t, c, (*cliparser.Component)(nil))

		p.Reset()
		p.Feed([]string{"--string", "=", "sub"})
		p.HintWithArg("string")
		p.HintCommand("sub")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "string", Arg: ""})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Command, Name: "sub"})
	})

	t.Run("OptCmdAfterArgs", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"--opt1", "arg1", "--opt2"})
		err := p.Parse()
		gotwant.TestError(t, err, nil)
		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "opt1", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "arg1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--opt2"})

		p.Reset()
		p.Feed([]string{"--string", "arg1?", "arg1", "sub"})
		p.HintWithArg("string")
		p.HintCommand("sub")
		err = p.Parse()
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "string", Arg: "arg1?"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "arg1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "sub"})
	})

	t.Run("DoubleDash", func(t *testing.T) {
		p := cliparser.New()
		p.Feed([]string{"--", "--opt1", "arg1", "--opt2"})
		p.HintWithArg("opt1")
		err := p.Parse()
		gotwant.TestError(t, err, nil)
		c := p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--opt1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "arg1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--opt2"})

		p.Reset()
		p.Feed([]string{"--opt1", "arg1", "--", "--opt2"})
		p.HintWithArg("opt1")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "opt1", Arg: "arg1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--opt2"})

		p.Reset()
		p.Feed([]string{"--opt1", "--", "arg1", "--opt2"})
		p.HintWithArg("opt1")
		err = p.Parse()
		gotwant.TestError(t, err, "without arguments")

		p.HintDisableDoubleHyphen()

		p.Reset()
		p.Feed([]string{"--", "--opt1", "arg1", "--opt2"})
		p.HintWithArg("opt1")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--opt1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "arg1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--opt2"})

		p.Reset()
		p.Feed([]string{"--opt1", "arg1", "--", "--opt2"})
		p.HintWithArg("opt1")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "opt1", Arg: "arg1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--opt2"})

		p.Reset()
		p.Feed([]string{"--opt1", "--", "arg1", "--opt2"})
		p.HintWithArg("opt1")
		err = p.Parse()
		gotwant.TestError(t, err, nil)
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Option, Name: "opt1", Arg: "--"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "arg1"})
		c = p.GetComponent()
		gotwant.Test(t, c, &cliparser.Component{Type: cliparser.Arg, Name: "", Arg: "--opt2"})
	})
}

func BenchmarkParse(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := cliparser.New()
		p.Feed([]string{"-a", "sub", "-b", "ccc", "subsub", "-d", "eee", "fff"})
		p.HintCommand("sub")
		p.HintCommand("subsub")
		p.HintWithArg("b")
		p.HintWithArg("d")

		p.Parse()
	}
}
