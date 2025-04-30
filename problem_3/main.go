/*
Create a program that demonstrates RAID data handling for RAID0, RAID1,
RAID10, RAID5, and RAID6.
a. Define the data type for RAID data storage, for example, an array(raid) of
array(disk) of byte of array(data stripe).
b. Write a string([]byte) into the RAID at position 0 with a length N, where N
should be greater than the size of one stripe.
c. Clear one of the disks in the RAID, setting it to zero.
d. Read the data in the RAID from 0 to N, convert it back to a string, and print it.
*/

package main

import (
	"fmt"
	"os"
)

func main() {
	content, err := os.ReadFile("inputSample.txt")
	if err != nil {
		fmt.Println("Unable to read file:", err)
		return
	}

	text := string(content)
	fmt.Println("File contents:\n", text)

	stripeSize := 12

	raid0Demo(content, stripeSize)
	raid1Demo(content, stripeSize)
	raid10Demo(content, stripeSize)
	raid5Demo(content, stripeSize)
	raid6Demo(content, stripeSize)
}

func raid0Demo(content []byte, stripeSize int) {
	numDisks := 4
	Disk := make([][][]byte, numDisks)
	for i := range Disk {
		Disk[i] = make([][]byte, 0)
	}
	fmt.Println("This is a simulation of RAID 0")

	for i := 0; i < len(content); i += stripeSize {
		end := i + stripeSize
		if end > len(content) {
			end = len(content)
		}
		stripe := content[i:end]

		// Use round-robin to distribute stripes across disks
		disknum := (i / stripeSize) % numDisks
		Disk[disknum] = append(Disk[disknum], stripe)
	}

	data := raid0Read(Disk, numDisks)
	fmt.Println("Rebuilt content:")
	fmt.Println(string(data))

	Disk[1] = nil

	brokenData := raid0Read(Disk, numDisks)
	fmt.Println("Rebuilt Broken data content:")
	fmt.Println(string(brokenData))
}

func raid0Read(Disk [][][]byte, numDisks int) []byte {
	var result []byte

	maxStripes := 0
	// Find the maximum number of stripes
	for _, disk := range Disk {
		if len(disk) > maxStripes {
			maxStripes = len(disk)
		}
	}

	for stripeIndex := 0; stripeIndex < maxStripes; stripeIndex++ {
		for diskIndex := 0; diskIndex < numDisks; diskIndex++ {
			if stripeIndex < len(Disk[diskIndex]) {
				result = append(result, Disk[diskIndex][stripeIndex]...)
			}
		}
	}

	return result
}

func raid1Demo(content []byte, stripeSize int) {
	numDisks := 2
	Disk := make([][][]byte, numDisks)
	for i := range Disk {
		Disk[i] = make([][]byte, 0)
	}
	fmt.Println("This is a simulation of RAID 1")

	for i := 0; i < len(content); i += stripeSize {
		end := i + stripeSize
		if end > len(content) {
			end = len(content)
		}
		stripe := content[i:end]

		disknum := (i / stripeSize) % (numDisks / 2)
		mirror := disknum + (numDisks / 2)

		Disk[disknum] = append(Disk[disknum], stripe)
		Disk[mirror] = append(Disk[mirror], stripe)
	}

	data := raid1Read(Disk, numDisks)
	fmt.Println("Rebuilt content:")
	fmt.Println(string(data))

	Disk[1] = nil

	brokenData := raid1Read(Disk, numDisks)
	fmt.Println("Rebuilt Broken data content:")
	fmt.Println(string(brokenData))
}

func raid1Read(Disk [][][]byte, numDisks int) []byte {
	var result []byte

	maxStripes := 0
	// Find the maximum number of stripes
	for _, disk := range Disk {
		if len(disk) > maxStripes {
			maxStripes = len(disk)
		}
	}
	half := numDisks / 2
	damageInspection := false

	for i := 0; i < half; i++ {
		if Disk[i] == nil {
			damageInspection = true
		}
	}

	if !damageInspection {
		for stripeIndex := 0; stripeIndex < maxStripes; stripeIndex++ {
			for diskIndex := 0; diskIndex < half; diskIndex++ {
				if stripeIndex < len(Disk[diskIndex]) {
					result = append(result, Disk[diskIndex][stripeIndex]...)
				}
			}
		}
	} else {
		for stripeIndex := 0; stripeIndex < maxStripes; stripeIndex++ {
			for diskIndex := 0 + half; diskIndex < numDisks; diskIndex++ {
				if stripeIndex < len(Disk[diskIndex]) {
					result = append(result, Disk[diskIndex][stripeIndex]...)
				}
			}
		}
	}

	return result
}

func raid10Demo(content []byte, stripeSize int) {
	numDisks := 4
	Disk := make([][][]byte, numDisks)
	for i := range Disk {
		Disk[i] = make([][]byte, 0)
	}
	fmt.Println("This is a simulation of RAID 10")

	for i := 0; i < len(content); i += stripeSize {
		end := i + stripeSize
		if end > len(content) {
			end = len(content)
		}
		stripe := content[i:end]

		disknum := (i / stripeSize) % (numDisks / 2)
		mirror := disknum + (numDisks / 2)

		Disk[disknum] = append(Disk[disknum], stripe)
		Disk[mirror] = append(Disk[mirror], stripe)
	}

	data := raid10Read(Disk, numDisks)
	fmt.Println("Rebuilt content:")
	fmt.Println(string(data))

	Disk[1] = nil

	brokenData := raid10Read(Disk, numDisks)
	fmt.Println("Rebuilt Broken data content:")
	fmt.Println(string(brokenData))
}

func raid10Read(Disk [][][]byte, numDisks int) []byte {
	var result []byte

	maxStripes := 0
	// Find the maximum number of stripes
	for _, disk := range Disk {
		if len(disk) > maxStripes {
			maxStripes = len(disk)
		}
	}
	half := numDisks / 2
	for stripeIndex := 0; stripeIndex < maxStripes; stripeIndex++ {
		for diskIndex := 0; diskIndex < half; diskIndex++ {
			if Disk[diskIndex] == nil {
				if stripeIndex < len(Disk[diskIndex+half]) {
					result = append(result, Disk[diskIndex+half][stripeIndex]...)
				}
			} else {
				if stripeIndex < len(Disk[diskIndex]) {
					result = append(result, Disk[diskIndex][stripeIndex]...)
				}
			}
		}
	}
	return result // Did not check the situation where both disks are broken
}

func raid5Demo(content []byte, stripeSize int) {
	numDisks := 3
	Disk := make([][][]byte, numDisks)
	for i := range Disk {
		Disk[i] = make([][]byte, 0)
	}
	fmt.Println("This is a simulation of RAID 5")
	stripes := [][]byte{}

	for i := 0; i < len(content); i += stripeSize {
		end := i + stripeSize
		if end > len(content) {
			end = len(content)
		}
		stripe := make([]byte, stripeSize)
		copy(stripe, content[i:end])
		stripes = append(stripes, stripe)

	}
	count := 0
	for i := 0; i < len(stripes)-1; i += 2 {
		stripe1 := stripes[i]
		stripe2 := stripes[i+1]
		parity := make([]byte, stripeSize)
		for j := 0; j < stripeSize; j++ {
			parity[j] = stripe1[j] ^ stripe2[j]
		} // Calculate XOR parity for two data stripes

		parityDisk := count % numDisks
		count++
		dataDisks := []int{}
		for k := 0; k < numDisks; k++ {
			if k != parityDisk {
				dataDisks = append(dataDisks, k)
			}
		}

		Disk[parityDisk] = append(Disk[parityDisk], parity)
		Disk[dataDisks[0]] = append(Disk[dataDisks[0]], stripe1)
		Disk[dataDisks[1]] = append(Disk[dataDisks[1]], stripe2)
	}

	data := raid5Read(Disk, numDisks)
	fmt.Println("Rebuilt content:")
	fmt.Println(string(data))

	Disk[1] = nil

	brokenData := raid5Read(Disk, numDisks)
	fmt.Println("Rebuilt Broken data content:")
	fmt.Println(string(brokenData))
}

func raid5Read(Disk [][][]byte, numDisks int) []byte {
	var result []byte

	maxStripes := 0
	// Find the maximum number of stripes
	for _, disk := range Disk {
		if len(disk) > maxStripes {
			maxStripes = len(disk)
		}
	}

	count := 0
	for i := 0; i < maxStripes; i++ {
		parityDisk := count % 3
		count++
		dataDisks := []int{}
		for k := 0; k < numDisks; k++ {
			if k != parityDisk {
				dataDisks = append(dataDisks, k)
			}
		}
		var stripe0, stripe1, parity []byte
		if Disk[dataDisks[0]] != nil && i < len(Disk[dataDisks[0]]) {
			stripe0 = Disk[dataDisks[0]][i]
		}
		if Disk[dataDisks[1]] != nil && i < len(Disk[dataDisks[1]]) {
			stripe1 = Disk[dataDisks[1]][i]
		}
		if Disk[parityDisk] != nil && i < len(Disk[parityDisk]) {
			parity = Disk[parityDisk][i]
		}

		if Disk[dataDisks[0]] == nil {
			stripe0 = make([]byte, len(parity))
			for j := 0; j < len(stripe1); j++ {
				stripe0[j] = parity[j] ^ stripe1[j]
			}
		}

		if Disk[dataDisks[1]] == nil {
			stripe1 = make([]byte, len(parity))
			for j := 0; j < len(stripe0); j++ {
				stripe1[j] = parity[j] ^ stripe0[j]
			}
		}

		result = append(result, stripe0...)
		result = append(result, stripe1...)
	}
	return result
}

func raid6Demo(content []byte, stripeSize int) {
	numDisks := 4
	Disk := make([][][]byte, numDisks)
	for i := range Disk {
		Disk[i] = make([][]byte, 0)
	}
	fmt.Println("This is a simulation of RAID 6")
	stripes := [][]byte{}

	for i := 0; i < len(content); i += stripeSize {
		end := i + stripeSize
		if end > len(content) {
			end = len(content)
		}
		stripe := make([]byte, stripeSize)
		copy(stripe, content[i:end])
		stripes = append(stripes, stripe)

	}
	count := 0

	for i := 0; i < len(stripes)-1; i += 2 {
		stripe1 := stripes[i]
		stripe2 := stripes[i+1]
		parity := make([]byte, stripeSize)
		qParity := make([]byte, stripeSize)
		for j := 0; j < stripeSize; j++ {
			parity[j] = stripe1[j] ^ stripe2[j]
		}
		for j := 0; j < stripeSize; j++ {
			qParity[j] = stripe1[j]<<1 ^ stripe2[j]<<1
		} // Calculate Q parity with simple linear transformation (for simulation only)

		parityDisk := count % numDisks
		qParityDisk := (parityDisk + 1) % numDisks
		count++
		dataDisks := []int{}
		for k := 0; k < numDisks; k++ {
			if k != parityDisk && k != qParityDisk {
				dataDisks = append(dataDisks, k)
			}
		}

		Disk[parityDisk] = append(Disk[parityDisk], parity)
		Disk[qParityDisk] = append(Disk[qParityDisk], qParity)
		Disk[dataDisks[0]] = append(Disk[dataDisks[0]], stripe1)
		Disk[dataDisks[1]] = append(Disk[dataDisks[1]], stripe2)
	}

	data := raid6Read(Disk, numDisks)
	fmt.Println("Rebuilt content:")
	fmt.Println(string(data))

	Disk[1] = nil

	brokenData := raid6Read(Disk, numDisks)
	fmt.Println("Rebuilt Broken data content:")
	fmt.Println(string(brokenData))
}

func raid6Read(Disk [][][]byte, numDisks int) []byte {
	var result []byte

	maxStripes := 0
	// Find the maximum number of stripes
	for _, disk := range Disk {
		if len(disk) > maxStripes {
			maxStripes = len(disk)
		}
	}

	count := 0
	for i := 0; i < maxStripes; i++ {
		parityDisk := count % 4
		qParityDisk := (parityDisk + 1) % 4
		count++
		dataDisks := []int{}
		for k := 0; k < numDisks; k++ {
			if k != parityDisk && k != qParityDisk {
				dataDisks = append(dataDisks, k)
			}
		}
		var stripe0, stripe1, parity, qParity []byte
		if Disk[dataDisks[0]] != nil && i < len(Disk[dataDisks[0]]) {
			stripe0 = Disk[dataDisks[0]][i]
		}
		if Disk[dataDisks[1]] != nil && i < len(Disk[dataDisks[1]]) {
			stripe1 = Disk[dataDisks[1]][i]
		}
		if Disk[parityDisk] != nil && i < len(Disk[parityDisk]) {
			parity = Disk[parityDisk][i]
		}
		if Disk[qParityDisk] != nil && i < len(Disk[qParityDisk]) {
			qParity = Disk[qParityDisk][i]
		}

		if Disk[dataDisks[0]] == nil {
			stripe0 = make([]byte, len(parity))
			for j := 0; j < len(stripe1); j++ {
				stripe0[j] = parity[j] ^ stripe1[j]
			}
		}

		if Disk[dataDisks[1]] == nil {
			stripe1 = make([]byte, len(parity))
			for j := 0; j < len(stripe0); j++ {
				stripe1[j] = parity[j] ^ stripe0[j]
			}
		}
		_ = qParity // Q parity is currently unused

		result = append(result, stripe0...)
		result = append(result, stripe1...)
	}
	return result
}
