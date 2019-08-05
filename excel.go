package txtbot

import (
	"errors"
	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

func ParseXls(path string) (string, error) {
	if xlFile, err := xls.Open(path, "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			brand = sheet1.Row(2).Col(1)
			xlsContent := ""
			for i := 4; i < 12; i++ {
				cell := sheet1.Row(i).Col(7)
				xlsContent = xlsContent + cell + "\n"
			}
			return xlsContent, nil
		}
	}
	return "", errors.New("Parse xls error.")
}

func ParseXlsx(path string) (string, error) {
	if xlFile, err := xlsx.OpenFile(path); err == nil {
		brand = xlFile.Sheets[0].Rows[2].Cells[1].Value
		xlsContent := ""
		for i := 4; i < 12; i++ {
			xlsContent = xlsContent + xlFile.Sheets[0].Rows[i].Cells[7].Value + "\n"
		}
		return xlsContent, nil
	}
	return "", errors.New("Parse xlsx error.")
}
