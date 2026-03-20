# PUZ Parser

A PUZ file decoder and encoder based on the format specified [here](https://code.google.com/archive/p/puz/wikis/FileFormat.wiki).

## Features

- Encodes and Decodes PUZ Files
- Supports Extra Sections
- Unscrambles and Re-scrambles PUZ files
- Preserves all data

## Installation

```sh
go get github.com/cqb13/puz-parser
```

## Basic Usage

```go
import (
    puz "github.com/cqb13/puz-parser"
)

func main() {
    // get the bytes from a file

    puzzle, err := puz.DecodePuz(fileBytes)
    if err != nil {
		panic(err)
    }

    encodedBytes, err := puz.EncodePuz(puzzle)
    if err != nil {
		panic(err)
    }
}

```

## Acknowledgments

This project would not be possible without the help of the following:

- [PUZ File Format Wiki](https://code.google.com/archive/p/puz/wikis/FileFormat.wiki)
- [PUZ File Format ](https://web.archive.org/web/20151028113448/https://code.google.com/p/puz/wiki/FileFormat)
- [Cryptic Crossword](https://www.muppetlabs.com/~breadbox/txt/acre.html)
- [Scrambling Algorithm](https://www.muppetlabs.com/~breadbox/txt/scramble-c.txt)
- [PuzPy](https://github.com/alexdej/puzpy)
