package main

import (
	"fmt"
	"regexp"
)

var replaceDataReg = regexp.MustCompile(`"/9j[A-Za-z0-9\+/=]+"`)
var replaceFeatureReg = regexp.MustCompile(`"Feature"(.+?)"(.+?)"`)

func main() {
	str := `"Data": "/9j/4AAQSkZJRgABAQAAAQABAAD/2wDFABALDA4MChAODQ4SERATGCgaGBYWGDEjJR0oOjM9PDkzODdAY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2N+mRJbiARSuJMfdYjkU0SSN/q4SfwzR5V22Moye5OKAGRGewfbKhaBu4521fBBAKkFSMg+tUnsnmQpLcsqnrs61PbW0VpAIYixVe7HJoAmooo70AGaQ/NnNHel4FADo5OiP17GpKgPIp6SfwOeexoAkopKWgAooooAKM0UlAC0UUlAC0UlFAC0lFGaACikozQAtJRmj86AP/9k="`
	// reqStr := string(content)
	str = data

	str = replaceDataReg.ReplaceAllString(str, `"/9j/ddd"`)
	fmt.Println(str)

	str = replaceFeatureReg.ReplaceAllString(str, `"Feature":"PRdx/w+rvQ=="`)
	fmt.Println(str)
}

var data = `{
	"Reference": "",
	"Seq": 9,
	"DeviceCode": "210235C3R0320B000985",
	"Timestamp": 1624364137,
	"NotificationType": 0,
	"FaceInfoNum": 1,
	"FaceInfoList": [
		{
			"ID": 9,
			"Timestamp": 1624364136,
			"CapSrc": 1,
			"FeatureNum": 0,
			"FeatureList": [
				{
					"FeatureVersion": "",
					"Feature": "QAAAQABAAD/2wDFABsSFBcUERsXFhceHB"
				},
				{
					"FeatureVersion": "",
					"Feature": "QAAAQABAAD/2wDFABsSFBcUERsXFhceHB"
				}
			],
			"Temperature": 0,
			"MaskFlag": 0,
			"PanoImage": {
				"Name": "1624364136_1_93.jpg",
				"Size": 134884,
				"Data": "/9j/4AAQSkZJRgABAQAAAQABAAD/2wDFABALDA4MChAODQ4SERATGCgaGBYWGDEjJR0oOjM9PDkzODdAY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2N+mRJbiARSuJMfdYjkU0SSN/q4SfwzR5V22Moye5OKAGRGewfbKhaBu4521fBBAKkFSMg+tUnsnmQpLcsqnrs61PbW0VpAIYixVe7HJoAmooo70AGaQ/NnNHel4FADo5OiP17GpKgPIp6SfwOeexoAkopKWgAooooAKM0UlAC0UUlAC0UlFAC0lFGaACikozQAtJRmj86AP/9k="
			},
			"FaceImage": {
				"Name": "1624364136_2_93.jpg",
				"Size": 12884,
				"Data": "/9j/4AAQSkZJRgABAQAAAQABAAD/2wDFABsSFBcUERsXFhceHBsgKEIrKCUlKFE6PTBCYFVlZF9VXVtqeJmBanGQc1tdhbWGkJ6jq62rZ4C8ybqmx5moq6QBHB4eKCMoTisrTqRuXW6kpKSkNMBp2aAHhSelG0+lCOF61J50fv+VICIqQM0bSRnFSNJGy4BpQ8YGA1AEe00bimAtOFIKWgBaVVLHAoRS5woq5FEEHvQKTikzQAtGaSnUwCiilpAFLRRQAtKKSlpgLRRRQAUopKKQC0006m0AKKdTRTqACkpCaKAFpabSimA4UtIKWgBaKSloEFFFLQAlBOKCcVGxyaAHb/ejf71HtJpdpoA//9k="
			},
			"FaceArea": {
				"LeftTopX": 2555,
				"LeftTopY": 5390,
				"RightBottomX": 4620,
				"RightBottomY": 6531
			}
		}
	],
	"CardInfoNum": 0,
	"CardInfoList": [],
	"GateInfoNum": 0,
	"GateInfoList": [],
	"LibMatInfoNum": 1,
	"LibMatInfoList": [
		{
			"ID": 9,
			"LibID": 3,
			"LibType": 3,
			"MatchStatus": 1,
			"MatchPersonID": 111,
			"MatchFaceID": 1,
			"MatchPersonInfo": {
				"PersonCode": "111",
				"PersonName": "住户1",
				"Gender": 0,
				"CardID": "",
				"IdentityNo": ""
			}
		}
	]
}`
