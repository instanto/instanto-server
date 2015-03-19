package main

import (
	"encoding/json"
	"fmt"
	lib "github.com/instanto/instanto-lib"
)

func HelperValidationError2JSON(verr *lib.ValidationError) (verrJSON []byte, err error) {
	verrJSON, err = json.Marshal(verr)
	if err != nil {
		return
	}
	return
}
func HelperValidationErrors2JSON(verrs []*lib.ValidationError) (verrsJSON []byte, err error) {
	data := make(map[string]interface{})
	data["errors"] = verrs
	verrsJSON, err = json.Marshal(data)
	if err != nil {
		return
	}
	return
}
func LogError(err error) {
	fmt.Println(err)
}
func LogInfo(things ...interface{}) {
	fmt.Println(things)
}
