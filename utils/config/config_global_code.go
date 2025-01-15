package config

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func CheckWorkingDirectory(workingDir string) (string, error) {

	if _, err := os.Stat(workingDir); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return "", err
		}
	}

	return workingDir, nil
}

func CheckGlobalCodeDirectory(workingDir string) (string, error) {
	wd, err := CheckWorkingDirectory(workingDir)
	if err != nil {
		return "", err
	}

	gcWorkingDir := wd + "/.abtasty"

	err = os.MkdirAll(gcWorkingDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return gcWorkingDir, nil
}

func AccountGlobalCodeFilePath(workingDir, accountID string) string {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	accountCodeDir := gcWorkingDir + "/" + accountID

	err = os.MkdirAll(accountCodeDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	jsFilePath := accountCodeDir + "/accountGlobalCode.js"
	return jsFilePath
}

func WriteAccountGlobalCode(workingDir, accountID, code string) (string, error) {
	jsFilePath := AccountGlobalCodeFilePath(workingDir, accountID)
	err := os.WriteFile(jsFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+jsFilePath)
	return jsFilePath, nil
}

func CampaignGlobalCodeFilePath(workingDir, accountID, campaignID string) string {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID

	err = os.MkdirAll(campaignCodeDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	jsFilePath := campaignCodeDir + "/campaignGlobalCode.js"
	return jsFilePath
}

func WriteCampaignGlobalCode(workingDir, accountID, campaignID, code string) (string, error) {
	jsFilePath := CampaignGlobalCodeFilePath(workingDir, accountID, campaignID)
	err := os.WriteFile(jsFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+jsFilePath)
	return jsFilePath, nil
}

func DeleteCampaignGlobalCodeDirectory(workingDir, accountID, campaignID string) (string, error) {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		return "", err
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID

	if _, err := os.Stat(campaignCodeDir); err == nil {
		err := os.RemoveAll(campaignCodeDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error deleting file: ", err)
			return "", err
		}
	}

	return "campaign directory for " + campaignID + " deleted", nil
}

func VariationGlobalCodeJSFilePath(workingDir, accountID, campaignID, variationID string) string {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID
	variationCodeDir := campaignCodeDir + "/" + variationID

	err = os.MkdirAll(variationCodeDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	jsFilePath := variationCodeDir + "/variationGlobalCode.js"
	return jsFilePath
}

func WriteVariationGlobalCodeJS(workingDir, accountID, campaignID, variationID, code string) (string, error) {
	jsFilePath := VariationGlobalCodeJSFilePath(workingDir, accountID, campaignID, variationID)
	err := os.WriteFile(jsFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+jsFilePath)
	return jsFilePath, nil
}

func VariationGlobalCodeCSSFilePath(workingDir, accountID, campaignID, variationID string) string {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID
	variationCodeDir := campaignCodeDir + "/" + variationID

	err = os.MkdirAll(variationCodeDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	filePath := variationCodeDir + "/variationGlobalCode.css"
	return filePath
}

func WriteVariationGlobalCodeCSS(workingDir, accountID, campaignID, variationID, code string) (string, error) {
	cssFilePath := VariationGlobalCodeCSSFilePath(workingDir, accountID, campaignID, variationID)
	err := os.WriteFile(cssFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+cssFilePath)
	return cssFilePath, nil
}

func ModificationCodeFilePath(workingDir, accountID, campaignID, variationID, modificationID string) string {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID
	variationCodeDir := campaignCodeDir + "/" + variationID
	elementCodeDir := variationCodeDir + "/" + modificationID

	err = os.MkdirAll(elementCodeDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	jsFilePath := elementCodeDir + "/element.js"
	return jsFilePath
}

func WriteModificationCode(workingDir, accountID, campaignID, variationID, modificationID string, code []byte) (string, error) {
	cssFilePath := ModificationCodeFilePath(workingDir, accountID, campaignID, variationID, modificationID)
	err := os.WriteFile(cssFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+cssFilePath)
	return cssFilePath, nil
}

func DeleteModificationCodeDirectory(workingDir, accountID, campaignID, variationID, modificationID string) (string, error) {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		return "", err
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID
	variationCodeDir := campaignCodeDir + "/" + variationID
	elementCodeDir := variationCodeDir + "/" + modificationID

	if _, err := os.Stat(elementCodeDir); err == nil {
		err := os.RemoveAll(elementCodeDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error deleting file: ", err)
			return "", err
		}
	}

	return "modification directory for " + modificationID + " deleted", nil

}

func AddHeaderSelectorComment(selector string, code []byte, re *regexp.Regexp) []byte {

	if !re.Match(code) {
		selectorComment := "/* Selector: " + selector + " */\n"
		headerComment := []byte(selectorComment)

		fileCode := append(headerComment, []byte(code)...)
		return fileCode
	}

	return code
}

func HashFile(filepath string) ([32]byte, error) {
	var zero [32]byte

	f, err := os.Open(filepath)
	if err != nil {
		return zero, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return zero, fmt.Errorf("failed to compute hash: %w", err)
	}

	var result [32]byte
	copy(result[:], h.Sum(nil))
	return result, nil
}

func HashString(s string) [32]byte {
	h := sha256.New()
	h.Write([]byte(s))
	var result [32]byte
	copy(result[:], h.Sum(nil))
	return result
}
