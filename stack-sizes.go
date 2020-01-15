// program stack-sizes reads out the stack sizes section of an elf compiled with the `-fstack-size-section` flag.
package main

import (
	"debug/elf"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"sort"
)

const STACK_SIZES_SECTION = ".stack_sizes"

var filename = flag.String("f", "", "Path to ELF file to analyze")

type StackSize struct {
	symbol uint64
	size   uint64
}

// decodeULEB128 decodes a string encoded in Unsigned-Little Endian base 128 format.
// See more in: https://en.wikipedia.org/wiki/LEB128
// This function assumes that the stack size fits in a uint64.
func decodeULEB128(r io.Reader) uint64 {
	result := uint64(0)
	shift := 0
	for {
		b := []byte{0x0}
		if _, err := r.Read(b); err != nil {
			log.Fatalf("failed to read: %v", err)
		}
		result |= uint64(b[0]&0x7F) << shift
		if b[0]&0x80 == 0x00 {
			break
		}
		shift += 7
	}

	return result
}

func main() {
	flag.Parse()
	if len(*filename) == 0 {
		log.Fatalf("Please specify a valid filename")
	}

	f, err := elf.Open(*filename)
	if err != nil {
		log.Fatalf("failed to open %q: %v", *filename, err)
	}
	defer f.Close()

	section := f.Section(STACK_SIZES_SECTION)
	if section == nil {
		log.Println("Missing Stack Sizes section.")
		return
	}

	// We need to translate the symbol values into their names for pretty-printing.
	symbols, err := f.Symbols()
	if err != nil {
		log.Fatalf("failed to parse object symbols")
	}
	symnames := make(map[uint64]string)
	for _, s := range symbols {
		symnames[s.Value] = s.Name
	}

	r := section.Open()
	sizes := make([]StackSize, 0)

	// The section contents consists of an array of [ symbol-value, stack-size ]
	// Where symbol-value is a 64bit number, and stack-size is a ULEB128 encoded value.
	for {
		// Read symbol value.
		var p uint64
		if err := binary.Read(r, binary.LittleEndian, &p); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to read: %v", err)
		}
		// Decode stack size.
		stackSize := decodeULEB128(r)
		sizes = append(sizes, StackSize{symbol: p, size: stackSize})
	}

	// Sort stack sizes in size-decreasing order.
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i].size > sizes[j].size
	})

	for _, v := range sizes {
		if _, ok := symnames[v.symbol]; !ok {
			log.Fatalf("could not find symbol name for %#x", v.symbol)
		}
		fmt.Printf("%#x\t%s: %d bytes\n", v.symbol, symnames[v.symbol], v.size)
	}
}
