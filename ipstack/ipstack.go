package ipstack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
)

type ipStackResponse struct {
	Success     bool   `json:"success"`
	Ip          net.IP `json:"ip"`
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
	PostalCode  string `json:"zip"`
	Location    ipStackLocation
}

type ipStackLocation struct {
	CountryFlagEmoji string `json:"country_flag_emoji"`
}

type Location struct {
	CountryName string
	RegionName  string
	City        string
	PostalCode  string
	CountryFlag string
}

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

var HTTPClient httpClient

func init() {
	HTTPClient = &http.Client{}
}

func GetLocation(ip string) (Location, error) {
	accessKey := os.Getenv("IPSTACK_ACCESS_KEY")
	baseUrl := "http://api.ipstack.com/"

	verifiedIp := net.ParseIP(ip)
	if verifiedIp == nil {
		return Location{}, errors.New("invalid IP")
	}

	url, err := url.Parse(fmt.Sprintf("%s%s", baseUrl, verifiedIp.String()))
	if err != nil {
		return Location{}, err
	}
	q := url.Query()
	q.Set("access_key", accessKey)
	url.RawQuery = q.Encode()

	res, err := HTTPClient.Get(url.String())
	if err != nil {
		return Location{}, err
	}
	defer res.Body.Close()

	ipStackResponse := ipStackResponse{Success: true}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Location{}, err
	}
	if err = json.Unmarshal(body, &ipStackResponse); err != nil {
		return Location{}, err
	}
	if !ipStackResponse.Success {
		return Location{}, errors.New("IP Stack call failed")
	}

	response := Location{
		CountryName: ipStackResponse.CountryName,
		RegionName:  ipStackResponse.RegionName,
		City:        ipStackResponse.City,
		PostalCode:  ipStackResponse.PostalCode,
		CountryFlag: ipStackResponse.Location.CountryFlagEmoji,
	}
	return response, nil
}
