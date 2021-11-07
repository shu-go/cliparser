// Package cliparser provides a simple parsing API for CLI.
package cliparser

import (
	"fmt"
	"strings"
)

// ComponentType represents the type of parsed component.
type ComponentType int

const (
	// Unknown is not used in this package.
	Unknown ComponentType = iota
	// Option for -o -opt --opt
	Option
	// Command for command (requires call of Parser.HintCommand)
	Command
	// Arg is not an option nor a command.
	Arg
)

// Component is a resultant type of this package.
type Component struct {
	Type ComponentType

	Name string
	Arg  string
}

type hint struct {
	name      string
	namespace []string
}

func (t ComponentType) String() string {
	switch t {
	case Option:
		return "Option"
	case Command:
		return "Command"
	case Arg:
		return "Arg"
	default:
		return "Unknown"
	}
}

func (c Component) String() string {
	return fmt.Sprintf("Component{Type:%v, Name:%v, Arg:%v}", c.Type, c.Name, c.Arg)
}

// Parser contains parsing configurations and methods.
type Parser struct {
	args []string

	result []Component

	currNS           []string
	aliasHints       []hint
	commandHints     []hint
	withArgHints     []hint
	longNameHints    []hint
	optsMaybeGrouped bool
}

// New makes a Parser.
func New() Parser {
	return Parser{
		result:           make([]Component, 0, 8),
		aliasHints:       make([]hint, 0, 8),
		commandHints:     make([]hint, 0, 8),
		withArgHints:     make([]hint, 0, 8),
		longNameHints:    make([]hint, 0, 8),
		optsMaybeGrouped: true,
	}
}

// Reset resets its parsing results, except hints.
// Next, call Feed and Parse.
func (p *Parser) Reset() {
	p.args = p.args[:0]
	p.result = p.result[:0]
	p.currNS = p.currNS[:0]
}

// Feed is called when you pass os.Args.
// On next step, call Parser.Parse.
func (p *Parser) Feed(args []string) {
	p.args = args
}

// HintAlias is for defining another name.
func (p *Parser) HintAlias(alias, name string, optNS ...[]string) {
	if alias == name {
		return
	}

	h := hint{name: alias + ":" + name}
	if len(optNS) > 0 {
		h.namespace = optNS[0]
	}
	p.aliasHints = append(p.aliasHints, h)
}

// HintCommand is for giving the parser hint that the name is command.
func (p *Parser) HintCommand(name string, optNS ...[]string) {
	h := hint{name: name}
	if len(optNS) > 0 {
		h.namespace = optNS[0]
	}
	p.commandHints = append(p.commandHints, h)
}

// HintWithArg is for giving the parser hint that the name is option and it requires an argument.
func (p *Parser) HintWithArg(name string, optNS ...[]string) {
	h := hint{name: name}
	if len(optNS) > 0 {
		h.namespace = optNS[0]
	}
	p.withArgHints = append(p.withArgHints, h)
}

// HintLongName is for giving the parser hint that the name is option has a long name even if ONE-HYPHEND (-hoge)
func (p *Parser) HintLongName(name string, optNS ...[]string) {
	h := hint{name: name}
	if len(optNS) > 0 {
		h.namespace = optNS[0]
	}
	p.longNameHints = append(p.longNameHints, h)
}

// HintNoOptionsGrouped disallows -abc -> -a -b -c
func (p *Parser) HintNoOptionsGrouped() {
	p.optsMaybeGrouped = false
}

func (p Parser) toPhysicalName(alias string) string {
	for ai := 0; ai < len(p.aliasHints); ai++ {
		if strings.HasPrefix(p.aliasHints[ai].name, alias+":") && len(p.currNS) == len(p.aliasHints[ai].namespace) {
			for i := 0; i < len(p.aliasHints[ai].namespace); i++ {
				if p.currNS[i] != p.aliasHints[ai].namespace[i] {
					return alias
				}
			}
			return p.aliasHints[ai].name[len(alias)+1:]
		}
	}
	return alias
}

func (p Parser) testCommand(name string) bool {
	for ci := 0; ci < len(p.commandHints); ci++ {
		if p.commandHints[ci].name == name && len(p.currNS) == len(p.commandHints[ci].namespace) {
			ok := true
			for i := 0; i < len(p.commandHints[ci].namespace); i++ {
				if p.currNS[i] != p.commandHints[ci].namespace[i] {
					ok = false
				}
			}
			if ok {
				return true
			}
		}
	}
	return false
}

func (p Parser) testWithArg(name string) bool {
	for wi := 0; wi < len(p.withArgHints); wi++ {
		if p.withArgHints[wi].name == name && len(p.currNS) == len(p.withArgHints[wi].namespace) {
			ok := true
			for i := 0; i < len(p.withArgHints[wi].namespace); i++ {
				if p.currNS[i] != p.withArgHints[wi].namespace[i] {
					ok = false
				}
			}
			if ok {
				return true
			}
		}
	}
	return false
}

func (p Parser) testLongName(name string) bool {
	for li := 0; li < len(p.longNameHints); li++ {
		if p.longNameHints[li].name == name && len(p.currNS) == len(p.longNameHints[li].namespace) {
			ok := true
			for i := range p.longNameHints[li].namespace {
				if p.currNS[i] != p.longNameHints[li].namespace[i] {
					ok = false
				}
			}
			if ok {
				return true
			}
		}
	}
	return false
}

// GetComponent returns a Component. At end of source stream, this returns nil.
func (p *Parser) GetComponent() *Component {
	if len(p.result) == 0 {
		return nil
	}

	c := &(p.result[0])
	p.result = p.result[1:]

	return c
}

// Parse parses given (at Parser.Feed) command line string.
// Call Parser.GetComponent-s serially to get results.
func (p *Parser) Parse() error {
	var optName string
	var eqGiven bool
	var argsGiven bool

	// clear result
	p.result = p.result[:0]

	for {
		t, l := token(&p.args)
		if l == 0 {
			break
		}

		if argsGiven {
			p.result = append(p.result, Component{
				Type: Arg,
				Name: "",
				Arg:  t,
			})
			continue
		}

		// option?
		if (optName == "" || !p.testWithArg(optName)) && strings.HasPrefix(t, "-") {
			// first, process the prev option (because curr token is not an arg)
			if optName != "" {
				p.result = append(p.result, Component{
					Type: Option,
					Name: p.toPhysicalName(optName),
					Arg:  "true",
				})
			}

			// long name?
			if strings.HasPrefix(t, "--") {
				optName = t[2:]
				eqGiven = false
				continue
			}

			// long name or short-named options ?
			optName = t[1:]
			eqGiven = false
			if p.testLongName(optName) {
				continue
			}

			if p.optsMaybeGrouped {
				// short names (-abc -> -a -b -c)

				names := optName
				optName = ""
				eqGiven = false
				for ni := 0; ni < len(names); ni++ {
					if optName != "" {
						if p.testWithArg(optName) {
							return fmt.Errorf("option %q without arguments", optName)
						}
						p.result = append(p.result, Component{
							Type: Option,
							Name: p.toPhysicalName(optName),
							Arg:  "true",
						})
					}

					optName = string(names[ni])
					eqGiven = false
				}
			}
			continue

		} else if t == "=" {
			eqGiven = true

			if optName == "" {
				return fmt.Errorf("appeared = while no option given")
			} else if !p.testWithArg(optName) {
				return fmt.Errorf("option %q must not have an argument", optName)
			}
			continue

		} else {
			if optName != "" {
				if p.testWithArg(optName) {
					if p.testCommand(t) {
						if eqGiven {
							// first, process the prev option (because curr token is not an arg)
							p.result = append(p.result, Component{
								Type: Option,
								Name: p.toPhysicalName(optName),
								Arg:  "",
							})
							optName = ""
							eqGiven = false

							// command
							p.result = append(p.result, Component{
								Type: Command,
								Name: p.toPhysicalName(t),
							})
							p.currNS = append(p.currNS, p.toPhysicalName(t))

						} else {
							return fmt.Errorf("option %q without arguments", optName)
						}
					} else {

						// argument for an option
						p.result = append(p.result, Component{
							Type: Option,
							Name: p.toPhysicalName(optName),
							Arg:  t,
						})
						optName = ""
						eqGiven = false
					}

					continue

				} else {
					p.result = append(p.result, Component{
						Type: Option,
						Name: p.toPhysicalName(optName),
						Arg:  "true",
					})
					optName = ""
					eqGiven = false
				}
			}

			// command or args
			if p.testCommand(t) {
				p.result = append(p.result, Component{
					Type: Command,
					Name: p.toPhysicalName(t),
				})
				p.currNS = append(p.currNS, p.toPhysicalName(t))
			} else {
				p.result = append(p.result, Component{
					Type: Arg,
					Name: "",
					Arg:  t,
				})
				argsGiven = true
			}
		}
	}

	if optName != "" {
		if p.testWithArg(optName) {
			if eqGiven {
				p.result = append(p.result, Component{
					Type: Option,
					Name: p.toPhysicalName(optName),
					Arg:  "",
				})
			} else {
				return fmt.Errorf("option %q without arguments", optName)
			}
		} else {
			p.result = append(p.result, Component{
				Type: Option,
				Name: p.toPhysicalName(optName),
				Arg:  "true",
			})
			//optName = ""
			//eqGiven = false
		}
	}

	return nil
}

func token(args *[]string) (t string, length int) {
	if len(*args) == 0 {
		return "", 0
	}

	src := (*args)[0]
	if len(src) == 0 {
		return "", 0
	}

	switch src[0] {
	case '=':
		t, length = "=", 1

	case '"':
		for i := 1; i < len(src); i++ {
			if src[i] == '"' {
				t, length = src[1:i], i-2 //+ 1
				break
			}
		}
		if length == 0 { // centinel
			t, length = src[1:], len(src)-1
		}

	default:
		for i := 0; i < len(src); i++ {
			if src[i] == '=' {
				t, length = src[:i], i //+ 1
				break
			}
		}
		if length == 0 { // centinel
			t, length = src, len(src)
		}
	}

	// consume curr token on args
	(*args)[0] = (*args)[0][length:]
	if len((*args)[0]) == 0 {
		*args = (*args)[1:]
	}
	return t, length
}
