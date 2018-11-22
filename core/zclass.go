package core

import (
	"fmt"
	"io"
	"strings"
)

type ZClassProp struct {
	VarName   ZString
	Default   Runnable
	Modifiers ZObjectAttr
}

type ZClassMethod struct {
	Name      ZString
	Modifiers ZObjectAttr
	Method    Callable
}

type ZClass struct {
	Name ZString
	l    *Loc
	attr ZClassAttr

	Parent *ZClass

	// string value of extend & implement (used previous to lookup)
	ExtendsStr    ZString
	ImplementsStr []ZString

	Constructor *ZClassMethod

	Extends     *ZClass
	Implements  []*ZClass
	Const       map[ZString]Runnable // class constants
	Props       []*ZClassProp
	Methods     map[ZString]*ZClassMethod
	StaticProps *ZHashTable
}

func (c *ZClass) Run(ctx Context) (*ZVal, error) {
	// TODO resolve extendstr/implementsstr
	return nil, ctx.Global().RegisterClass(c.Name, c)
}

func (c *ZClass) Loc() *Loc {
	return c.l
}

func (c *ZClass) Dump(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%sclass %s {", c.attr, c.Name)
	if err != nil {
		return err
	}
	// TODO
	_, err = fmt.Fprintf(w, "TODO }")
	return err
}

func (c *ZClass) BaseName() ZString {
	// rturn class name without namespaces/etc
	pos := strings.LastIndexByte(string(c.Name), '\\')
	if pos == -1 {
		return c.Name
	}
	return c.Name[pos+1:]
}

func (c *ZClass) getStaticProps(ctx Context) (*ZHashTable, error) {
	if c.StaticProps == nil {
		c.StaticProps = NewHashTable()
		for _, p := range c.Props {
			if !p.Modifiers.IsStatic() {
				continue
			}
			if p.Default == nil {
				c.StaticProps.SetString(p.VarName, ZNULL.Dup())
				continue
			}
			v, err := p.Default.Run(ctx)
			if err != nil {
				c.StaticProps = nil
				return nil, err
			}
			c.StaticProps.SetString(p.VarName, v)
		}
	}
	return c.StaticProps, nil
}
