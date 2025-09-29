/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package audience

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateAudience(dataRaw []byte) ([]byte, error) {
	var audiencePayload models.AudiencePayload
	err := json.Unmarshal(dataRaw, &audiencePayload)
	if err != nil {
		return nil, err
	}

	for _, v := range audiencePayload.Groups {
		for _, t := range v {

			targetingType := t.Type
			b, err := json.Marshal(t.Conditions)
			if err != nil {
				return nil, err
			}

			dec := json.NewDecoder(bytes.NewReader(b))
			dec.DisallowUnknownFields()
			switch targetingType.(type) {
			case string:
				switch targetingType {
				case COOKIE:
					var cookieModel []models.Cookie
					if err := dec.Decode(&cookieModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case ACTION_TRACKING:
					var actionTrackingModel []models.ActionTracking
					if err := dec.Decode(&actionTrackingModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case SESSION_NUMBER:
					var sessionNumberModel []models.SessionNumber
					if err := dec.Decode(&sessionNumberModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case RETURNING_VISITOR:
					var returningVisitorModel []models.NewOrReturningVisitorPayload
					if err := dec.Decode(&returningVisitorModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case PROVIDERS:
					var providerModel []models.Provider
					if err := dec.Decode(&providerModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case PAGES_INTEREST:
					var pagesInterestModel []models.PageInterest
					if err := dec.Decode(&pagesInterestModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case CAMPAIGN_EXPOSITION:
					var campaignExpositionModel []models.CampaignExposure
					if err := dec.Decode(&campaignExpositionModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case CUSTOM_VARIABLE:
					var customVariableModel []models.CustomVariable
					if err := dec.Decode(&customVariableModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case PAGE_VIEW:
					var pageViewModel []models.PageViewPayload
					if err := dec.Decode(&pageViewModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case DEVICE:
					var deviceModel []models.Device
					if err := dec.Decode(&deviceModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case IP:
					var ipModel []models.IPRange
					if err := dec.Decode(&ipModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case DATALAYER:
					var dataLayerModel []models.DataLayer
					if err := dec.Decode(&dataLayerModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case CODE:
					var codeModel []models.Code
					if err := dec.Decode(&codeModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case GEOLOCALISATION:
					var geolocalisationModel []models.GeoLocation
					if err := dec.Decode(&geolocalisationModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case URL_PARAMETER:
					var urlParameterModel []models.UrlParameter
					if err := dec.Decode(&urlParameterModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case SELECTOR:
					var selectorModel []models.Selector
					if err := dec.Decode(&selectorModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case LANDING_PAGE:
					var landingPageModel []models.LandingPage
					if err := dec.Decode(&landingPageModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case NUMBER_PAGE_VIEWED:
					var numberPageViewedModel []models.NumberPageView
					if err := dec.Decode(&numberPageViewedModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case JS_VARIABLE:
					var jsVariableModel []models.JSVariable
					if err := dec.Decode(&jsVariableModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case BROWSER:
					var browserModel []models.Browser
					if err := dec.Decode(&browserModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case SCREEN_SIZE:
					var screenSizeModel []models.ScreenSize
					if err := dec.Decode(&screenSizeModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				case PREVIOUS_PAGE:
					var previousPageModel []models.PreviousPage
					if err := dec.Decode(&previousPageModel); err != nil {
						return nil, fmt.Errorf("%v", err)
					}
				default:
					return nil, fmt.Errorf("Type not supported")
				}

			default:
				return nil, fmt.Errorf("Type format not supported")
			}

		}
	}

	audienceHeader, err := httprequest.AudienceRequester.HTTPCreateAudience(dataRaw)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(audienceHeader), "/")
	audienceID := parts[len(parts)-1]
	body, err := httprequest.AudienceRequester.HTTPGetAudience(audienceID)
	if err != nil {
		return nil, err
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bodyByte, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create an audience",
	Long:  `Create an audience`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := CreateAudience([]byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %s", err)

		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(resp))
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your audience, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	AudienceCmd.AddCommand(createCmd)

}

// Trigger and Segment
const COOKIE string = "COOKIE"                           // 20
const RETURNING_VISITOR string = "RETURNING_VISITOR"     // 24
const PROVIDERS string = "PROVIDERS"                     // 53
const PAGES_INTEREST string = "PAGES_INTEREST"           // 52
const CAMPAIGN_EXPOSITION string = "CAMPAIGN_EXPOSITION" // 29
const CUSTOM_VARIABLE string = "CUSTOM_VARIABLE"         // 41
const PAGE_VIEW string = "PAGE_VIEW"                     // 51

// Trigger Only
const DEVICE string = "DEVICE"                         // 17
const IP string = "IP"                                 // 18
const DATALAYER string = "DATALAYER"                   // 44
const CODE string = "CODE"                             // 40
const GEOLOCALISATION string = "GEOLOCALISATION"       // 19
const URL_PARAMETER string = "URL_PARAMETER"           // 39
const SELECTOR string = "SELECTOR"                     // 43
const LANDING_PAGE string = "LANDING_PAGE"             // 22
const NUMBER_PAGE_VIEWED string = "NUMBER_PAGE_VIEWED" // 31
const JS_VARIABLE string = "JS_VARIABLE"               // 28
const BROWSER string = "BROWSER"                       // 23
const SCREEN_SIZE string = "SCREEN_SIZE"               // 27
const PREVIOUS_PAGE string = "PREVIOUS_PAGE"           // 26
const SESSION_NUMBER string = "SESSION_NUMBER"         // 34
const ACTION_TRACKING string = "ACTION_TRACKING"       // 42
