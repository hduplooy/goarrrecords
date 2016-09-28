# hduplooy/goarrrecords

## Utility function to copy between string slices and records and records to records

Typical places where this can be used with is for example encoding/csv and github.com/hduplooy/gofixedwidth where slices of slices of strings are being generated or output. This can be used to map those slices to records and vice versa.

### Arr2Record(arr []string, rec interface{}) error

The slice of input strings are mapped to the fields within the provided structure (rec). An error is returned if there are either too many entries in the input are the input field could not be parsed according to the type of a field.

### Arr2Records(arr [][]string, itf interface{}) (interface{}, error)

The slice of slice of input strings are mapped to the fields within the provided structure (rec) per row and a slice based on the provided interface (itf) is returned. An error is returned if there are either too many entries in the input are the input field could not be parsed according to the type of a field.

### Record2Arr(rec interface{}) ([]string, error)

The fields of the provided record is converted to strings and returned as a slice of strings.

### Records2Arr(recs interface{}) ([][]string, error)

For each record of the provided slice the fields are converted to a slice of strings and all rows are returned as a slice of a slice of strings.

### CopyRec(dst interface{}, src interface{}) error

The fields in src are copied to the fields with the same name in the dst. These can be different structures but the fields with the same name must have the same type.

### CopyRecs(dstif interface{}, src interface{}) (interface{}, error)

For each record in the provided src slice a new record is created based on dstif. The fields are copied by name.

An example of use:


        package main

        import (
	        "encoding/csv"
	        "fmt"
	        "os"
	        "strings"

	        ar "github.com/hduplooy/goarrrecords"
	        fw "github.com/hduplooy/gofixedwidth"
        )

        type Person struct {
	        Name string
	        Id   int
        }

        type Employee struct {
	        Name   string
	        Id     int
	        Salary float64
        }

        func main() {
	        input := `This is a header line to be skipped
        # The following is info for the men
          John   1245
          Peter  5545
        # The following is info for certain women
          Susan  6784
          Sarah  4321
        `
	        sr := strings.NewReader(input)
	        r := fw.NewReader(sr)

	        r.Comment = '#'
	        r.SkipLines = 1
	        r.SkipStart = 2
	        r.TrimFields = true
	        r.FieldLengths = []int{7, 4}
	        tmp, _ := r.ReadAll()

	        persons, err := ar.Arr2Records(tmp, Person{})
	        if err != nil {
		        fmt.Printf("Arr2Records Error=%s\n", err.Error())
		        return
	        }
	        tmp2, err := ar.CopyRecs(Employee{}, persons)
	        if err != nil {
		        fmt.Printf("CopyRecs Error=%s\n", err.Error())
		        return
	        }
	        employees := tmp2.([]Employee)
	        for i, val := range employees {
		        if val.Id > 5000 {
			        employees[i].Salary = 15000
		        } else {
			        employees[i].Salary = 20000
		        }
	        }

	        tmp3, err := ar.Records2Arr(employees)
	        if err != nil {
		        fmt.Printf("Records2Arr Error=%s\n", err.Error())
		        return
	        }
	        w := csv.NewWriter(os.Stdout)
	        err = w.WriteAll(tmp3)
	        if err != nil {
		        fmt.Printf("Write Error=%s\n", err.Error())
	        }
        }

The output for this code is:

    John,1245,20000.000000
    Peter,5545,15000.000000
    Susan,6784,15000.000000
    Sarah,4321,20000.000000
    
