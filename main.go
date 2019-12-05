package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func returnInum(file string) (inum string) {
	var stat syscall.Stat_t
	if err := syscall.Stat(file, &stat); err != nil {
		//fmt.Printf("%s: %v\n", file, err)
		inum = ""
		return inum
	}
	inum = strconv.FormatUint(stat.Ino, 10)
	return inum
}

func convertHexToDec(hexaNumber string) (dec int64) {
	dec, err := strconv.ParseInt(hexaNumber, 16, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	return dec
}

func walkTheFS(rootpath string) (allInums map[string]string) {
	allInums = make(map[string]string)
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		allInums[returnInum(path)] = path
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return allInums
}

func readFDINFO(filename string) (inums []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 9 {
			inumField := fields[2]
			hexInum := strings.Split(inumField, ":")[1]
			int64Inum := convertHexToDec(hexInum)
			inumString := strconv.FormatInt(int64Inum, 10)
			inums = append(inums, inumString)
		}
	}
	return inums
}

func huntInums(inumsToHunt []string, allInumsOnFS map[string]string) {
	for _, inum := range inumsToHunt {
		for fsInum, filename := range allInumsOnFS {
			if inum == fsInum {
				fmt.Printf("%s\n", filename)
				break
			}
		}
	}
}

func main() {
	fdinfo := flag.String("fdinfo", "", "Path to fdinfo file, example: /proc/123031/fdinfo/23")
	flag.Parse()
	inumsToHunt := readFDINFO(*fdinfo)
	fmt.Printf("Number of inotify to hunt down: %d\n", len(inumsToHunt))
	allInumsOnFS := walkTheFS("/")
	fmt.Printf("All inodes on filesystem: %d\n", len(allInumsOnFS))
	huntInums(inumsToHunt, allInumsOnFS)

}
