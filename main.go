package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/vertextau/txtcrusher/pastebin"
)

var (
	createPaste    = flag.String("c", "", "Create a new paste")
	createFromFile = flag.String("f", "", "Create a new paste with the content from an input file")
	createKey      = flag.Bool("k", false, "Create a new user key")
	listPastes     = flag.Int("l", 0, "List pastes created by a user")
	deletePaste    = flag.String("d", "", "Delete a paste")
	echoInfo       = flag.Bool("i", false, "Get a user info")
	getUserPaste   = flag.String("g", "", "Get a paste created by a user")
	grabPaste      = flag.String("p", "", "Grab a paste without a config file")
	helpFlag       = flag.Bool("help", false, "Manual about the program")

	// flags for a new paste
	guestFlag   = flag.Bool("guest", false, "Create a new paste under a guest")
	titleFlag   = flag.String("title", "", "Title for a new paste")
	formatFlag  = flag.String("format", "Text", "Highlight format for a new paste")
	expDateFlag = flag.String("expire", "N", "Expire date for a new paste")
	modFlag     = flag.Int("mod", 0, "Visibility modificator for a new paste")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: txtcrusher [OPTION] INPUT\nTry 'txtcrusher -help' for more information.")
		os.Exit(1)
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if len(*grabPaste) > 0 {
		data, err := pastebin.GetPaste(*grabPaste)
		if err != nil {
			log.Fatalf("[%s] %s", *grabPaste, err)
		}
		fmt.Println(*data)
		os.Exit(0)
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
		if err != nil {
			log.Fatalf("[%s] %s", *createPaste, err)
		}

		fmt.Println(*data)
	} else if len(*createFromFile) > 0 {
		fileData, err := ioutil.ReadFile(*createFromFile)
		if err != nil {
			log.Fatalf("[%s] %s", *createFromFile, err)
		}

		text := string(fileData)

		data, err := u.CreateNewPaste(&text, *guestFlag, *titleFlag, *formatFlag, *expDateFlag, *modFlag)
		if err != nil {
			log.Fatalf("[%s] %s", *createFromFile, err)
		}

		fmt.Println(*data)
	} else if *createKey {
		if len(flag.Args()) != 2 {
			log.Fatal("Usage: txtcrusher -k username password")
		}

		data, err := u.GetUserKey(flag.Arg(0), flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(*data)
	} else {
		switch {
		case len(*deletePaste) > 0:
			data, err := u.DeleteUserPaste(*deletePaste)
			if err != nil {
				log.Fatalf("[%s] %s", *deletePaste, err)
			}

			fmt.Println(*data)

		case *echoInfo:
			data, err := u.GetUserInfo()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(*data)

		case len(*getUserPaste) > 0:
			data, err := u.GetUserPaste(*getUserPaste)
			if err != nil {
				log.Fatalf("[%s] %s", *getUserPaste, err)
			}

			fmt.Println(*data)

		case *listPastes > 0:
			data, err := u.ListUserPastes(*listPastes)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(*data)
		}
	}
}
