package main

import (
	"github.com/codegangsta/cli"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"context"
	"log"
	"github.com/amoniacou/hubspot"
)

func main() {
	app := cli.NewApp()
	app.Name = "Hubspot To Google Spreadsheet"
	app.Usage = "Hubspot to google doc Contacts sync"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "hapikey",
			EnvVar: "HAPIKEY",
			Usage:  "Hubspot API Key",
		},
		cli.StringFlag{
			Name:   "googledoc_url",
			EnvVar: "GOOGLE_DOC_URL",
			Usage:  "Google Doc URL",
		},
	}
	app.Action = func(c *cli.Context) error {
		ctx := context.Background()
		client, err := google.DefaultClient(ctx, sheets.SpreadsheetsScope)
		if err != nil {
			log.Println("Something went wrong", err)
		}
		sheetsClient, err := sheets.New(client)
		if err != nil {
			log.Println("Something went wrong", err)
		}
		// Headers
		var vr sheets.ValueRange
		headerValues := []interface{}{"First Name", "Last Name", "Status", "Campaign", "Country", "Title"}
		hubspotValues := []string{
			"firstname", "lastname", "hs_lead_status", "lead_generation_campaign", "country", "jobtitle",
		}
		vr.Values = append(vr.Values, headerValues)
		hubspotClient := hubspot.NewHAPIClient(c.GlobalString("hapikey"))

		contacts, err := hubspotClient.GetContacts(0, 100, "all", hubspotValues)
		if err != nil {
			log.Fatal(err)
		}
		for {
			for _, contact := range contacts.Contacts {
				var tmp []interface{}
				for _, key := range hubspotValues {
					tmp = append(tmp, contact.Properties[key].Value)
				}
				vr.Values = append(vr.Values, tmp)
			}
			if contacts.Next() {
				contacts.GetNext()
			} else {
				break
			}
		}
		// Update Sheet
		_, err = sheetsClient.Spreadsheets.Values.Update(c.GlobalString("googledoc_url"), "A1", &vr).ValueInputOption("RAW").Do()
		if err != nil {
			log.Println("Unable to retrieve data from sheet. %v", err)
		}
		return nil
	}
	app.Run(os.Args)
}
