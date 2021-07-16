package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

func main() {
	var str = `POST /LAPI/V1.0/PACS/Controller/HeartReportInfo HTTP/1.1\r\nHost: 82.157.107.86\r\nContent-Type: application/json\r\nContent-Length: 180\r\n\r\n{\r\n\"RefId\": \"8eb4f13c-59ed-4dc6-9074-898d8408c97d\",\r\n\"Time\": \"2021-07-06 03:26:02\",\r\n\"NextTime\": \"2021-07-06 03:26:12\",\r\n\"DeviceCode\": \"210235C3R0320B000985\",\r\n\"DeviceType\": 5\r\n}\r\n`

	str = `GET http://82.157.107.86:8080/api/gb28181/invite?id=34020000001330000002&channel=34020000001310000003 HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: 82.157.107.86:8080
Connection: close
User-Agent: Paw/3.2.2 (Macintosh; OS X/11.4.0) GCDHTTPRequest

`

	str = `POST /VIID/Faces HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: 127.0.0.1:8085
Connection: close
User-Agent: Paw/3.2.2 (Macintosh; OS X/11.4.0) GCDHTTPRequest
Content-Length: 1171

{"FaceListObject":{"FaceObject":[{"FaceID":"111110220200710143217001770600178","InfoKind":1,"SourceID":"11111022020071014321700177","DeviceID":"11111","ShotTime":"20200710143217","LeftTopX":512,"LeftTopY":369,"RightBtmX":749,"RightBtmY":707,"LocationMarkTime":"20200710143217","FaceAppearTime":"20200710143217","FaceDisAppearTime":"20200710143217","GenderCode":"1","AgeUpLimit":28,"AgeLowerLimit":28,"GlassStyle":"99","Emotion":"1","IsDriver":2,"IsForeigner":2,"IsSuspectedTerrorist":2,"IsCriminalInvolved":2,"IsDetainees":2,"IsVictim":2,"IsSuspiciousPerson":2,"Similaritydegree":0,"SubImageList":{"SubImageInfoObject":[{"ImageID":"11111022020071014321700177","EventSort":10,"DeviceID":"11111","StoragePath":"","Type":"14","FileFormat":"Jpeg","ShotTime":"20200710143217","Width":1920,"Height":1264,"Data":"\u56fe\u7247\u6570\u636e"},{"ImageID":"11111022020071014321700180","EventSort":10,"DeviceID":"11111","StoragePath":"","Type":"11","FileFormat":"Jpeg","ShotTime":"20200710143217","Width":896,"Height":700,"Data":"\u56fe\u7247\u6570\u636e"}]},"RelatedType":"01","RelatedList":{"RelatedObject":[{"RelatedType":"01","RelatedID":"111110220200710143217001770100179"}]}}]}}`

	r := bytes.NewReader([]byte(str))
	reader := bufio.NewReader(r)
	req, err := http.ReadRequest(reader)
	fmt.Println(err)

	body, err := ioutil.ReadAll(req.Body)
	bodystr := string(body)
	fmt.Println(bodystr)

}
