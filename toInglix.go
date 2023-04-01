package main

import (
	byt "bytes"
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v2"
	osx "os"
	sts "strings"
	uni "unicode"
	utf "unicode/utf8"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz")

const dictionaryFile = "./dictionaries/English.bali"

func notInAlphabet(r rune) bool {
	return !byt.ContainsRune(alphabet, uni.ToLower(r))
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
	var dictionary = Dictionary(dictionaryFile)

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
		var word = string(english[index:next])

		// Translate the next word.
		var translation = dictionary.GetValue(sts.ToLower(word))
		if len(translation) == 0 {
			// Prompt for a new translation.
			fmt.Printf("Enter translation for %s: ", word)
			fmt.Scanln(&translation)
			if len(translation) > 0 {
				// Add a new word to the dictionary.
				dictionary.SetValue(sts.ToLower(word), sts.ToLower(translation))
			} else {
				// Keep the word untranslated.
				translation = word
			}
		}

		// Set the capitalization correctly.
		if uni.IsUpper(r) {
			translation = sts.Title(translation)
		}

		// Append the translated word to the Inglix text.
		buffer.WriteString(translation)
		index = next
	}

	// Write out the Inglix text.
	var inglix = buffer.Bytes()
	err = osx.WriteFile(osx.Args[2], inglix, 0644)
	if err != nil {
		panic(err)
	}

	// Write out the updated dictionary.
	dictionary.Save()
}

const (
	EOL = "\n" // The POSIX end of line character.
)

func Dictionary(file string) *dictionary {
	var v = col.Catalog[string, string]()
	var bytes, err = osx.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var lines = sts.Split(string(bytes), EOL)
	lines = lines[1:len(lines)-2]  // Remove the brackets.
	for _, line := range lines {
		var strings = sts.Split(line, `"`)
		v.SetValue(strings[1], strings[3])  // ----"key", "value"
	}
	return &dictionary{v, file}
}

type dictionary struct {
	col.CatalogLike[string, string]
	file string
}

func (v *dictionary) Save() {
	v.SortValues()
	var builder sts.Builder
	builder.WriteString("[" + EOL)
	var iterator = col.Iterator[col.Binding[string, string]](v)
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		builder.WriteString(`    "` + key + `": "` + value + `"` + EOL)
	}
	builder.WriteString("]" + EOL)
	var bytes = []byte(builder.String())
	var err = osx.WriteFile(v.file, bytes, 0644)
	if err != nil {
		panic(err)
	}
}
