package main

import (
	"encoding/csv"
	"net/netip"
	"os"
)

type IPAndCountryMappingDatabase struct {
	Ipv4 []*IPRangeData
	Ipv6 []*IPRangeData
}

type IPRangeData struct {
	IpRange     netip.Prefix
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
			CountryCode: row[4],
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

func buildIpv4Data(countryData *[]*rawCountryData) []*IPRangeData {
	f, err := os.Open("GeoLite2-Country-Blocks-IPv4.csv")
	panicIfError(err)
	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	panicIfError(err)

	ipRangeList := make([]*IPRangeData, len(rows)-1) // 減掉標頭列
	i := 0
	for _, row := range rows {
		ipRange, err := netip.ParsePrefix(row[0])
		if err != nil && err.Error() == "netip.ParsePrefix(\"network\"): no '/'" {
			continue
		}
		panicIfError(err)

		ipRangeList[i] = &IPRangeData{
			IpRange:     ipRange,
			CountryCode: findCountryCode(row[1], *countryData),
		}
		i++
	}

	return ipRangeList
}

func buildIpv6Data(countryData *[]*rawCountryData) []*IPRangeData {
	f, err := os.Open("GeoLite2-Country-Blocks-IPv6.csv")
	panicIfError(err)
	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	panicIfError(err)

	ipRangeList := make([]*IPRangeData, len(rows)-1) // 減掉標頭列
	i := 0
	for _, row := range rows {
		ipRange, err := netip.ParsePrefix(row[0])
		if err != nil && err.Error() == "netip.ParsePrefix(\"network\"): no '/'" {
			continue
		}
		panicIfError(err)

		ipRangeList[i] = &IPRangeData{
			IpRange:     ipRange,
			CountryCode: findCountryCode(row[1], *countryData),
		}
		i++
	}

	return ipRangeList
}

func BuildDatabase() *IPAndCountryMappingDatabase {
	countryData := readCountryData()
	Ipv4 := buildIpv4Data(&countryData)
	Ipv6 := buildIpv6Data(&countryData)
	return &IPAndCountryMappingDatabase{
		Ipv4,
		Ipv6,
	}
}
