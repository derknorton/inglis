package main

import (
	byt "bytes"
	fmt "fmt"
	bal "github.com/bali-nebula/go-component-framework/v2/bali"
	osx "os"
	utf "unicode/utf8"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const dictionaryFile = "./dictionaries/English.bali"

func notInAlphabet(r rune) bool {
	return !byt.ContainsRune(alphabet, r)
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
	var english = bytes

	// Load in the dictionary.
	bytes, err = osx.ReadFile(dictionaryFile)
	if err != nil {
		panic(err)
	}
	var dictionary = bal.ParseDocument(bytes).ExtractCatalog()

	// Translate the English text.
	var buffer byt.Buffer
	var index = 0
	var size = len(english)
	for index < size {
		// Find the next word.
		var r, length = utf.DecodeRune(english[index:])
		if notInAlphabet(r) {
			// Append the non-letter rune to the Inglix text.
			buffer.WriteRune(r)
			index += length
			continue
		}

		// Extract the next word.
		var next = index + byt.IndexFunc(english[index:], notInAlphabet)
		var word = bal.Quote(`"` + string(english[index:next]) + `"`)
		var translation = word
		var value = dictionary.GetValue(word)

		// Translate the next word.
		if value != nil {
			translation = value.ExtractQuote()
			if translation.IsEmpty() {
				fmt.Println("The dictionary is empty for:", word)
				translation = word
			}
		} else {
			fmt.Println("The dictionary is missing:", word)
		}

		// Append the translated word to the Inglix text.
		buffer.WriteString(translation.AsString())
		index = next
	}

	// Write out the Inglix text.
	var inglix = buffer.Bytes()
	err = osx.WriteFile(osx.Args[2], inglix, 0644)
	if err != nil {
		panic(err)
	}
}
