package test

import (
	"fmt"
	"testing"
)

func TestDemo(t *testing.T) {
	data := map[string]int64{"demo":1, "demov2":222, "demov3":3333}
	
	for key, val := range data {
		fmt.Println(key, val)
		if val == 222 {
			delete(data, key)
		}
	}
	
	for key, val := range data {
		fmt.Println(key, val)
		if val == 222 {
			delete(data, key)
		}
	}
}