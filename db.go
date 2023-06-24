package main

import (
	"encoding/csv"
	"net"
	"os"
)

type IPRange struct {
	IpRange     *net.IPNet
	CountryCode string
}

type rawCountryData struct {
	CountryCode string
	Id          string
}

func readCountryData() []*rawCountryData {
	var countryDataList []*rawCountryData
	f, err := os.Open("GeoLite2-Country-Locations-en.csv")
	panicIfError(err)

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	panicIfError(err)

	for _, row := range rows {
		countryDataList = append(countryDataList, &rawCountryData{
			CountryCode: row[2],
			Id:          row[0],
		})
	}

	return countryDataList
}

func findCountryCode(countryId string, rawCountryDataList []*rawCountryData) string {
	for _, data := range rawCountryDataList {
		if data.Id == countryId {
			return data.CountryCode
		}
	}
	return ""
}

func CreateNetworkDb() []*IPRange {
	countryData := readCountryData()
	f, err := os.Open("GeoLite2-Country-Blocks-IPv4.csv")
	panicIfError(err)
	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	panicIfError(err)

	ipRangeList := make([]*IPRange, len(rows)-1) // 減掉標頭列
	i := 0
	for _, row := range rows {
		_, ipRange, err := net.ParseCIDR(row[0])
		if err != nil && err.Error() == "invalid CIDR address: network" {
			continue
		}
		panicIfError(err)

		ipRangeList[i] = &IPRange{
			IpRange:     ipRange,
			CountryCode: findCountryCode(row[1], countryData),
		}
		i++
	}

	return ipRangeList
}
