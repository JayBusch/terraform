package main

import (
	"C"
	"log"
	"fmt"
	"reflect"
	"unicode"
	"strconv"
//	"encoding/xml"
	"github.com/hashicorp/terraform/internal/configs"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform/internal/configs/configload"
	"github.com/hashicorp/terraform/internal/addrs"
	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/zclconf/go-cty/cty"

)

type System struct {
	backend *Backend
	resources map[string]*InfraResource
}

type Backend struct {
	backend_type string
	attributes map[string]string
}

type InfraResource struct {
	resource_mode *ResourceMode
	providers map[string]*Provider
	name string
	resource_type string
	config string

	attributes map[string]string

	children []*InfraResource
	parent *InfraResource
}

type Provider struct {
	provider_type string
	namespace string
	hostname string
}

type ResourceMode struct {
	name string
	resource_mode_type string
	attributes map[string]string
}

var system System

/*
func (m Module) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var s string
	switch a {
	default:
		s = "unknown"
	case Gopher:
		s = "gopher"
	case Zebra:
		s = "zebra"
	}
	return e.EncodeElement(s, start)
}
*/

func getTypeName(data interface{}) string {
	if data == nil {
	    return "nil"
	}

	if t := reflect.TypeOf(data); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func nextVal(data interface{}, fieldName string, path []string) interface{} {

	fmt.Printf("path: %#+v\n\n\n", path)

	//fmt.Printf("getTypeName(data): %#+v\n\n", getTypeName(data))

	switch getTypeName(data) {

	case "*Module":
		fmt.Printf("found Module\n")
		moduleData := data.(*configs.Module)
		vData := reflect.ValueOf(*moduleData)
		tData := reflect.TypeOf(*moduleData)

		for i:=0;i<tData.NumField();i++{
			//fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				return nil
			}
			newPath := append(path,tData.Field(i).Name)
			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}
		fmt.Printf("finished Module\n")
		return nil

	case "*Backend":
		fmt.Printf("found Backend\n")
		backendData := data.(*configs.Backend)
		vData := reflect.ValueOf(*backendData)
		tData := reflect.TypeOf(*backendData)

		system.backend = &Backend{
			attributes : make(map[string]string),
		}

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}
		fmt.Printf("finished Backend\n")
		return nil

	case "*Body":
		fmt.Printf("found Body\n")
		bodyData := data.(*hclsyntax.Body)
		vData := reflect.ValueOf(*bodyData)
		tData := reflect.TypeOf(*bodyData)

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}
		fmt.Printf("finished Body\n")
		return nil

	case "Blocks":
		fmt.Printf("found Blocks\n")
		blocksData := data.(hclsyntax.Blocks)
		vData := reflect.ValueOf(blocksData)
		tData := reflect.TypeOf(blocksData)

		fmt.Printf("BLOCKS - tData: %#+v\n vData: %#+v\n",tData,vData)

/*
		for i:=0;i<tData.NumField();i++{
			fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}

			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)

		}
		*/
		fmt.Printf("finished Blocks\n")
		return nil

	case "*RequiredProviders":
		fmt.Printf("found RequiredProviders\n")
		rpData := data.(*configs.RequiredProviders)
		vData := reflect.ValueOf(*rpData)
		tData := reflect.TypeOf(*rpData)

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}
		fmt.Printf("finished RequiredProviders\n")
		return nil

	case "Resource":
		fmt.Printf("found Resource\n")
		rData := data.(configs.Resource)
//		fmt.Printf("rData: %#+v",rData)

		vData := reflect.ValueOf(rData)
		tData := reflect.TypeOf(rData)

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			path = append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, path)
		}

		fmt.Printf("finished Resources\n")
		return nil

	case "Provider":
		fmt.Printf("found Provider\n")
		pData := data.(addrs.Provider)
//		fmt.Printf("pData: %#+v\n\n",pData)

		vData := reflect.ValueOf(pData)
		tData := reflect.TypeOf(pData)

		system.resources[path[1]].providers = make(map[string]*Provider)

		for i:=0;i<tData.NumField();i++{
			//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}

		fmt.Printf("finished Provider\n")
		return nil

	case "Attributes":
		fmt.Printf("found Attributes\n")
		aData := data.(hclsyntax.Attributes)
//		fmt.Printf("aData: %#+v\n\n",aData)

		for k,v := range aData {
			newPath := append(path,k)

			nextVal(v,k, newPath)
		}

/*
		vData := reflect.ValueOf(aData)
		tData := reflect.TypeOf(aData)

		for i:=0;i<tData.NumField();i++{
			fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				return nil
			}

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name)
		}
*/
		fmt.Printf("finished Attributes\n")
		return nil

	case "*Attribute":
		fmt.Printf("found *Attribute\n")
		aData := data.(*hclsyntax.Attribute)
//		fmt.Printf("aData: %#+v\n\n",*aData)

		vData := reflect.ValueOf(*aData)
		tData := reflect.TypeOf(*aData)

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}

		fmt.Printf("finished *Attribute\n")
		return nil

	case "*TemplateExpr":
		fmt.Printf("found *TemplateExpr\n")
		teData := data.(*hclsyntax.TemplateExpr)
		fmt.Printf("teData: %#+v\n\n",*teData)

		vData := reflect.ValueOf(*teData)
		tData := reflect.TypeOf(*teData)

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}

		fmt.Printf("finished *TemplateExpr\n")
		return nil

	case "*LiteralValueExpr":
		fmt.Printf("found *LiteralExpr\n")
		leData := data.(*hclsyntax.LiteralValueExpr)
		fmt.Printf("leData: %#+v\n\n",*leData)

		vData := reflect.ValueOf(*leData)
		tData := reflect.TypeOf(*leData)

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}

		fmt.Printf("finished *LiteralExpr\n")
		return nil

	case "Value":
		fmt.Printf("found Value\n\n")

		valData := data.(cty.Value)

		vData := reflect.ValueOf(valData)
		tData := reflect.TypeOf(valData)

		//fmt.Printf("valData.Type().GoString(): %s\n\n", valData.Type().GoString())

		switch valData.Type().GoString() {
		case "cty.Bool":
			fmt.Printf(" valData.True(): %t\n", valData.True())
		case "cty.String":
			fmt.Printf(" valData.AsString(): %s\n", valData.AsString())
			if len(path) > 6 {
				if path[len(path)- 6] == "Attributes" {
					if path[0] == "Backend" {
						system.backend.attributes[path[len(path)-5]]=valData.AsString()
						fmt.Printf("system.backend: %#+v\n", system.backend)
					} else if path[0] == "ManagedResources" {
						system.resources[path[1]].attributes[path[len(path)-5]] = valData.AsString()
						fmt.Printf("system.resources[path[1]: %#+v\n", system.resources[path[1]])
					}
				}
			}
		default:
			fmt.Printf(" UNKNOWN\n")
		}

		fmt.Printf(" path[len(path)-6]: %s\n", path[len(path)-6])


		//fmt.Printf("vData.GoString(): %s\n\n",string(vData))


		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}


		fmt.Printf("finished Value\n\n")
		return nil

	case "*ObjectConsExpr":
		fmt.Printf("found *ObjectConstExpr\n")

		oceData := data.(*hclsyntax.ObjectConsExpr)
		//fmt.Printf("oceData: %#+v\n\n",*oceData)

		vData := reflect.ValueOf(*oceData)
		tData := reflect.TypeOf(*oceData)

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}


		fmt.Printf("finished *ObjectConstExpr\n")

		return nil
	case "ResourceMode":

		fmt.Printf("found ResourceMode\n")
		fmt.Printf("path: %s\n", path)

		rmData := data.(addrs.ResourceMode)
		//fmt.Printf("rmData: %#+v\n\n",rmData)

		vData := reflect.ValueOf(rmData)
		fmt.Printf("ResourceMode: %s\n\n", vData.Interface())

		if path[0] == "ManagedResources" {
			if len(path) == 3 {
				if path[2] == "Mode" {
					system.resources[path[1]].resource_mode = &ResourceMode{}
				}
			}
		}

		return nil
/*		tData := reflect.TypeOf(rmData)

		for i:=0;i<tData.NumField();i++{
			fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}


		return nil
*/

	case "Range":
		fmt.Printf("found Range\n")
		rngData := data.(hcl.Range)
		vData := reflect.ValueOf(rngData)
		tData := reflect.TypeOf(rngData)

		for i:=0;i<tData.NumField();i++{
		//	fmt.Printf("processing: %s\n\n", tData.Field(i).Name)
			if tData.Field(i).Name != "" && !unicode.IsUpper([]rune(tData.Field(i).Name)[0]){
				continue
			}
			newPath := append(path,tData.Field(i).Name)

			nextVal(vData.Field(i).Interface(),tData.Field(i).Name, newPath)
		}

		fmt.Printf("finished Range\n")
		return nil
	default:
		fmt.Printf("***Unknown Type***: %#+v\n\n", getTypeName(data) )
	}

	switch fieldName {
	case "SourceDir":
		fmt.Printf("SourceDir %#+v\n\n",data)
	case "CoreVersionConstraints":
		fmt.Printf("CoreVersionConstraints %#+v\n\n",data)
	case "ActiveExperiments":
		fmt.Printf("ActiveExperiments %#+v\n\n",data)
	case "Backend":
		fmt.Printf("Backend %#+v\n\n",data)
	case "CloudConfig":
		fmt.Printf("CloudConfig %#+v\n\n",data)
	case "ProviderConfigs":
		fmt.Printf("ProviderConfigs %#+v\n\n", data)
	case "ProviderRequirements":
		fmt.Printf("ProviderRequirements %#+v\n\n",data)
	case "ProviderLocalNames":
		fmt.Printf("ProviderLocalNames %#+v\n\n",data)
	case "ProverMetas":
		fmt.Printf("ProviderMetas %#+v\n\n",data)
	case "Variables":
		fmt.Printf("Variables %#+v\n\n",data)
	case "Locals":
		fmt.Printf("Locals %#+v\n\n",data)
	case "Outputs":
		fmt.Printf("Outputs %#+v\n\n",data)
	case "ModuleCalls":
		fmt.Printf("ModuleCalls %#+v\n\n",data)
	case "ManagedResources":
		fmt.Printf("ManagedResources %#+v\n\n",data)
		fmt.Printf("found ManagedResources\n")
		mrData := data.(map[string]*configs.Resource)
//		fmt.Printf("mrData: %#+v",mrData)

		system.resources = make(map[string]*InfraResource)

		for k,v := range mrData{
			//			fmt.Println(*v)
			newPath := append(path,k)

			system.resources[k] = &InfraResource{
				name: k,
				attributes : make(map[string]string),
			}
			fmt.Printf("system.resources: %#+v\n\n", system.resources)

			nextVal(*v,k,newPath)
		}

		fmt.Printf("finished ManagedResources\n")

	case "DataResources":
		fmt.Printf("DataResources %#+v\n\n",data)
	case "Moved":
		fmt.Printf("Moved %#+v\n\n",data)
	case "ResourceMode":
		fmt.Printf("ResourceMode %#+v\n\n", data)
	case "Parts":
		fmt.Printf("Parts %#+v\n\n", data)
		fmt.Printf("found Parts\n")
		ptsData := data.([]hclsyntax.Expression)
//		fmt.Printf("ptsData: %#+v",ptsData)

		for i, v := range ptsData{
			//			fmt.Println(*v)
			newPath := append(path,"Part " + fmt.Sprint(i))

			nextVal(v,"Part " + fmt.Sprint(i),newPath)
		}

		fmt.Printf("finished Parts\n")

	case "Name":
		fmt.Printf("Name %s\n\n", data.(string))
		fmt.Printf("path: %s\n\n", path)

		if len(path) > 1 {
			if path[0] == "ManagedResources" {
				if len(path) == 4 {
					if path[2] == "Mode" {
						if path[3] == "Name" {
							system.resources[path[1]].resource_mode.name = data.(string)
							fmt.Printf("system.resources[path[1].resource_mode: %#+v\n\n", system.resources[path[1]].resource_mode)
						}
					}
				}
			}
		}

	case "Type":
		fmt.Printf("Type %#+v\n\n",data)

		if len(path) > 1  {
			if path[0] == "Backend" {
				fmt.Printf("System.backend.type found\n")
				system.backend.backend_type = data.(string)
				fmt.Printf("system: %#+v\n", system)
			} else if path[len(path)-2] == "Provider" {
				system.resources[path[1]].providers[strconv.Itoa(len(system.resources[path[1]].providers))] = &Provider{provider_type: data.(string)}
				fmt.Printf("system.resources[path[1]].providers: %#+v\n", system.resources[path[1]].providers)
			}
		}
	case "Namespace":
		fmt.Printf("Namespace %#+v\n\n",data)

		if len(path) > 1 {
			if path[len(path)-2] == "Provider" {
				system.resources[path[1]].providers[strconv.Itoa(len(system.resources[path[1]].providers)-1)].namespace = data.(string)
				fmt.Printf("system.resources[path[1]].providers[strconv.Itoa(len(system.resources[path[1]].providers)-1)].namespace: %#+v\n", system.resources[path[1]].providers[strconv.Itoa(len(system.resources[path[1]].providers)-1)].namespace)
			}
		}
	case "Hostname":
		fmt.Printf("Hostname %#+v\n\n",data)
		if len(path) > 1 {
			if path[len(path)-2] == "Provider" {
//				newPath := append(path,"Provider")
//				nextVal(data,"Provider ",newPath)

				system.resources[path[1]].providers[strconv.Itoa(len(system.resources[path[1]].providers)-1)].hostname = data.(svchost.Hostname).String()

				fmt.Printf("system.resources[path[1]].providers[strconv.Itoa(len(system.resources[path[1]].providers)-1)].hostname: %#+v\n", system.resources[path[1]].providers[strconv.Itoa(len(system.resources[path[1]].providers)-1)].hostname)

			}
		}

	default:
		fmt.Printf("Default %#+v\n\n", data)
	}


	return data
}

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

	fmt.Printf("Config is %#+v\n\n", loadedConfig)
	fmt.Printf("Module is %#+v\n\n", loadedConfig.Module)

	_, diags := loadedConfig.Module.Backend.Config.JustAttributes()

//	fmt.Printf("attributes[\"bucket\"].Expr.Value(nil): %#+v\n\n", attributes["bucket"].Expr.Value(nil))

	fmt.Printf("diags: %#+v\n\n", diags)

	//fmt.Printf("Backend.Config.Content().Attributes is %#+v\n\n", loadedConfig.Module.Backend.Config.Content().Attributes)

	nextVal(loadedConfig.Module, "Module", []string{})

	return 0, C.CString(fmt.Sprintf("%v", loadedConfig.Module)), nil
}

func main() {

}

