package txtbot

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Body struct {
	count int
	body  string
}

type AllBody []*Body

func Head(n int) string {
	a := "Hi Supplier,\n\n"
	b := "The following "
	c := " DRAFT ARTWORK files are for your FIRST approval:\n\n"
	if n == 1 {
		return a + "The following DRAFT ARTWORK file is for your FIRST approval:\n\n"
	} else {
		return a + b + strconv.Itoa(n) + c
	}
}

func SortTxt(path string) (int, *string, error) {
	fileinfo, _ := os.Stat(path)
	txtFile, err := os.Open(path)
	defer txtFile.Close()
	if err != nil {
		return 0, nil, err
	}
	// LineBreak in Mackintosh is CR(0xd), convert to linux LF(0xa).
	fileSize := fileinfo.Size()
	buf := make([]byte, fileSize)
	txtFile.Read(buf)
	for n, v := range buf {
		if v == 0xd {
			buf[n] = 0xa
		}
	}
	buftxt := bytes.NewBuffer(buf)

	sortedTxt := ""
	count := 0
	reader := bufio.NewScanner(buftxt)
	for reader.Scan() {
		sortedTxt = sortedTxt + strings.TrimSuffix(reader.Text(), ".pdf") + "\n"
		count++
	}
	return count, &sortedTxt, nil
}

func input() []byte {
	var inputLen int
	var inputer *bufio.Reader
	p := make([]byte, 12)
	for inputLen < 7 {
		fmt.Printf("Please input the job# no less than 6 digits: ")
		inputer = bufio.NewReader(os.Stdin)
		inputLen, _ = inputer.Read(p)
	}
	return p[0 : inputLen-1]
}

func FetchBody(pathes []string, bodylen int) (*AllBody, error) {
	var allbody AllBody
	allbody = make([]*Body, bodylen)
	for c, v := range pathes {
		linesNum, bdtxt, err := SortTxt(v)
		if err != nil {
			return nil, err
		}
		allbody[c] = &Body{count: linesNum, body: *bdtxt}
	}
	return &allbody, nil
}

func ConstructPDFName(DFjobpath string) (*AllBody, error) {
	count, files, err := SearchFile(DFjobpath, ".pdf")
	if err != nil {
		count, files, err = SearchFile(DFjobpath, ".PDF") /* not consider if ".PDF" and ".pdf" formats both exist. It rarely happens. */
	}
	if err != nil {
		return nil, err
	}
	var filenames string
	for _, value := range files {
		value = filepath.Base(value)
		if strings.HasSuffix(value, ".pdf") {
			filenames = filenames + strings.TrimSuffix(value, ".pdf") + "\n"
		} else {
			filenames = filenames + strings.TrimSuffix(value, ".PDF") + "\n"
		}
	}
	var allbody AllBody = make([]*Body, 1)
	allbody[0] = &Body{count: count, body: filenames}

	return &allbody, nil
}

func TitleSplit(path string) (string, string, string) {
	base := filepath.Base(path)
	tjob := "X"
	tcode := "X"
	if len(base) > 10 {
		base = strings.TrimSpace(base)
		splited := strings.Split(base, "_")
		if len(splited) == 2 {
			tjob = splited[0]
			tcode = splited[1]
		}
	}

	var separators = [...]string{"/", "-", "／"}
	if len(brand) > 4 {
		brand = strings.ToUpper(brand)
		brand = strings.TrimSpace(brand)
		before_sep := brand
		for _, separator := range separators {
			splitedbrand := strings.Split(brand, separator)
			if len(splitedbrand) == 2 { brand = splitedbrand[0] ; break}
		} 

		if before_sep == brand { brand = "X" }
	} else { brand = "X" }

	return brand, tcode, tjob
}
