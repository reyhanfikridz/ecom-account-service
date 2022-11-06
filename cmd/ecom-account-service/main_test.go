/*
Package main the executeable file
*/
package main

import "testing"

// TestInitAPI test InitAPI
func TestInitAPI(t *testing.T) {
	// init API
	_, err := InitAPI()
	if err != nil {
		t.Errorf("Initialization of API failed => " + err.Error())
	}
}
