package config

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
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

func AccountGlobalCodeFilePath(workingDir, accountID string) (string, error) {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		return "", err
	}

	accountCodeDir := gcWorkingDir + "/" + accountID

	err = os.MkdirAll(accountCodeDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	jsFilePath := accountCodeDir + "/accountGlobalCode.js"
	return jsFilePath, nil
}

func WriteAccountGlobalCode(workingDir, accountID, code string) (string, error) {
	jsFilePath, err := AccountGlobalCodeFilePath(workingDir, accountID)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(jsFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+jsFilePath)
	return jsFilePath, nil
}

func CampaignGlobalCodeFilePath(workingDir, accountID, campaignID string) (string, error) {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		return "", err
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID

	err = os.MkdirAll(campaignCodeDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	jsFilePath := campaignCodeDir + "/campaignGlobalCode.js"
	return jsFilePath, nil
}

func WriteCampaignGlobalCode(workingDir, accountID, campaignID, code string) (string, error) {
	jsFilePath, err := CampaignGlobalCodeFilePath(workingDir, accountID, campaignID)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(jsFilePath, []byte(code), os.ModePerm)
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

func VariationGlobalCodeJSFilePath(workingDir, accountID, campaignID, variationID string) (string, error) {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		return "", err
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID
	variationCodeDir := campaignCodeDir + "/" + variationID

	err = os.MkdirAll(variationCodeDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	jsFilePath := variationCodeDir + "/variationGlobalCode.js"
	return jsFilePath, nil
}

func WriteVariationGlobalCodeJS(workingDir, accountID, campaignID, variationID, code string) (string, error) {
	jsFilePath, err := VariationGlobalCodeJSFilePath(workingDir, accountID, campaignID, variationID)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(jsFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+jsFilePath)
	return jsFilePath, nil
}

func VariationGlobalCodeCSSFilePath(workingDir, accountID, campaignID, variationID string) (string, error) {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		return "", err
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID
	variationCodeDir := campaignCodeDir + "/" + variationID

	err = os.MkdirAll(variationCodeDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	filePath := variationCodeDir + "/variationGlobalCode.css"
	return filePath, nil
}

func WriteVariationGlobalCodeCSS(workingDir, accountID, campaignID, variationID, code string) (string, error) {
	cssFilePath, err := VariationGlobalCodeCSSFilePath(workingDir, accountID, campaignID, variationID)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(cssFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+cssFilePath)
	return cssFilePath, nil
}

func ModificationCodeFilePath(workingDir, accountID, campaignID, variationID, modificationID string) (string, error) {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		return "", err
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID
	variationCodeDir := campaignCodeDir + "/" + variationID
	elementCodeDir := variationCodeDir + "/" + modificationID

	err = os.MkdirAll(elementCodeDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	jsFilePath := elementCodeDir + "/element.js"
	return jsFilePath, nil
}

func WriteModificationCode(workingDir, accountID, campaignID, variationID, modificationID string, code []byte) (string, error) {
	cssFilePath, err := ModificationCodeFilePath(workingDir, accountID, campaignID, variationID, modificationID)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(cssFilePath, []byte(code), os.ModePerm)
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
