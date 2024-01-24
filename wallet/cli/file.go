package cli

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"goXdagWallet/components"
	xdagoUtils "goXdagWallet/xdago/utils"
	"os"
)

var validateFile = func(input string) error {
	if !xdagoUtils.FileExists(input) {
		return errors.New("input path not exists")
	}

	fileInfo, _ := os.Stat(input)
	if fileInfo.IsDir() {
		return errors.New("path to import file error")
	} else {
		return nil
	}
}

var validateFolder = func(input string) error {
	if !xdagoUtils.FileExists(input) {
		return errors.New("input path not exists")
	}
	fileInfo, _ := os.Stat(input)
	if fileInfo.IsDir() {
		if components.CheckOldWallet(input) {
			return nil
		} else {
			return errors.New("no wallet data in the folder")
		}
	} else {
		return errors.New("path to import folder error")
	}
}

func InputFilePath() string {
	prompt := promptui.Prompt{
		Label:    "Path to Mnemonic text file",
		Validate: validateFile,
	}

	result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input path failed %v\n", err)
		result, err = prompt.Run()
	}
	return result

}

func inputFolderPath() string {
	prompt := promptui.Prompt{
		Label:    "Path to Non Mnemonic data folder",
		Validate: validateFolder,
	}

	result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input path failed %v\n", err)
		result, err = prompt.Run()
	}
	return result
}
