package txtbot

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	dfpath = "/Volumes/datavolumn_bmkserver_Pub"
	wks    = "/Volumes/datavolumn_bmkserver_Pub/新做稿/未开始"
	jxz    = "/Volumes/datavolumn_bmkserver_Pub/新做稿/进行中"
)

var (
	now       = time.Now()
	today     = now.Format("0102")
	yesterday = now.AddDate(0, 0, -1).Format("0102")
	month     = now.Format("200601")
	re        *regexp.Regexp
	job       string
	emailSave = filepath.Join(os.Getenv("HOME"), "Desktop", "draftartwork.txt")
	cancel    bool /*flag to skip executing bash command "open"*/
	help      bool /*flag, help*/
	trimTail  bool /*trim the tail*/
)

func init() {
	if len(os.Args) == 4 {
		if len(os.Args[3]) == 6 {
			job = os.Args[3]
		}
	}
	if len(os.Args) == 3 {
		if len(os.Args[2]) == 6 {
			job = os.Args[2]
		}
	}
	if len(os.Args) == 2 {
		if len(os.Args[1]) == 6 {
			job = os.Args[1]
		}
	}
	flag.BoolVar(&cancel, "c", false, "Cancel Open folder.")
	flag.BoolVar(&help, "h", false, "Help.")
	flag.BoolVar(&trimTail, "t", false, "Ignore PF mode, txt with out PF info.")
	flag.Usage = usage
}

func Run() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	if !TestConnect() {
		log.Println("Connection Errors, Please check if the server is connected at: " + dfpath)
		return
	}

	if len(job) == 0 {
		job = string(input())
	}
	job = strings.ToUpper(job)

	re = regexp.MustCompile(job)

	DFjobpath, err := FetchJobPath()
	if err != nil {
		log.Println(err)
		log.Println("Please check today and yesterday, if job folder is existed!")
		return
	}

	PFpath, err := FetchPFpath()
	if trimTail { /*shield the PF error if trimTail mode*/
		err = nil
	}
	if err != nil {
		log.Println(err)
		return
	}

	txtCount, txtFilePath, err := FetchTxtpath(DFjobpath)
	if err != nil {
		log.Println(err)
		if err.Error() == "nofile of .txt" {
			log.Println("No txt file in job folder, creating base on existing pdf files.")
		}
	}

	var allbody *AllBody
	switch len(txtFilePath) {
	case 0:
		allbody, err = ConstructPDFName(DFjobpath)
		if err != nil {
			log.Println(err)
			return
		}

	default:
		allbody, err = FetchBody(txtFilePath, txtCount)
		if err != nil {
			log.Println(err)
			return
		}
	}

	tail, err := FetchTail(PFpath)
	if err != nil {
		log.Println(err)
		return
	}

	emailTxt := CombineAll(allbody, &tail)
	WriteEmail(emailSave, emailTxt)

	log.Println("TxtSave", "------>", emailSave)

	if !cancel {
		cmd := exec.Command("open", DFjobpath)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func TestConnect() bool {
	result := Exists(dfpath)
	return result
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func SearchJob(path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	for _, file := range files {
		if file.IsDir() {
			if re.MatchString(file.Name()) {
				return filepath.Join(path, file.Name()), nil
			}
		}
	}
	return "", errors.New("no job of " + job + "on path: " + path)

}

func SearchFile(path, suf string) (int, []string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		return 0, nil, err
	}
	var found []string
	for _, file := range files {
		if !file.IsDir() && file.Name()[0] != '.' {
			if strings.HasSuffix(file.Name(), suf) {
				found = append(found, filepath.Join(path, file.Name()))
				continue
			}
		}
	}
	length := len(found)
	if length > 0 {
		return length, found, nil
	}
	return 0, nil, errors.New("nofile of " + suf)
}

func FetchJobPath() (string, error) {
	tpath := filepath.Join(dfpath, month, today)
	ypath := filepath.Join(dfpath, month, yesterday)
	tStatus := Exists(tpath)
	yStatus := Exists(ypath)

	if tStatus {
		if jobpath, err := SearchJob(tpath); err == nil {
			return jobpath, nil
		}
	}

	if yStatus {
		if jobpath, err := SearchJob(ypath); err == nil {
			return jobpath, nil
		}
	}

	return "", errors.New("Today and Yesterday, can not be located the job: " + job)
}

func FetchPFpath() (string, error) {
	jobpath, err := SearchJob(wks)
	if err != nil {
		jobpath, err = SearchJob(jxz)
	}
	if err != nil {
		return "", errors.New("PF job folder may not existed, please check!")
	}
	_, PFpath, err := SearchFile(jobpath, ".xls")
	if err != nil {
		_, PFpath, err = SearchFile(jobpath, ".xlsx")
	}
	if err != nil {
		return "", errors.New("PF sheet file is not located, please check!")
	}
	return PFpath[0], nil

}

func FetchTxtpath(jobpath string) (int, []string, error) {
	count, txtpath, err := SearchFile(jobpath, ".txt")
	if err != nil {
		return 0, nil, err
	}
	return count, txtpath, nil
}

func FetchTail(path string) (string, error) {
	if trimTail { /*flag value, trim tail mode*/
		return "", nil
	}
	if strings.HasSuffix(path, ".xls") {
		return ParseXls(path)
	} else {
		return ParseXlsx(path)
	}
}

func CombineAll(allbody *AllBody, tail *string) *string {
	emailtxt := ""
	for _, v := range *allbody {
		emailtxt = emailtxt + "\n" + Head(v.count) + v.body + "\n" + *tail + "\n"
	}
	return &emailtxt
}

func WriteEmail(path string, content *string) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}
	emailWriter := bufio.NewWriter(file)
	_, err = emailWriter.WriteString(*content) // Default 4096 byte.
	if err != nil {
		log.Println("write string err!")
		return
	}
	emailWriter.Flush()
	return
}

func usage() {
	fmt.Fprintf(os.Stderr, `txtbot version: txtbot/1.10.1
Usage: txtbot [-c job#]

Contact: bafelem@gmail.com
Project address: https://github.com/FengliangChen/txtbot

Options:
`)
	flag.PrintDefaults()
}
