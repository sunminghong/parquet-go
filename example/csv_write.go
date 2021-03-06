package main

import (
	"fmt"
	. "github.com/xitongsys/parquet-go/ParquetHandler"
	. "github.com/xitongsys/parquet-go/Plugin/CSVWriter"
	"log"
	"os"
)

type MyFile struct {
	FilePath string
	File     *os.File
}

func (self *MyFile) Create(name string) (ParquetFile, error) {
	file, err := os.Create(name)
	myFile := new(MyFile)
	myFile.File = file
	return myFile, err

}
func (self *MyFile) Open(name string) (ParquetFile, error) {
	var (
		err error
	)
	if name == "" {
		name = self.FilePath
	}

	myFile := new(MyFile)
	myFile.FilePath = name
	myFile.File, err = os.Open(name)
	return myFile, err
}

func (self *MyFile) Seek(offset int, pos int) (int64, error) {
	return self.File.Seek(int64(offset), pos)
}

func (self *MyFile) Read(b []byte) (n int, err error) {
	return self.File.Read(b)
}

func (self *MyFile) Write(b []byte) (n int, err error) {
	return self.File.Write(b)
}

func (self *MyFile) Close() {
	self.File.Close()
}

func main() {
	md := []MetadataType{
		{Type: "UTF8", Name: "Name"},
		{Type: "INT32", Name: "Age"},
		{Type: "INT64", Name: "Id"},
		{Type: "FLOAT", Name: "Weight"},
		{Type: "BOOLEAN", Name: "Sex"},
	}

	var f ParquetFile
	f = &MyFile{}

	//write flat
	f, _ = f.Create("csv.parquet")
	ph := NewCSVWriterHandler()
	ph.WriteInit(md, f, 10, 30)

	num := 10
	for i := 0; i < num; i++ {
		data := []string{
			"StudentName",
			fmt.Sprintf("%d", 20+i%5),
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%f", 50.0+float32(i)*0.1),
			fmt.Sprintf("%t", i%2 == 0),
		}
		rec := make([]*string, len(data))
		for j := 0; j < len(data); j++ {
			rec[j] = &data[j]
		}

		ph.Write(rec)
	}
	ph.WriteStop()
	log.Println("Write Finished")
	f.Close()

}
