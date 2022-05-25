package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Node struct {
	Name     string
	Level    int
	TypeName string
	Children []*Node
}

func searchGen(wd string) {
	if err := generateFiles(wd); err != nil {
		panic(err)
	}
}

func generateFiles(outDir string) error {

	data, err := readDataJson(path.Join(outDir, "gen", "types.json"))
	if err != nil {
		return err
	}

	for kType, value := range data {
		fmt.Printf("Generating %s\n", kType)
		capType := cases.Title(language.English).String(kType)

		file := getFileHeader()
		generateFunc(file, "lex"+capType, kType, "Kind"+capType, value)

		fileName := fmt.Sprintf("gen_lex_%s.go", kType)

		fmt.Printf("Writing to file %s\n", fileName)

		outFile := path.Join(outDir, fileName)
		outWriter, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
		if err != nil {
			return err

		}

		if err := file.Render(outWriter); err != nil {
			return err
		}

		if err := genTest(kType, value, outFile); err != nil {
			return err
		}

	}
	return nil
}

func readDataJson(jsonPath string) (map[string]map[string]string, error) {
	f, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make(map[string]map[string]string)

	err = json.NewDecoder(f).Decode(&data)
	return data, err
}

// generate best way to search for a token from string
func generateFunc(file *jen.File, name, goType, kind string, data map[string]string) {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	consts := []jen.Code{}
	for _, k := range keys {
		consts = append(consts, jen.Id(k).Op("=").Lit(data[k]))
	}

	file.Const().Defs(consts...)

	const (
		source    = "source"
		srcCursor = "srcCursor"
		cursor    = "cur"
	)

	f := file.Func().Id(name).Params(
		jen.Id(source).String(),
		jen.Id(srcCursor).Qual("internal/pkg/parser/lexer", "cursor"),
	).Params(
		jen.Qual("internal/pkg/parser/lexer", "*Token"),
		jen.Qual("internal/pkg/parser/lexer", "cursor"),
		jen.Bool(),
	)

	tree := generatePrefixTree(data)

	f.BlockFunc(func(g *jen.Group) {
		g.Var().Id(cursor).Qual("internal/pkg/parser/lexer", "cursor").Op("=").Id(srcCursor)

		tree.generateSearchCode(g, kind, source, srcCursor, "cur")

	})
}

func generatePrefixTree(data map[string]string) *Node {
	root := &Node{
		Level: 0,
	}
	for k, v := range data {
		cur := root
		for i, c := range v {
			child := cur.findChild(string(c))
			if child == nil {
				child = &Node{
					Name:  string(c),
					Level: cur.Level + 1,
				}
				if i == len(v)-1 {
					child.TypeName = k
					// fmt.Printf("%s => %s\n", v, child.TypeName)
				}

				cur.Children = append(cur.Children, child)
			}
			cur = child
		}
	}

	root.fixEmpty()
	root.printTree(0)

	return root
}

func (n *Node) findChild(name string) *Node {
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (n *Node) printTree(indent int) {
	fmt.Printf("%s%s => %s\n", strings.Repeat(" ", indent), n.Name, n.TypeName)
	for _, c := range n.Children {
		c.printTree(indent + 2)
	}
}

func (n *Node) fixEmpty() {
	if len(n.Children) == 0 {
		n.Children = nil
	}

	for _, c := range n.Children {
		c.fixEmpty()
	}

	if len(n.Children) == 1 && n.TypeName == "" {
		if len(n.Children[0].Children) == 0 {
			n.Name += n.Children[0].Name
			n.TypeName = n.Children[0].TypeName
			n.Children = nil
		}
	}

}

func (n *Node) generateSearchCode(jenGroup *jen.Group, kind, source, srcCur, cur string) {
	var runeChilren []*Node
	var stringChildren []*Node

	for _, c := range n.Children {
		if len(c.Name) == 1 {
			runeChilren = append(runeChilren, c)
		} else {
			stringChildren = append(stringChildren, c)
		}
	}

	switch {
	case len(n.Name) == 0:
		for _, c := range stringChildren {
			c.generateSearchCode(jenGroup, kind, source, srcCur, cur)
		}

		jenGroup.Switch(jen.Id(source).Index(jen.Id(cur).Dot("pointer"))).BlockFunc(func(s *jen.Group) {
			for _, c := range runeChilren {
				c.generateSearchCode(s, kind, source, srcCur, cur)
			}
		})

		jenGroup.Return(
			jen.Nil(),
			jen.Id(srcCur),
			jen.False(),
		)

	case len(n.Name) == 1:
		rLower := unicode.ToLower(rune(n.Name[0]))
		rUpper := unicode.ToUpper(rune(n.Name[0]))

		cases := []jen.Code{jen.LitRune(rune(rLower))}
		if rLower != rUpper {
			cases = append(cases, jen.LitRune(rune(rUpper)))
		}

		jenGroup.Case(cases...).BlockFunc(func(s *jen.Group) {
			s.Id(cur).Dot("pointer").Op("++")
			s.Id(cur).Dot("loc").Dot("Col").Op("++")

			if len(stringChildren) > 0 || len(runeChilren) > 0 {
				ifLen := jen.If(jen.Len(jen.Id(source).Index(jen.Id(cur).Dot("pointer").Id(":"))).Op(">=").Lit(len(n.Name))).Comment("/* sub group */").BlockFunc(func(ifGroup *jen.Group) {
					if len(runeChilren) > 0 {
						ifGroup.Switch(jen.Id(source).Index(jen.Id(cur).Dot("pointer"))).BlockFunc(func(s *jen.Group) {
							for _, c := range runeChilren {
								c.generateSearchCode(s, kind, source, srcCur, cur)
							}
						})
					}

					if len(stringChildren) > 0 {
						for _, c := range stringChildren {
							c.generateSearchCode(ifGroup, kind, source, srcCur, cur)
						}
					}
				})
				s.Line().Add(ifLen).Line()

			}

			if n.TypeName != "" {
				s.Return(
					jen.Qual("internal/pkg/parser/lexer", "&Token").Values(
						jen.Id("Kind").Op(":").Id(kind),
						jen.Id("Value").Op(":").String().Call(jen.Id(n.TypeName)),
						jen.Id("Loc").Op(":").Id(srcCur).Dot("loc"),
					),
					jen.Id(cur),
					jen.True(),
				)
			}

		})
	default:
		ifLen := jen.If(jen.Len(jen.Id(source).Index(jen.Id(cur).Dot("pointer").Id(":"))).Op(">=").Lit(len(n.Name))).BlockFunc(func(g *jen.Group) {
			var1 := jen.Id(source).Index(jen.Id(cur).Dot("pointer").Id(":").Id(cur).Dot("pointer").Op("+").Lit(len(n.Name)))
			var2 := jen.Lit(n.Name)

			compare := jen.Qual("strings", "EqualFold").Call(var1, var2)

			g.If(compare).BlockFunc(func(s *jen.Group) {
				s.Id(cur).Dot("pointer").Op("+=").Lit(len(n.Name))
				s.Id(cur).Dot("loc").Dot("Col").Op("+=").Lit(len(n.Name))
				if n.TypeName != "" {
					s.Return(
						jen.Qual("internal/pkg/parser/lexer", "&Token").Values(
							jen.Id("Kind").Op(":").Id(kind),
							jen.Id("Value").Op(":").String().Call(jen.Id(n.TypeName)),
							jen.Id("Loc").Op(":").Id(srcCur).Dot("loc"),
						),
						jen.Id(cur),
						jen.True(),
					)
				} else {

					if len(runeChilren) > 0 {
						s.If(jen.Len(jen.Id(source).Index(jen.Id(cur).Dot("pointer").Id(":"))).Op(">=").Lit(len(n.Name))).BlockFunc(func(ifGroup *jen.Group) {
							ifGroup.Switch(jen.Id(source).Index(jen.Id(cur).Dot("pointer"))).BlockFunc(func(s *jen.Group) {
								for _, c := range runeChilren {
									c.generateSearchCode(s, kind, source, srcCur, cur)
								}
							})
						})
					}

					if len(stringChildren) > 0 {

						for _, c := range stringChildren {
							ifLen := jen.If(jen.Len(jen.Id(source).Index(jen.Id(cur).Dot("pointer").Id(":"))).Op(">").Lit(len(n.Name))).BlockFunc(func(g *jen.Group) {
								c.generateSearchCode(g, kind, source, srcCur, cur)
							})
							s.Line().Add(ifLen).Line()
						}
					}
				}
			},
			)

		})
		jenGroup.Line().Add(ifLen).Line()

	}
}
