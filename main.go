package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/vertextau/txtcrusher/pastebin"
	"io/ioutil"
)

var (
	createPaste    = flag.String("c", "", "Create a new paste")
	createFromFile = flag.String("f", "", "Create a new paste with the content from input file")
	createKey      = flag.Bool("k", false, "Create a new user key")
	listPastes     = flag.Int("l", 0, "List pastes created by a user")
	deletePaste    = flag.String("d", "", "Delete a paste")
	echoInfo       = flag.Bool("i", false, "Get a user info")
	getUserPaste   = flag.String("g", "", "Get a paste created by a user")
	grabPaste      = flag.String("p", "", "Grab a paste without a config file")
	helpFlag       = flag.Bool("help", false, "Manual about the program")

	// flags for a new paste
	guestFlag   = flag.Bool("guest", false, "Post under a guest")
	titleFlag   = flag.String("title", "", "Title for a paste")
	formatFlag  = flag.String("format", "Text", "Highlight format")
	expDateFlag = flag.String("expire", "N", "Expire date for a paste")
	modFlag     = flag.Int("mod", 0, "Modificator for a paste")
)

func main() {
	flag.Parse()

	if len(*grabPaste) > 0 {
		data, err := pastebin.GetPaste(*grabPaste)
		checkError(err)
		fmt.Println(*data)
		return
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/txtcrusher")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error: %s", err)
	}

	u := pastebin.Pastebin{}
	u.DeveloperKey = viper.GetString("pastebin.api_dev_key")
	u.UserKey = viper.GetString("pastebin.api_user_key")

	if len(*createPaste) > 0 {
		data, err := u.CreateNewPaste(createPaste, *guestFlag, *titleFlag, *formatFlag, *expDateFlag, *modFlag)
		checkError(err)
		fmt.Println(*data)
	} else if len(*createFromFile) > 0 {
		fileData, err := ioutil.ReadFile(*createFromFile)
		checkError(err)
		text := string(fileData)
		data, err := u.CreateNewPaste(&text, *guestFlag, *titleFlag, *formatFlag, *expDateFlag, *modFlag)
		checkError(err)
		fmt.Println(*data)
	}

	switch {
	case *createKey:
		data, err := u.GetUserKey(flag.Arg(0), flag.Arg(1))
		checkError(err)
		fmt.Println(*data)
	case len(*deletePaste) > 0:
		data, err := u.DeleteUserPaste(*deletePaste)
		if err != nil {
			log.Fatalf("[%s] %s", *deletePaste, err)
		}
		fmt.Println(*data)
	case *echoInfo:
		data, err := u.GetUserInfo()
		checkError(err)
		fmt.Println(*data)
	case len(*getUserPaste) > 0:
		data, err := u.GetUserPaste(*getUserPaste)
		if err != nil {
			log.Fatalf("[%s] %s", *getUserPaste, err)
		}
		fmt.Println(*data)
	case *listPastes > 0:
		data, err := u.ListUserPastes(*listPastes)
		checkError(err)
		fmt.Println(*data)
	case *helpFlag:
		flag.Usage()
	default:
		fmt.Println()
		//fmt.Println("Usage: txtcrusher [OPTION] INPUT\nTry 'txtcrusher -help' for more information.")
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
