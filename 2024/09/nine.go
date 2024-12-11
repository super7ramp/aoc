package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type File struct {
	id   int
	size byte
}

type Disk struct {
	files      []File
	freeSpaces []byte
}

func ParseDisk(input string) *Disk {
	files := make([]File, 0, len(input)/2+1)
	freeSpaces := make([]byte, 0, len(input)/2)
	for i, block := range []rune(input) {
		if i%2 == 0 {
			file := File{id: i / 2, size: byte(block - '0')}
			files = append(files, file)
		} else {
			freeSpaces = append(freeSpaces, byte(block-'0'))
		}
	}
	return &Disk{files, freeSpaces}
}

func (d *Disk) String() string {
	sb := strings.Builder{}
	for i := range d.freeSpaces {
		file := d.files[i]
		sb.WriteString(strings.Repeat(strconv.Itoa(file.id), int(file.size)))
		sb.WriteString(strings.Repeat(".", int(d.freeSpaces[i])))
	}
	lastFile := d.files[len(d.files)-1]
	sb.WriteString(strings.Repeat(strconv.Itoa(lastFile.id), int(lastFile.size)))
	return sb.String()
}

func (d *Disk) Compact() {
	for i := 0; i < len(d.freeSpaces); i++ {
		lastFile := d.files[len(d.files)-1]
		freeSpace := d.freeSpaces[i]
		if freeSpace > lastFile.size {
			// Move the entire last file to the free space
			d.freeSpaces[i] = 0
			d.freeSpaces = slices.Insert(d.freeSpaces, i+1, freeSpace-lastFile.size)
			d.files = d.files[:len(d.files)-1]
			d.files = slices.Insert(d.files, i+1, lastFile)
			d.freeSpaces = d.freeSpaces[:len(d.freeSpaces)-1]
		} else {
			// Move last file blocks that fit in the free space
			d.freeSpaces[i] = 0
			d.freeSpaces = slices.Insert(d.freeSpaces, i+1, 0)
			d.files = slices.Insert(d.files, i+1, File{id: lastFile.id, size: freeSpace})
			d.files[len(d.files)-1].size -= freeSpace
			i++
		}
		//fmt.Println(d.String())
	}
}

func (d *Disk) Checksum() int {
	checksum := 0
	blockPosition := 0
	for _, file := range d.files {
		for range file.size {
			checksum += blockPosition * file.id
			blockPosition += 1
		}
	}
	return checksum
}

//go:embed input-example.txt
var input string

func main() {
	disk := ParseDisk(input)
	fmt.Println(disk)
	disk.Compact()
	fmt.Println(disk)
	fmt.Println("Checksum", disk.Checksum())
}
