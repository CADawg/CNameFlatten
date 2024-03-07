package main

import (
	"context"
	"encoding/json"
	_ "github.com/joho/godotenv/autoload"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type MailInABoxDNSRecord struct {
	Source string `json:"qname"`
	Type   string `json:"rtype"`
	Target string `json:"value"`
	Zone   string `json:"zone"`
	Parent string `json:"-"`
}

var dnsResolver = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: time.Millisecond * time.Duration(10000),
		}
		conn, err := d.DialContext(ctx, network, "1.1.1.1:53")

		if err != nil {
			return d.DialContext(ctx, network, "8.8.8.8:53")
		}

		return conn, err
	},
}

// This is to keep track of the domains we've sent a PUT to as every subsequent request must be a POST instead if there are multiple IP addresses for each protocol
var mailInABoxAlreadyPutURLs []string

func main() {
	username := os.Getenv("MAILINABOX_USER")
	password := os.Getenv("MAILINABOX_PASSWORD")
	hostname := os.Getenv("MAILINABOX_HOSTNAME")

	// get the list of dns records (we need to look for _cname_flatten TXT records)
	res, err := GetRequestWithAuth(username, password, "https://"+hostname+"/admin/dns/custom")

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// parse the response
	var records []MailInABoxDNSRecord

	err = json.NewDecoder(res.Body).Decode(&records)

	if err != nil {
		panic(err)

	}

	var interestingRecords []MailInABoxDNSRecord

	// loop through the records and print out the ones we are interested in
	for _, record := range records {
		src := strings.Split(record.Source, ".")
		if record.Type == "TXT" && src[0] == "_cname_flatten" {
			// needs to be above 2 because if it's _cname_flatten.com then we can't set the A and AAAA records for .com
			if len(src) >= 2 {
				record.Parent = strings.Join(strings.Split(record.Source, ".")[1:], ".")
				interestingRecords = append(interestingRecords, record)
			}
		}
	}

	// update the A and AAAA records based on the TXT value's DNS records
	for _, record := range interestingRecords {
		// need to do a dns lookup
		records, err := dnsResolver.LookupIPAddr(context.Background(), record.Target)

		if err != nil {
			panic(err)
		}

		for _, ip := range records {
			isV4 := ip.IP.To4() != nil

			if isV4 {
				// update Whispering (A) records
				_, err := SetMailInABoxAnswer(username, password, "https://"+hostname+"/admin/dns/custom/"+record.Parent+"/A", ip.IP.String())

				if err != nil {
					panic(err)
				}
			} else {
				// update Screaming (AAAA) records
				_, err := SetMailInABoxAnswer(username, password, "https://"+hostname+"/admin/dns/custom/"+record.Parent+"/AAAA", ip.IP.String())
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func SetMailInABoxAnswer(username, password, url, answer string) (*http.Response, error) {
	client := &http.Client{}

	var method = http.MethodPut

	if contains(mailInABoxAlreadyPutURLs, url) {
		method = http.MethodPost
	}

	var answerReader = strings.NewReader(answer)

	req, err := http.NewRequest(method, url, answerReader)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	return client.Do(req)
}

func contains(ls []string, url string) bool {
	for _, l := range ls {
		if l == url {
			return true
		}
	}
	return false
}

func GetRequestWithAuth(username, password, url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	return client.Do(req)
}