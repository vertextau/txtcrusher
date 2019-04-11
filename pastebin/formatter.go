package pastebin

import (
	"encoding/xml"
	"fmt"
	"log"
)

type UserInfo struct {
	XMLName     xml.Name `xml:"user"`
	UsrName     string   `xml:"user_name"`
	UsrFormat   string   `xml:"user_format_short"`
	UsrExpr     string   `xml:"user_expiration"`
	UsrAvtrLink string   `xml:"user_avatar_url"`
	UsrPrvMode  string   `xml:"user_private"`
	UsrWebSite  string   `xml:"user_website"`
	UsrEmail    string   `xml:"user_email"`
	UsrLoc      string   `xml:"user_location"`
	UsrAccType  string   `xml:"user_account_type"`
}

func UnmarshalData(data *string) *UserInfo {
	v := UserInfo{}

	err := xml.Unmarshal([]byte(*data), &v)
	if err != nil {
		log.Fatal(err)
	}

	return &v
}

func PrintData(data *UserInfo) {
	fmt.Printf("Username:\t%q\n", data.UsrName)
	fmt.Printf("User format:\t%q\n", data.UsrFormat)
	fmt.Printf("Expiration:\t%q\n", data.UsrExpr)
	fmt.Printf("Avatar link:\t%q\n", data.UsrAvtrLink)
	fmt.Printf("Private mode:\t%q\n", data.UsrPrvMode)
	fmt.Printf("Website:\t%q\n", data.UsrWebSite)
	fmt.Printf("Email:\t\t%q\n", data.UsrEmail)
	fmt.Printf("Location:\t%q\n", data.UsrLoc)
	fmt.Printf("Account type:\t%q\n", data.UsrAccType)
}
