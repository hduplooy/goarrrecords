// github.com/hduplooy/goarrrecords
// Author: Hannes du Plooy
// Revision Date: 28 Sep 2016
// This is a utility library to convert slice(s) of strings to a structure (or slice of structures),
// or from a structure (or slice of structures) to a slice(s) of strings
// Furthermore based on the field names data can be copied between different structures
package goarrrecords

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Arr2Record will take a slice of strings and covert them to the fields of the provided structure
// Only int, float64 and string is accomodated for at the moment
// Errors occur when either the provided interface is not a structure, or more entries are in the
// provided input than the fields in the record or fields cannot be parsed to the types of the fields of the record
func Arr2Record(arr []string, rec interface{}) error {
	valr := reflect.Indirect(reflect.ValueOf(rec))
	if valr.Kind() != reflect.Struct {
		return errors.New("the provided record must be a structure")
	}
	if valr.NumField() < len(arr) {
		return errors.New("too many fields in input for destination record")
	}

	for i := 0; i < len(arr); i++ {
		fld := valr.Field(i)
		switch fld.Kind() {
		case reflect.Int:
			tmp, err := strconv.Atoi(arr[i])
			if err != nil {
				return errors.New(fmt.Sprintf("Error at field %d. Not an integer."))
			}
			fld.SetInt(int64(tmp))
		case reflect.Float64:
			tmp, err := strconv.ParseFloat(arr[i], 64)
			if err != nil {
				return errors.New(fmt.Sprintf("Error at field %d. Not a float64."))
			}
			fld.SetFloat(tmp)
		case reflect.String:
			fld.SetString(arr[i])
		}
	}
	return nil
}

// Arr2Records will take a slice of rows of string slices and convert them to a slice of the
// type of the provided interface.
// Errors occur when the fields in a provided row is more than the fields of the interface,
// when parsing could not be done or the interface is not a structure
func Arr2Records(arr [][]string, itf interface{}) (interface{}, error) {
	valr := reflect.ValueOf(itf)
	if valr.Kind() != reflect.Struct {
		return nil, errors.New("the provided record must be a structure")
	}
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(itf)), 0, 100)
	for row := 0; row < len(arr); row++ {
		rowarr := arr[row]
		newval := reflect.Indirect(reflect.New(reflect.TypeOf(itf)))
		if newval.NumField() < len(rowarr) {
			return nil, errors.New(fmt.Sprintf("too many fields in input for destination record (row %d)", row))
		}
		for i := 0; i < len(rowarr); i++ {
			fld := newval.Field(i)
			switch fld.Kind() {
			case reflect.Int:
				tmp, err := strconv.Atoi(rowarr[i])
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Error at field %d. Not an integer."))
				}
				fld.SetInt(int64(tmp))
			case reflect.Float64:
				tmp, err := strconv.ParseFloat(rowarr[i], 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Error at field %d. Not a float64."))
				}
				fld.SetFloat(tmp)
			case reflect.String:
				fld.SetString(rowarr[i])
			}
		}
		slice = reflect.Append(slice, newval)
	}
	return slice.Interface(), nil
}

// Record2Arr will take the fields in the record and convert them to a slice of strings
func Record2Arr(rec interface{}) ([]string, error) {
	valr := reflect.Indirect(reflect.ValueOf(rec))
	if valr.Kind() != reflect.Struct {
		return nil, errors.New("the provided record must be a structure")
	}
	result := make([]string, 0, valr.NumField())
	for i := 0; i < valr.NumField(); i++ {
		fld := valr.Field(i)
		switch fld.Kind() {
		case reflect.Int:
			result = append(result, strconv.Itoa(int(fld.Int())))
		case reflect.Float64:
			result = append(result, fmt.Sprintf("%f", fld.Float()))
		case reflect.String:
			result = append(result, fld.String())
		}
	}
	return result, nil
}

// Records2Arr will take the slice of records and for each convert them to a slice of strings
func Records2Arr(recs interface{}) ([][]string, error) {
	valr := reflect.ValueOf(recs)
	if valr.Kind() != reflect.Slice {
		return nil, errors.New("the provided record must be a slice")
	}
	result := make([][]string, 0, valr.Len())
	if valr.Len() == 0 {
		return result, nil
	}
	if valr.Index(0).Kind() != reflect.Struct {
		return nil, errors.New("the elements of the array must be structures")
	}
	for row := 0; row < valr.Len(); row++ {
		rec := valr.Index(row)
		tmp := make([]string, 0, rec.NumField())
		for i := 0; i < rec.NumField(); i++ {
			fld := rec.Field(i)
			switch fld.Kind() {
			case reflect.Int:
				tmp = append(tmp, strconv.Itoa(int(fld.Int())))
			case reflect.Float64:
				tmp = append(tmp, fmt.Sprintf("%f", fld.Float()))
			case reflect.String:
				tmp = append(tmp, fld.String())
			}
		}
		result = append(result, tmp)
	}
	return result, nil
	return nil, nil
}

// CopyRec will copy the fields in structure src to the fields with the same name in structure dst
func CopyRec(dst interface{}, src interface{}) error {
	vals := reflect.Indirect(reflect.ValueOf(src))
	if vals.Kind() != reflect.Struct {
		return errors.New("the source must be a structure")
	}
	vald := reflect.Indirect(reflect.ValueOf(dst))
	if vald.Kind() != reflect.Struct {
		return errors.New("the destination must be a structure")
	}
	vals2 := reflect.TypeOf(src)
	for i := 0; i < vals.NumField(); i++ {
		fldsrc := vals.Field(i)
		flddst := vald.FieldByName(vals2.Field(i).Name)
		if flddst.IsNil() {
			return errors.New(fmt.Sprintf("Field not found n destination %s", vals2.Field(i).Name))
		}
		if fldsrc.Kind() != flddst.Kind() {
			return errors.New(fmt.Sprintf("not same types for field %s", vals2.Field(i).Name))
		}
		flddst.Set(fldsrc)
	}
	return nil
}

// CopyRecs will take a slice of records as input and generate a slice of records
// based on dstif with the fields copied by name
func CopyRecs(dstif interface{}, src interface{}) (interface{}, error) {
	vals := reflect.ValueOf(src)
	if vals.Kind() != reflect.Slice {
		return nil, errors.New("the source record must be an array")
	}
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(dstif)), 0, 100)

	if vals.Len() == 0 {
		return slice.Interface(), nil
	}
	if vals.Index(0).Kind() != reflect.Struct {
		return slice, errors.New("The source must be an array of structures")
	}

	vald := reflect.Indirect(reflect.ValueOf(dstif))
	if vald.Kind() != reflect.Struct {
		return nil, errors.New("the destination record must be a structure")
	}
	vals2 := vals.Index(0).Type()
	for row := 0; row < vals.Len(); row++ {
		fld := vals.Index(row)
		newval := reflect.Indirect(reflect.New(reflect.TypeOf(dstif)))
		for i := 0; i < fld.NumField(); i++ {
			fldsrc := fld.Field(i)
			flddst := newval.FieldByName(vals2.Field(i).Name)
			if fldsrc.Kind() != flddst.Kind() {
				return nil, errors.New(fmt.Sprintf("not same types for field %s", vals2.Field(i).Name))
			}
			flddst.Set(fldsrc)
		}
		slice = reflect.Append(slice, newval)
	}
	return slice.Interface(), nil
}
