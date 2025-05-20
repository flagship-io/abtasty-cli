/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/audience"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/campaign"
	favorite_url "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/favorite-url"
	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	resourceFile    string
	outputFile      string
	inputParams     string
	inputParamsFile string
)

var inputParamsMap map[string]interface{}

type Data interface {
	getName() string
	Save(data []byte) ([]byte, error)
	Delete(id string) error
}

type ResourceData struct {
	Id string `json:"id"`
}

type CampaignData struct {
	*models.CampaignWEResourceLoader
}

// getName implements Data.
func (f *CampaignData) getName() string {
	return "Campaign"
}

// DeleteWithParent implements Data.
func (f *CampaignData) DeleteWithParent(parentId string, id string) error {
	panic("unimplemented")
}

// SaveWithParent implements Data.
func (f *CampaignData) SaveWithParent(parentId string, data []byte) ([]byte, error) {
	panic("unimplemented")
}

func (f *CampaignData) Save(data []byte) ([]byte, error) {
	return campaign.CreateCampaign(data), nil
}

func (f *CampaignData) Delete(id string) error {
	return http_request.CampaignWERequester.HTTPDeleteCampaign(id)
}

func (f *CampaignData) Switch(id, state string) error {
	return http_request.CampaignWERequester.HTTPSwitchStateCampaign(id, state)
}

type AudienceData struct {
	*models.AudienceResourceLoader
}

// getName implements Data.
func (f *AudienceData) getName() string {
	return "Audience"
}

func (f *AudienceData) Save(data []byte) ([]byte, error) {
	return audience.CreateAudience(data), nil
}

// Can't delete audience for this moment
func (f *AudienceData) Delete(id string) error {
	return nil
}

type FavoriteUrlData struct {
	*models.FavoriteURL
}

// Can't create favorite URL for this moment
func (f *FavoriteUrlData) Save(data []byte) ([]byte, error) {
	return favorite_url.CreateFavoriteURL(data), nil
}

// Can't delete favorite URL for this moment
func (f *FavoriteUrlData) Delete(id string) error {
	return nil
}

/* type ModificationData struct {
	*models.Modification
}

// Delete implements Data.
func (f *ModificationData) Delete(id string) error {
	panic("unimplemented")
}

// Save implements Data.
func (f *ModificationData) Save(data string) ([]byte, error) {
	panic("unimplemented")
}

func (f *ModificationData) SaveWithParent(campaignId, data string) ([]byte, error) {
	var modelData models.ModificationCodeCreateStruct
	campaignIdInt, err := strconv.Atoi(campaignId)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	err = json.Unmarshal([]byte(data), &modelData)

	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	return http_request.ModificationRequester.HTTPCreateModification(campaignIdInt, modelData)
}

func (f *ModificationData) DeleteWithParent(campaignId, id string) error {
	campaignIdInt, err := strconv.Atoi(campaignId)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	return http_request.ModificationRequester.HTTPDeleteModification(campaignIdInt, idInt)
} */

type VariationData struct {
	*models.VariationWE
}

func (f *VariationData) Save(campaignId int, data models.VariationWE) ([]byte, error) {
	return http_request.VariationWERequester.HTTPCreateVariation(campaignId, data)
}

func (f *VariationData) Delete(campaignId, id int) error {
	return http_request.VariationWERequester.HTTPDeleteVariation(campaignId, id)
}

func resolveVariables(data interface{}, resourceVariables map[string]interface{}) interface{} {
	switch val := data.(type) {
	case map[string]interface{}:
		for k, v := range val {
			val[k] = resolveVariables(v, resourceVariables)
		}

	case []interface{}:
		for i, v := range val {
			val[i] = resolveVariables(v, resourceVariables)
		}

	case string:
		if strings.Contains(val, "$") {
			vTrim := strings.Trim(val, "$")
			for k_, variable := range resourceVariables {
				script, _ := tengo.Eval(context.Background(), vTrim, map[string]interface{}{
					k_: variable,
				})
				if script == nil {
					continue
				}
				// Update the string value with the result
				if resultStr, ok := script.(string); ok {
					return resultStr
				}
			}
		}
	}
	return data
}

type ResourceType int

const (
	Campaign ResourceType = iota
	Audience
	FavoriteURL
	Modification
	Variation
)

var resourceTypeMap = map[string]ResourceType{
	"audience":     Audience,
	"favorite_url": FavoriteURL,
	"modification": Modification,
	"campaign":     Campaign,
	"variation":    Variation,
}

type Resource struct {
	Name             ResourceType
	Data             Data
	ResourceVariable string
	Method           string
}

var cred common.RequestConfig

func Init(credL common.RequestConfig) {
	cred = credL
}

type ResourceCmdStruct struct {
	Name             string `json:"name,omitempty"`
	ResourceVariable string `json:"resource_variable,omitempty"`
	Response         string `json:"response,omitempty"`
	Method           string `json:"method,omitempty"`
	Error            string `json:"error,omitempty"`
}

func UnmarshalConfig(filePath string) ([]Resource, error) {
	var config struct {
		Resources []struct {
			Name             string
			Data             json.RawMessage
			ResourceVariable string
			Method           string
		}
	}

	bytes, err := os.ReadFile(resourceFile)

	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	var resources []Resource
	for _, r := range config.Resources {
		name, ok := resourceTypeMap[r.Name]
		if !ok {
			return nil, fmt.Errorf("invalid resource name: %s", r.Name)
		}

		var data Data = nil
		var err error = nil

		switch name {

		case Audience:
			audienceData := AudienceData{}
			err = json.Unmarshal(r.Data, &audienceData)
			data = &audienceData

		/* case Modification:
		modificationData := ModificationData{}
		err = json.Unmarshal(r.Data, &modificationData)
		data = &modificationData */

		case Campaign:
			campaignData := CampaignData{}
			err = json.Unmarshal(r.Data, &campaignData)
			data = &campaignData

		}

		if err != nil {
			return nil, err
		}

		resources = append(resources, Resource{Name: name, Data: data, ResourceVariable: r.ResourceVariable, Method: r.Method})
	}

	return resources, nil
}

var gResources []Resource

// LoadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [--file=<file>]",
	Short: "Load your resources",
	Long:  `Load your resources`,
	Run: func(cmd *cobra.Command, args []string) {
		var params []byte

		if resourceFile == "" && !utils.CheckSingleFlag(inputParams != "", inputParamsFile != "") {
			log.Fatalf("error occurred: %s", "1 flag is required. (input-params, input-params-file)")
		}

		if inputParams != "" {
			params = []byte(inputParams)

		}

		if inputParamsFile != "" {
			fileContent, err := os.ReadFile(inputParamsFile)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			params = fileContent
		}

		if params != nil {
			err := json.Unmarshal(params, &inputParamsMap)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error: %s", err)
				return
			}
		}

		jsonBytes := ScriptResource(cmd, gResources, inputParamsMap)
		if outputFile != "" {
			os.WriteFile(outputFile, jsonBytes, os.ModePerm)
			fmt.Fprintf(cmd.OutOrStdout(), "File created at %s\n", outputFile)
			return
		}

		if viper.GetString("output_format") == "json" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s", string(jsonBytes))
		}

	},
}

func init() {
	cobra.OnInitialize(initResource)

	loadCmd.Flags().StringVarP(&resourceFile, "file", "", "", "resource file that contains your resource")

	if err := loadCmd.MarkFlagRequired("file"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	loadCmd.Flags().StringVarP(&outputFile, "output-file", "", "", "result of the command that contains all resource information")

	loadCmd.Flags().StringVarP(&inputParams, "input-params", "", "", "params to replace resource loader file")
	loadCmd.Flags().StringVarP(&inputParamsFile, "input-params-file", "", "", "file that contains params to replace resource loader file")

	ResourceCmd.AddCommand(loadCmd)
}

func initResource() {

	// Use config file from the flag.
	var err error
	if resourceFile != "" {
		gResources, err = UnmarshalConfig(resourceFile)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	}
}

func ScriptResource(cmd *cobra.Command, resources []Resource, inputParamsMap map[string]interface{}) []byte {

	resourceVariables := make(map[string]interface{})
	var loadResultJSON []string
	var loadResultOutputFile []ResourceCmdStruct

	for _, resource := range resources {
		var response []byte
		var resultOutputFile ResourceCmdStruct
		var resourceData map[string]interface{}
		var responseData interface{}

		var resourceName = resource.Data.getName()
		const color = "\033[0;33m"
		const colorNone = "\033[0m"

		data, err := json.Marshal(resource.Data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error occurred marshal data: %v\n", err)
		}

		var httpMethod string = "POST"

		if resource.Method == "delete" {
			httpMethod = "DELETE"
		}

		if resource.Method == "switch" {
			httpMethod = "PATCH"
		}

		err = json.Unmarshal(data, &resourceData)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error occurred unmarshal resourceData: %v\n", err)
		}

		if inputParamsMap != nil {
			for k, vInterface := range resourceData {
				v, ok := vInterface.(string)
				if ok {
					if strings.Contains(v, "$") {
						vTrim := strings.Trim(v, "$")
						vTrimL := strings.Split(vTrim, ".")
						value, err := getNestedValue(inputParamsMap, vTrimL)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Error: %s\n", err)
						}

						if value != nil {
							resourceData[k] = value
						}
					}
				}
			}
		}

		resourceData = resolveVariables(resourceData, resourceVariables).(map[string]interface{})

		dataResource, err := json.Marshal(resourceData)
		if err != nil {
			log.Fatalf("error occurred http call: %v\n", err)
		}

		if resource.Method == "switch" {
			if resource.Name == Campaign {
				fmt.Println("--id="+fmt.Sprintf("%v", resourceData["id"]), "--state="+fmt.Sprintf("%v", resourceData["state"]))
				if !(fmt.Sprintf("%v", resourceData["state"]) == "active" || fmt.Sprintf("%v", resourceData["state"]) == "paused" || fmt.Sprintf("%v", resourceData["state"]) == "interrupted") {
					fmt.Fprintln(cmd.OutOrStdout(), "Status can only have 3 values: active or paused or interrupted")
				} else {
					err := http_request.CampaignFERequester.HTTPSwitchStateCampaign(fmt.Sprintf("%v", resourceData["id"]), fmt.Sprintf("%v", resourceData["state"]))
					if err != nil {
						log.Fatalf("error occurred: %v", err)
					}
					fmt.Fprintf(cmd.OutOrStdout(), "campaign status set to %s\n", fmt.Sprintf("%v", resourceData["state"]))
				}
			}
		}

		if httpMethod == "POST" {
			response, err = resource.Data.Save(dataResource)

			resultOutputFile = ResourceCmdStruct{
				Name:             resourceName,
				Response:         string(response),
				ResourceVariable: resource.ResourceVariable,
				Method:           httpMethod,
			}

			if err != nil {
				resultOutputFile.Error = err.Error()
			}

			loadResultOutputFile = append(loadResultOutputFile, resultOutputFile)
		}

		if httpMethod == "DELETE" {
			//_, err = common.HTTPRequest[ResourceData](httpMethod, utils.GetWebExperimentationHost()+"/v1/accounts/"+cred.AccountID+url+"/"+fmt.Sprintf("%s", resourceData["id"]), nil)
			err = resource.Data.Delete(fmt.Sprintf("%s", resourceData["id"]))

			if err == nil && viper.GetString("output_format") != "json" {
				response = []byte("The id: " + fmt.Sprintf("%v", resourceData["id"]) + " deleted successfully")
			}
		}

		if err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "%s - %s: %s %s\n", color, resourceName, colorNone, err.Error())
			continue
		}

		if viper.GetString("output_format") != "json" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s - %s: %s %s\n", color, resourceName, colorNone, string(response))
		}

		if httpMethod != "DELETE" && httpMethod != "PATCH" {
			err = json.Unmarshal(response, &responseData)

			if err != nil {
				fmt.Fprintf(os.Stderr, "error occurred unmarshal responseData: %v\n", err)
			}

			if responseData == nil {
				fmt.Fprintf(os.Stderr, "error occurred not response data: %s\n", string(response))
				continue
			}

			resourceVariables[resource.ResourceVariable] = responseData
		}

		loadResultJSON = append(loadResultJSON, string(response))
	}

	var jsonBytes []byte
	var jsonString any

	if outputFile != "" {
		jsonString = loadResultOutputFile
	} else {
		jsonString = loadResultJSON
	}

	jsonBytes, err := json.Marshal(jsonString)

	if err != nil {
		log.Fatalf("Error marshaling struct: %v", err)
	}

	return jsonBytes
}

func getNestedValue(data map[string]interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	current := data
	for i, key := range path {
		value, ok := current[key]
		if !ok {
			continue
		}
		if i == len(path)-1 {
			return value, nil
		}
		next, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("value at key '%s' is not an object", key)
		}
		current = next
	}

	return nil, nil
}
