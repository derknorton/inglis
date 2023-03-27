package main

import (
	fmt "fmt"
	bal "github.com/bali-nebula/go-component-framework/v2/bali"
	osx "os"
	sts "strings"
	utf "unicode/utf8"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const dictionaryFile = "./dictionaries/English.bali"

func notInAlphabet(r rune) bool {
	return !sts.ContainsRune(alphabet, r)
}

func main() {
	// Validate the commandline arguments.
	if len(osx.Args) != 3 {
		panic("toInglix <english file> <inglix file>")
	}

	// Read in the English text.
	var bytes []byte
	var err error
	bytes, err = osx.ReadFile(osx.Args[1])
	if err != nil {
		panic(err)
	}
	var english = string(bytes)

	// Load in the dictionary.
	bytes, err = osx.ReadFile(dictionaryFile)
	if err != nil {
		panic(err)
	}
	var dictionary = bal.ParseDocument(bytes).ExtractCatalog()

	// Translate to Inglix text.
	var builder sts.Builder
	var index = 0
	var size = len(english)
	for index < size {
		var r, length = utf.DecodeRune([]byte(english[index:]))
		if notInAlphabet(r) {
			builder.WriteRune(r)
			index += length
			continue
		}
		var next = index + sts.IndexFunc(english[index:], notInAlphabet)
		var word = bal.Quote(`"` + english[index:next] + `"`)
		var translation = word
		var value = dictionary.GetValue(word)
		if value != nil {
			translation = value.ExtractQuote()
			if translation.IsEmpty() {
				fmt.Println("The dictionary is empty for:", word)
				translation = word
			}
		} else {
			fmt.Println("The dictionary is missing:", word)
		}
		builder.WriteString(translation.AsString())
		index = next
	}
	var inglix = builder.String()

	// Write out the Inglix text.
	err = osx.WriteFile(osx.Args[2], []byte(inglix), 0644)
	if err != nil {
		panic(err)
	}
}
