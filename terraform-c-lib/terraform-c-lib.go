package main

import (
	"C"
	"log"
	"fmt"
	//	"github.com/hashicorp/terraform/internal/configs"
	"github.com/hashicorp/terraform/internal/configs/configload"

)

//export LoadModule
func LoadModule(modulePath string) (rc int, config_c *C.char, errstr *C.char) {
	log.Printf("modulePath: %s", modulePath)
	modulePathClean := modulePath[:len(modulePath)-1]

	config := &configload.Config{}
	loader, loaderErr := configload.NewLoader(config)
	if loaderErr != nil {
		log.Printf("Error creating loader: %s", loaderErr )
		return -1, nil, C.CString(fmt.Sprintf("%v", loaderErr))
	}


	loadedConfig, configErr := loader.LoadConfig(modulePathClean)
	if configErr != nil {
		log.Printf("Error loading config: %s", configErr )
		return -2, nil, C.CString(fmt.Sprintf("%v", configErr))
	}
	log.Printf("Module is %#v", loadedConfig.Module)
	return 0, C.CString(fmt.Sprintf("%v", loadedConfig.Module)), nil
}

func main() {

	rc, config, err := LoadModule("../terraform-c-lib/test_tf_files/")

	log.Printf("rc: %d", rc)
	log.Printf("config: %s", config)
	log.Printf("err: %v", err)

}

