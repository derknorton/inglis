package main

import (
	buf "bufio"
	byt "bytes"
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v2"
	osx "os"
	sts "strings"
	uni "unicode"
	utf "unicode/utf8"
)

// MAIN PROGRAM

func main() {
	// Validate the commandline arguments.
	if len(osx.Args) < 3 {
		panic("toInglıś <english file> <ınglıś file> [learn]")
	}
	var learning = len(osx.Args) > 3

	// Read in the English text.
	var bytes []byte
	var err error
	bytes, err = osx.ReadFile(osx.Args[1])
	if err != nil {
		panic(err)
	}
	var english = bytes

	// Load in the dictionary.
	var dictionary = Dictionary(filename)

	// Translate the English text.
	var scanner = buf.NewScanner(osx.Stdin)
	var buffer byt.Buffer
	var index = 0
	var size = len(english)
	for index < size {
		// Find the next word.
		var r, length = utf.DecodeRune(english[index:])
		if notInAlphabet(r) {
			// Append the non-letter rune to the Inglıś text.
			buffer.WriteRune(r)
			index += length
			continue
		}

		// Extract the next word.
		var next = index + byt.IndexFunc(english[index:], notInAlphabet)
		var word = string(english[index:next])

		// Translate the next word.
		var translation = dictionary.GetValue(sts.ToLower(word))
		if learning && len(translation) == 0 {
			// Prompt for a new translation.
			fmt.Printf("Enter translation for %s: ", word)
			scanner.Scan()
			translation = scanner.Text()
			if len(translation) > 0 {
				// Add a new word to the dictionary.
				dictionary.SetValue(sts.ToLower(word), sts.ToLower(translation))
			}
		}
		if len(translation) == 0 {
			// Keep the word untranslated.
			translation = word
		}
		if uni.IsUpper(r) {
			translation = sts.Title(translation)
		}

		// Append the translated word to the Inglıś text.
		buffer.WriteString(translation)
		index = next
	}

	// Write out the Inglıś text.
	var ınglıś = buffer.Bytes()
	err = osx.WriteFile(osx.Args[2], ınglıś, 0644)
	if err != nil {
		panic(err)
	}

	// Write out the updated dictionary (if necessary).
	if learning {
		dictionary.Save()
	}
}

const (
	EOL      = "\n" // The POSIX end of line character.
	filename = "./dictionaries/English.txt"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz'")

func notInAlphabet(r rune) bool {
	return !byt.ContainsRune(alphabet, uni.ToLower(r))
}

// DICTIONARY IMPLEMENTATION

func Dictionary(file string) *dictionary {
	var v = col.Catalog[string, string]()
	var bytes, err = osx.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var lines = sts.Split(string(bytes), EOL)
	lines = lines[1 : len(lines)-2] // Remove the brackets.
	for _, line := range lines {
		var strings = sts.Split(line, `"`)
		v.SetValue(strings[1], strings[3]) // ----"key": "value"
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
