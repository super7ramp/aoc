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
	if len(d.files) > len(d.freeSpaces) {
		lastFile := d.files[len(d.files)-1]
		sb.WriteString(strings.Repeat(strconv.Itoa(lastFile.id), int(lastFile.size)))
	}
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

func (d *Disk) CompactFiles() {
	for fileToMoveIndex := len(d.files) - 1; fileToMoveIndex >= 0; fileToMoveIndex-- {
		fileToMove := d.files[fileToMoveIndex]
		for freeSpaceIndex := 0; freeSpaceIndex < fileToMoveIndex; freeSpaceIndex++ {
			freeSpace := d.freeSpaces[freeSpaceIndex]
			if freeSpace >= fileToMove.size {

				// Remove file from its current position and update space around the position
				d.files = slices.Delete(d.files, fileToMoveIndex, fileToMoveIndex+1)
				d.freeSpaces[fileToMoveIndex-1] += fileToMove.size
				if len(d.freeSpaces) > len(d.files) {
					// Merge contiguous free spaces
					d.freeSpaces[fileToMoveIndex-1] += d.freeSpaces[fileToMoveIndex]
					d.freeSpaces = slices.Delete(d.freeSpaces, fileToMoveIndex, fileToMoveIndex+1)
				}

				// Insert file in the free space
				d.files = slices.Insert(d.files, freeSpaceIndex+1, fileToMove)
				d.freeSpaces[freeSpaceIndex] = 0
				d.freeSpaces = slices.Insert(d.freeSpaces, freeSpaceIndex+1, freeSpace-fileToMove.size)

				fileToMoveIndex++
				//fmt.Println(d.String())
				break
			}
		}
	}
}

func (d *Disk) Checksum() int {
	checksum := 0
	blockPosition := 0
	for i, file := range d.files {
		for range file.size {
			checksum += blockPosition * file.id
			blockPosition += 1
		}
		if i < len(d.freeSpaces) {
			blockPosition += int(d.freeSpaces[i])
		}
	}
	return checksum
}

//go:embed input.txt
var input string

func main() {
	disk := ParseDisk(input)

	disk.Compact()
	fmt.Println("(Part 1) Compacted disk:", disk)
	fmt.Println("(Part 1) Checksum:", disk.Checksum())

	disk = ParseDisk(input)
	disk.CompactFiles()
	fmt.Println("(Part 2) Compacted disk:", disk)
	fmt.Println("(Part 2) Checksum:", disk.Checksum())
}
