package cliparser

import (
	"testing"

	"bitbucket.org/shu/gotwant"
)

func TestParser(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		p := New()
		p.Feed([]string{})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, (*Component)(nil))
	})

	t.Run("Option1", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-a"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
	})

	t.Run("Option2", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-a", "-b"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "b", Arg: "true"})
	})

	t.Run("OptionWithArg", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-a", "hoge"})
		p.HintWithArg("a")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "hoge"})
	})

	t.Run("OptionWithArgErr", func(t *testing.T) {
		//p := New([]string{"-a", "hoge"})
		p := New()
		p.Feed([]string{"-a", "-b"})
		p.HintWithArg("a")

		err := p.Parse()
		gotwant.TestError(t, err, "without")
	})

	t.Run("OptionLong", func(t *testing.T) {
		p := New()
		p.Feed([]string{"--abc"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "abc", Arg: "true"})
	})

	t.Run("OptionShortLong", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-abc"})
		p.HintLongName("abc")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "abc", Arg: "true"})
	})

	t.Run("OptionShortConcat", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-abc"})
		//p.HintLongName("abc")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "b", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "c", Arg: "true"})
	})

	t.Run("OptionShortConcatArg", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-abc", "ccc"})
		p.HintWithArg("c")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "b", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "c", Arg: "ccc"})
	})

	t.Run("SLSL", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-a", "--bb", "-c", "--dd"})

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "bb", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "c", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "dd", Arg: "true"})
	})

	t.Run("Command", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-a", "sub", "-b"})
		p.HintCommand("sub")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Command, Name: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "b", Arg: "true"})
	})

	t.Run("Arg", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-a", "sub", "-b"})
		//p.HintCommand("sub")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Arg, Arg: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "b", Arg: "true"})
	})

	t.Run("CommandArg", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-a", "sub", "-b", "ccc", "ddd"})
		p.HintCommand("sub")
		p.HintWithArg("b")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Command, Name: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "b", Arg: "ccc"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Arg, Arg: "ddd"})
	})

	t.Run("SubCommandArg", func(t *testing.T) {
		p := New()
		p.Feed([]string{"-a", "sub", "-b", "ccc", "subsub", "-d", "eee", "fff"})
		p.HintCommand("sub")
		p.HintCommand("subsub")
		p.HintWithArg("b")
		p.HintWithArg("d")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "a", Arg: "true"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Command, Name: "sub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "b", Arg: "ccc"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Command, Name: "subsub"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Option, Name: "d", Arg: "eee"})
		c = p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Arg, Arg: "fff"})
	})

	t.Run("Command1", func(t *testing.T) {
		p := New()
		p.Feed([]string{"a"})
		p.HintCommand("a")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Command, Name: "a"})
	})

	t.Run("Arg1", func(t *testing.T) {
		p := New()
		p.Feed([]string{"a"})
		//p.HintCommand("a")

		err := p.Parse()
		gotwant.TestError(t, err, nil)

		c := p.GetComponent()
		gotwant.Test(t, c, &Component{Type: Arg, Arg: "a"})
	})
}

func BenchmarkParse(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := New()
		p.Feed([]string{"-a", "sub", "-b", "ccc", "subsub", "-d", "eee", "fff"})
		p.HintCommand("sub")
		p.HintCommand("subsub")
		p.HintWithArg("b")
		p.HintWithArg("d")

		p.Parse()
	}
}
