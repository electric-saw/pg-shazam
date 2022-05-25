package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func genTest(kType string, value map[string]string, outFile string) error {
	file := getFileHeader()

	outFile = strings.ReplaceAll(outFile, ".go", "_test.go")
	outWriter, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err

	}

	varName := fmt.Sprintf("lex%sTestData", strings.Title(kType))

	dict := jen.Dict{}
	for k, v := range value {
		dict[jen.Id(k)] = jen.Lit(v)
	}

	file.Var().Id(varName).Op("=").Map(jen.String()).String().Values(dict)
	lexFunc := fmt.Sprintf("lex%s", cases.Title(language.English).String(kType))
	kindType := fmt.Sprintf("Kind%s", cases.Title(language.English).String(kType))
	funcName := fmt.Sprintf("TestLex%s", cases.Title(language.English).String(kType))

	file.Func().Id(funcName).Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
		jen.For(jen.List(jen.Id("tk"), jen.Id("s")).Op(":=").Range().Id(varName)).Block(
			jen.Id("token, _, found").Op(":=").Id(lexFunc).Call(jen.Id("s"), jen.Id("cursor{}")),

			jen.If(jen.Id("!found")).Block(
				jen.Id("t").Op(".").Id("Errorf").Call(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("token %s (%s) not found"), jen.Id("tk"), jen.Id("s")),
				),
			).Else().If(jen.Id("token").Dot("Kind").Op("!=").Id(kindType)).Block(

				jen.Id("t").Op(".").Id("Errorf").Call(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("token is %d, expected %d"), jen.Id("token.Kind"), jen.Id(kindType)),
				),
			).Else().If(jen.Id("token").Dot("Value").Op("!=").Id("tk")).Block(
				jen.Id("t").Op(".").Id("Errorf").Call(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("token value is %s, expected %s"), jen.Id("token.Value"), jen.Id("tk")),
				),
			),
		),
	)

	if err := file.Render(outWriter); err != nil {
		return err
	}

	return nil
}
