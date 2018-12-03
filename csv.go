package csv

import (
	"bufio"
	//"encoding/binary"
	"encoding/json"
	"fmt"
	//"log"
	"go/types"
	"os"
	"strconv"
	//"reflect"
)

type RenderType int64

type Csv struct {
	Filename  string
	Delimiter string
	Data      []Row
	Render    RenderType
}

type Row struct {
	Fields []Field
}

type Field struct {
	Value    string
	DataType types.BasicKind
}

const (
	Normal   RenderType = 1
	Quotes   RenderType = 2
	NoQuotes RenderType = 3
)

func (c *Csv) WriteCsv() {
	file := CreateFile(c.Filename)
	defer file.Close() //make sure this closes

	w := bufio.NewWriter(file)

	for _, r := range c.Data {
		c.WriteRow(r, w)
	}
}

func (c *Csv) WriteRow(row Row, w *bufio.Writer) {
	count := len(row.Fields)

	for i, r := range row.Fields {
		c.WriteField(r, w)

		if i != count-1 {
			c.WriteDelimiter(w)
		}
	}
	WriteEnd(w)

	w.Flush()
}

func (c *Csv) WriteField(field Field, w *bufio.Writer) {

	if field.DataType == types.String {
		str := fmt.Sprintf("\"%s\"", field.Value)
		w.Write([]byte(str))
	}

	if field.DataType == types.Float64 {
		w.Write([]byte(field.Value))
	}
}

func (c *Csv) StructMap(i interface{}) {
	var mapped map[string]interface{}
	var fields []Field

	j, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(j, &mapped)

	for _, val := range mapped {
		var field Field
		FormatType(val, &field)
		fields = append(fields, field)
	}

	var row Row
	row.Fields = fields
	c.Data = append(c.Data, row)
}

func FormatType(value interface{}, field *Field) {
	switch v := value.(type) {
	case string:
		field.Value = v
		field.DataType = types.String
	case int64:
		str := strconv.FormatInt(v, 10)
		field.Value = str
		field.DataType = types.Int64
	case float64:
		str := strconv.FormatFloat(v, 'E', -1, 64)
		field.Value = str
		field.DataType = types.Float64
	case bool:
		//idk
	}
}

func (c *Csv) WriteDelimiter(w *bufio.Writer) {
	w.Write([]byte(c.Delimiter))
}

func WriteEnd(w *bufio.Writer) {
	end := "\r\n"
	w.Write([]byte(end))
}

//str := fmt.Sprintf("name type %v", t)

/*
func ReadCsv(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close() //make sure this closes

	csvFile := bufio.NewReader(f)
	//csvFile.Comma = 59 //semicolon

	//get all rows
	rows, err := csvFile.ReadAll()

	if err != nil {
		panic(err)
	}

	return rows
}
*/

func CreateFile(filename string) *os.File {
	file, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	return file
}
