package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/vertextau/txtcrusher/pastebin"
)

var (
	createPaste  = flag.String("c", "", "Create a new paste")
	createKey    = flag.Bool("k", false, "Create a new user key")
	listPastes   = flag.Int("l", 50, "List pastes created by a user")
	deletePaste  = flag.String("d", "", "Delete a paste")
	echoInfo     = flag.Bool("i", false, "Get a user info")
	getUserPaste = flag.String("g", "", "Get a paste created by a user")
	grabPaste    = flag.String("p", "", "Grab a paste without a config file")

	// flags for a new paste
	guestFlag   = flag.Bool("guest", false, "Post under a guest")
	titleFlag   = flag.String("title", "", "Title for a paste")
	formatFlag  = flag.String("format", "Text", "Highlight format")
	expDateFlag = flag.String("expire", "N", "Expire date for a paste")
	modFlag     = flag.Int("mod", 0, "Modificator for a paste")
)

func chck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if len(*grabPaste) > 0 {
		data, err := pastebin.GetPaste(*grabPaste)
		chck(err)
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

	switch {
	case len(*createPaste) > 0:
		data, err := u.CreateNewPaste(createPaste, *guestFlag, *titleFlag, *formatFlag, *expDateFlag, *modFlag)
		chck(err)
		fmt.Println(*data)
	case *createKey:
		data, err := u.GetUserKey(flag.Arg(0), flag.Arg(1))
		chck(err)
		fmt.Println(*data)
	case len(*deletePaste) > 0:
		data, err := u.DeleteUserPaste(*deletePaste)
		chck(err)
		fmt.Println(*data)
	case *echoInfo:
		data, err := u.GetUserInfo()
		chck(err)
		fmt.Println(*data)
	case len(*getUserPaste) > 0:
		data, err := u.GetUserPaste(*getUserPaste)
		chck(err)
		fmt.Println(*data)
	case *listPastes > 0:
		data, err := u.ListUserPastes(*listPastes)
		chck(err)
		fmt.Println(*data)
	}
}
