package csvreadwrite

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	//CSV Data读取包
	"encoding/csv"
)

func parseStruct(str string) interface{} {
	//like [wall|5]
	var res []interface{}

	begPos := 0
	endPos := len(str)

	if strings.HasPrefix(str, "[") {
		begPos += 1
	}

	if strings.HasSuffix(str, "]") {
		endPos -= 1
	}

	if str == "" {
		return res
	}

	strs := strings.Split(str[begPos:endPos], "|")

	for _, v := range strs {

		res = append(res, v)

	}

	return res
}
func parseObject(typ string, str string) interface{} {
	isArray := false
	isStruct := false

	if strings.HasSuffix(typ, "[]") { // array
		isArray = true
	}

	if strings.HasPrefix(typ, "[") { //struct
		isStruct = true
		fmt.Println("isStruct = true")
	}
	//fmt.Println("1111",isArray)

	if isArray {
		var res []interface{}
		if isStruct {
			sep := "]|["
			if str != "" {
				eles := strings.Split(str, sep)
				for _, v := range eles {
					res = append(res, parseStruct(v))
				}
			}

		} else {
			sep := "|"
			if str != "" {
				eles := strings.Split(str, sep)
				for _, v := range eles {
					res = append(res, v)

				}
			}

		}

		return res

	} else {

		if isStruct {
			return parseStruct(str)
		} else {
			return str
		}

	}

}
func GetCSVData(path string) [][]interface{} {
	c, _ := ioutil.ReadFile(path)
	r := csv.NewReader(strings.NewReader(string(c)))
	rawcsvdata, err := r.ReadAll()

	if err != nil {
		fmt.Println(err)
	}
	if len(rawcsvdata) < 2 {
		fmt.Println(rawcsvdata)
		fmt.Println("Error: " + path + "is empty")
	}
	types := rawcsvdata[1]
	//fmt.Println(types)
	var rawdata [][]interface{}

	//dataStart := false
	for _, row := range rawcsvdata[1:len(rawcsvdata)] {

		var coldata []interface{}

		for i, v := range row {

			coldata = append(coldata, parseObject(types[i], v))
		}

		rawdata = append(rawdata, coldata)

	}

	return rawdata
}

func WriteCSVData(path string) bool {

	records := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
		return false
	}
	fout, err := os.Create(path)
	defer fout.Close()
	if err != nil {
		fmt.Println(path, err)
		return false
	}
	fout.WriteString(buf.String())

	return true
}

/*
func main() {

	successNomal := WriteCSVData("./xxx.csv")

	if successNomal {
		fmt.Println("写入成功")
	} else {
		fmt.Println("写入失败")
	}
	//fmt.Println(math.MaxUint64-1)
	data := GetCSVData("./xxx.csv")
	fmt.Println("file data raw length =", len(data))
	fmt.Println(data[0])
	//OutPut
	//写入成功
	//file data raw length = 3
	//[Rob Pike rob]

}
*/
