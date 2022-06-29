package main

import (
	"testing"
)

func TestLoadModule(t *testing.T){

	rc, config, err := LoadModule("../terraform-c-lib/test_tf_files/")

	if err != nil {
		t.Fatalf("rc: %d, config: %v, err: %v",rc,config,err)
	}

}
