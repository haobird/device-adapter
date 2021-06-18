package faceguard

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/haobird/goutils"
	"github.com/tidwall/gjson"
)

//Image 图片
type Image struct {
	Data string  `json:"Data"`
	Name string  `json:"Name"`
	Size float64 `json:"Size"`
}

//CardInfo 卡片
type CardInfo struct {
	Birthday   string  `json:"Birthday"`
	CapSrc     int64   `json:"CapSrc"`
	CardID     string  `json:"CardID"`
	CardStatus int64   `json:"CardStatus"`
	CardType   int64   `json:"CardType"`
	Gender     int64   `json:"Gender"`
	Size       float64 `json:"Size"`
	ID         int64   `json:"ID"`
	IDImage    Image   `json:"IDImage"`
	IdentityNo string  `json:"IdentityNo"`
	Name       string  `json:"Name"`
}

//FaceInfo 图片信息
type FaceInfo struct {
	CapSrc   int64 `json:"CapSrc"`
	FaceArea struct {
		LeftTopX     float64 `json:"LeftTopX"`
		LeftTopY     float64 `json:"LeftTopY"`
		RightBottomX float64 `json:"RightBottomX"`
		RightBottomY float64 `json:"RightBottomY"`
	} `json:"FaceArea"`
	FaceImage   Image `json:"FaceImage"`
	FeatureList []struct {
		Feature        string `json:"Feature"`
		FeatureVersion string `json:"FeatureVersion"`
	} `json:"FeatureList"`
	FeatureNum  int64   `json:"FeatureNum"`
	ID          int64   `json:"ID"`
	MaskFlag    int64   `json:"MaskFlag"`
	PanoImage   Image   `json:"PanoImage"`
	Temperature float64 `json:"Temperature"`
	Timestamp   int64   `json:"Timestamp"`
}

//PersonInfo 个人信息
type PersonInfo struct {
	CardID     string `json:"CardID"`
	Gender     int64  `json:"Gender"`
	IdentityNo string `json:"IdentityNo"`
	PersonCode string `json:"PersonCode"`
	PersonName string `json:"PersonName"`
}

//LibMatInfo 库匹配信息
type LibMatInfo struct {
	ID              int64      `json:"ID"` // 记录 ID
	LibID           int64      `json:"LibID"`
	LibType         int64      `json:"LibType"`
	MatchFaceID     int64      `json:"MatchFaceID"`
	MatchPersonID   int64      `json:"MatchPersonID"`
	MatchPersonInfo PersonInfo `json:"MatchPersonInfo"`
	MatchStatus     int64      `json:"MatchStatus"`
}

//GateInfo 门的信息
type GateInfo struct {
}

//PersonVerifyInfo 校验信息
type PersonVerifyInfo struct {
	NotificationType int64        `json:"NotificationType"`
	Reference        string       `json:"Reference"`
	Seq              int64        `json:"Seq"`
	Timestamp        int64        `json:"Timestamp"`
	DeviceCode       string       `json:"DeviceCode"`
	CardInfoNum      int64        `json:"CardInfoNum"`
	FaceInfoNum      int64        `json:"FaceInfoNum"`
	LibMatInfoNum    int64        `json:"LibMatInfoNum"`
	GateInfoNum      int64        `json:"GateInfoNum"`
	GateInfoList     []GateInfo   `json:"GateInfoList"`
	LibMatInfoList   []LibMatInfo `json:"LibMatInfoList"`
	CardInfoList     []CardInfo   `json:"CardInfoList"`
	FaceInfoList     []FaceInfo   `json:"FaceInfoList"`
}

type Sdk struct{}

//_parse 解析包  topic payload
func (this *Sdk) _parse(buf []byte) (string, string, error) {
	eventType := "unknown"
	// 分割字符串
	str := string(buf)
	pos := strings.Index(str, "{")
	if pos < 1 {
		fmt.Println("未读取数据")
		return eventType, str, errors.New("未读取数据")
	}
	var input map[string]interface{}
	content := str[pos:]
	err := json.Unmarshal([]byte(content), &input)
	if err != nil {
		return eventType, content, errors.New("解析数据错误")
	}
	if strings.Contains(str, "HeartReportInfo") {
		// 处理心跳
		eventType = "heart"
	} else if strings.Contains(str, "PersonVerification") && !strings.Contains(str, "Notifications") {
		eventType = "open"
	} else if strings.Contains(str, "Response") {
		// 暂时不处理
		eventType = "response"
	}
	return eventType, content, nil
}

//ReadPacket 处理 消息, 返回 messageType, requestId , topic , payload
func (this *Sdk) ParsePacket(buf []byte) (string, string, string, []byte, error) {
	topic := ""
	messageType := Unknown
	var requestid string = ""

	// 分割字符串
	str := string(buf)
	pos := strings.Index(str, "{")
	if pos < 1 {
		fmt.Println("未读取数据")
		return messageType, requestid, topic, nil, errors.New("未读取数据")
	}
	var input map[string]interface{}
	var resp []byte
	content := str[pos:]
	err := json.Unmarshal([]byte(content), &input)
	if err != nil {
		return messageType, requestid, topic, nil, errors.New("解析数据错误")
	}

	// 提取 requestid
	requestID := between(str, "RequestID: ", "\r\n")

	if strings.Contains(str, "HeartReportInfo") {
		return Heart, requestid, topic, nil, nil
	} else if strings.Contains(str, "PersonVerification") && !strings.Contains(str, "Notifications") {
		messageType = Publish
		topic = "personVerification"

		var info PersonVerifyInfo
		err = json.Unmarshal([]byte(content), &info)
		if err == nil {
			respInfo := map[string]interface{}{
				"persionID":   info.LibMatInfoList[0].MatchPersonID,
				"persionName": info.LibMatInfoList[0].MatchPersonInfo.PersonName,
				"openType":    "1",
				"panoImage":   info.FaceInfoList[0].PanoImage,
				"faceImage":   info.FaceInfoList[0].FaceImage,
				"timestamp":   info.Timestamp,
			}
			resp, _ = json.Marshal(respInfo)
		}
	} else if strings.Contains(str, "DeviceBasicInfo") {
		messageType = Publish
		topic = "deviceInfoFace"

		fmt.Println("DeviceBasicInfo:", input)

		respInfo := map[string]interface{}{
			"mac":        gjson.Get(content, "Response.Data.MAC").String(),
			"ip":         gjson.Get(content, "Response.Data.Address").String(),
			"deviceCode": gjson.Get(content, "Response.Data.SerialNumber").String(),
			"name":       gjson.Get(content, "Response.Data.DeviceModel").String(),
		}
		resp, _ = json.Marshal(respInfo)

	} else if strings.Contains(str, "Response") {
		messageType = PubAck
		resp, _ = json.Marshal(input)
	}

	return messageType, requestID, topic, resp, nil
}

//ParseConnect 解析连接信息
func (this *Sdk) ParseConnect(buf []byte) (string, error) {
	eventType, content, err := this._parse(buf)
	if err != nil {
		return "", err
	}
	if eventType != "heart" {
		return "", errors.New("不正确的数据包:" + eventType)
	}
	var input map[string]interface{}
	json.Unmarshal([]byte(content), &input)
	if deviceCode, ok := input["DeviceCode"]; ok {
		clientID := goutils.String(deviceCode)
		return clientID, nil
	}
	return "", errors.New("设备编号不存在")
}

// Handle 传入消息包，返回: 消息类型， 消息id，topic, payload, error : 如果 是 publish 、 puback, 则进行其它处理
func (this *Sdk) HandlePacket(buf []byte) (*Package, error) {
	messageType, requestid, topic, body, err := this.ParsePacket(buf)
	if err != nil {
		return nil, err
	}

	pack := &Package{
		MessageType: messageType,
		RequestID:   requestid,
	}

	if messageType == Heart {
		return pack, nil
	}

	if messageType == PubAck {
		// 处理 提取 关键 信息
		pack.Topic = topic
		pack.Data = body
		return pack, nil

	}

	if messageType == Publish {
		pack.Topic = topic
		pack.Data = body
		return pack, nil
	}

	return pack, errors.New("not match  messageType")
}

// 获取设备信息
func (this *Sdk) PackDeviceInfoReq() string {
	header := map[string]string{"RequestID": goutils.String(time.Now().Unix())}
	return this.PackReq("GET", "/LAPI/V1.0/System/DeviceBasicInfo", "", header)
}

//RespPing 封装相应的响应
func (this *Sdk) Pong() string {
	str := `{"ResponseURL": "/LAPI/V1.0/PACS/Controller/HeartReportInfo", "Code": 0,"Data": {"Time": "%s" }}`
	str = fmt.Sprintf(str, goutils.GetNormalTimeString(time.Now()))
	return this.PackResp(str, nil)
}

func (this *Sdk) PackResp(body string, headers map[string]string) string {
	respHeader := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain;charset=ISO-8859-1\r\n"
	for key, val := range headers {
		respHeader = respHeader + key + ": " + val + "\r\n"
	}
	//拼装http返回的消息
	resp := respHeader + fmt.Sprintf("Content-Length: %d\r\n", len(body)) + "\r\n" + body
	return resp
}

func (this *Sdk) PackReq(method, url, body string, headers map[string]string) string {
	strHeader := method + " " + url + " HTTP/1.1\r\n" +
		"Content-Type: application/json\r\n" +
		"Connection: close\r\n"
	for key, val := range headers {
		strHeader = strHeader + key + ": " + val + "\r\n"
	}

	text := strHeader
	length := len(body)
	if length > 0 {
		text = text + fmt.Sprintf("Content-Length: %d\r\n", len(body))
	}
	text = text + "\r\n" + body

	return text
}

//Trans 处理 对应的 数据 提取
func (this *Sdk) Trans(pack Package) []byte {
	topic := pack.Topic
	// 设备白名单
	if topic == "personAuthorized" {
		body := this.dealAuthorized(pack.Data)
		// body := string(pack.Data)
		header := map[string]string{"RequestID": pack.RequestID}
		content := this.PackReq("POST", "/LAPI/V1.0/PeopleLibraries/3/People", body, header)
		// fmt.Println(content)
		return []byte(content)
	}

	// 删除白名单
	if topic == "personAuthorizedCancel" {
		header := map[string]string{"RequestID": pack.RequestID}
		persionID := gjson.GetBytes(pack.Data, "persionID").String()
		timestamp := time.Now().Unix()
		url := fmt.Sprintf("/LAPI/V1.0/PeopleLibraries/3/People/%s?LastChange=%d", persionID, timestamp)
		content := this.PackReq("DELETE", url, "", header)
		// fmt.Println(content)
		return []byte(content)
	}

	if topic == "personSearch" {
		body := this.personSearch(pack.Data)
		header := map[string]string{"RequestID": pack.RequestID}
		content := this.PackReq("POST", "/LAPI/V1.0/PeopleLibraries/3/People/Info", body, header)
		return []byte(content)
	}

	// 远程开门
	if topic == "deviceOpenFace" {
		header := map[string]string{"RequestID": pack.RequestID}
		content := this.PackReq("PUT", "/LAPI/V1.0/PACS/Controller/RemoteOpened", "", header)
		return []byte(content)
	}

	// 设备信息
	if topic == "deviceInfoFace" {
		header := map[string]string{"RequestID": pack.RequestID}
		content := this.PackReq("GET", "/LAPI/V1.0/System/DeviceBasicInfo", "", header)
		return []byte(content)
	}

	return nil
}

func (this *Sdk) personSearch(buf []byte) string {
	str := `{ "Num": 0, "QueryInfos": [ {"QryType": 27, "QryCondition": 0, "QryData": "1001" },{"QryType": 55, "QryCondition": 0, "QryData": "Uniview" }],"Limit": 10, "Offset": 0 }`
	return str
}

func (this *Sdk) dealAuthorized3(buf []byte) string {
	str := ""
	str = `{"Num": 1,"PersonInfoList": [{"PersonID": 22,"LastChange": 1602329484,"PersonCode": "5hh","PersonName": "我的陌生人","Remarks": "陌生人的尝试哈哈哈哈哈嘎","TimeTemplateNum": 0,"ImageNum": 1,"ImageList": [{"FaceID": 1,"Name": "1_1.jpg","Size": 3196,"Data": "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"}]}]}`
	str = `{"Num": 1,"PersonInfoList": [{"PersonID": 22,"LastChange": 1602329484,"PersonCode": "5hh","PersonName": "我的陌生人","Remarks": "陌生人的尝试哈哈哈哈哈嘎","TimeTemplateNum": 0,"ImageNum": 1,"ImageList": [{"FaceID": 1,"Name": "1_1.jpg","Size": 3196,"Data": "/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAUDBAQEAwUEBAQFBQUGBwwIBwcHBw8LCwkMEQ8SEhEPERETF"}]}]}`

	return str
}

func (this *Sdk) dealAuthorized2(buf []byte) string {
	str := `{
		"Num": 1,
		"PersonInfoList": [
		  {
			"PersonID": 22,
			"LastChange": 1602329484,
			"PersonCode": "5hh",
			"PersonName": "我的陌生人",
			"Remarks": "陌生人的尝试哈哈哈哈哈嘎",
			"TimeTemplateNum": 0,
			"IdentificationNum": 0,
			"ImageNum": 1,
			"ImageList": [
			  {
				"FaceID": 1,
				"Name": "1_1.jpg",
				"Size": 3196,
				"Data": "/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAUDBAQEAwUEBAQFBQUGBwwIBwcHBw8LCwkMEQ8SEhEPERETFhwXExQaFRERGCEYGh0dHx8fExciJCIeJBweHx7/2wBDAQUFBQcGBw4ICA4eFBEUHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh7/wAARCABAADIDAREAAhEBAxEB/8QAGQAAAgMBAAAAAAAAAAAAAAAABgcEBQgD/8QANhAAAQMCBQIFAwIDCQAAAAAAAQIDBAURAAYSITETQQciUWFxFDKBCKEVI/BicpGSorHB4fH/xAAbAQACAwEBAQAAAAAAAAAAAAADBAECBQAGB//EADIRAAEDAwIEAgkEAwAAAAAAAAEAAgMEESESMQUTQWFRcQYUIoGRobHB8CMy0eFSkvH/2gAMAwEAAhEDEQA/AErBhJKQQN7emH0wrBqEQkkp3HJtjlynU+kvS30R47SnHnVBDaQAdSibAD84g4C5dM55fquXW5MeZDdjym0FQSpFwq3dJGxHuMc1wcuBuoGXqXKrVAdrEdkNRGFFKy+4kLJAGo+g+L/98TpNiovmyr6hBIsFJAN8SpVWqM5qPlTz6Y5Si+hp+sjIkJZQ3dsJsgHzEE3VvbYgjsNwccVCMcqZIr+YVpRTKbIebKvM6UaW0/KjtijpGt3VSQE+fDvwvgZTZFQmpTUKzYBOhN0ME7eXa/yojjgeqkkxdgbIT5LBT8/5RXXaYYryUSgtxvYJ0li4spaST25tb1BvhKXngh0TtiMeITdK6G1pB0Oe/RZf8QsuZi8O6yqElLzURT31LRa+xSrKR1EdrFJKSk7EEg2NiNZhEjboeFHXHpL+XYs1yoxVz5K1FMeKmyG0DbSUnzAi19+L7XFjiWYOkKLm6piy0CRuPxgllK0v+m6FDfySp52nRZK0y3ApSmgpabBFhuN+b4SnJ1Ibk2pRlLQgRZKWGkmykhF1qPZNyfL35HxgbNPVXjdE0HW0k9M4+mfiqTOFacyxlNl+OsSZi3g11HRY6jdRUU+1uOPxhasn5bdQC3OC8Oj4tXFrhpYBew7WAF/ugKg57zA7LZamym3W33UoJcYB06ja9kaSbel8ZkVbLqAJ+S9lX+ivDxE50TCCATgnp53HyRlmPLsPNlBfodcXEW0u6oLyWFNOtK7nQsk23Tf1v22tswyOabkr5zUxsaLxNItvkOHbIA3ysoZqyLUcpZvdiz0WQ3dTZH2rSbgFJ7p/85vjSY4OFwlwbhQOij0/2xZSnl+mGUlml5ggl09Vx9lYTq2TrSUXH+T9sY/FpOUy4wXWaPf/AAj07NbwDsMp6MLSpao9y5rslKCSNCUmxUo+p/ewwjFO8SaW5JwOwBsSe5z547qHxNLdRx9yeg8vkoWfKLErVIUh4uJcZ/moUkBRuAe3e47fGHKiESssU1wPiUlBUhzNnYN0HSMjCFKh1ehp/jUZCw4tlZ0qUjsUHYHb+u2EjR6CHsyF6mP0m9ZZJS1X6TiLAjoe/VFdKpOl2HIizp6YrOpLkeagLVexGylDUnk7gkHtthxjNiCceK8zV1uJI5WNLnWs5pt8hg+RAI65VP4qZRh5qozkF9goeaBVFl3F2lEfuk7XH/IBxPrckUhAZjxuAFkNwsxysiZoYkusGlzFltZQVIbKkmxIuCOR74fFdT/5j4hXuiHwFnLRm6VD12EmIV+xKFAC/v51f44xPShpNM09/sVqcLI5jvJaToMaV9IkyXAwyT9vC3B2BPYc8b4S4TRzmEB5s0/EjuegQuIVMQkOnJ+QP3K4VOdPjpVMLao8qMySuIp0Fp1scqQfUbbn4tj2MMERAjGWnY9QfArzEs8oJecEbi+COpCraDm1rMVJrdNW8zHqMEBxsoV00OMOedlY3228iv7SVYDGx0NSGFt7fn9pmdwlpi8OtfqiqnTm51JamJUkhxAUoavtPcH4wGaIxPLT0V4JWyxh4KCK14g02PIe68VC4EcHXMcWAkAdwD2/O/tjOkDJnW0AnvlIN4uHzcuNhI8UASPH3KLb7jbdKq7yEqIS4mEgBYB5F1XsffBBwyW37h/qFs6Sl14HyGY/ibRVK+xxbjSibWF21Ef6gnGhWRsfEQ8fnj7kQPc25aVrRx2PNpag7F+v07lrY3UDzuQL9+cJ0FTzLOY+x2J/54pOdjXNOpt+yFc006c/TRCpNMqTCOoFLS86FJCd7lI1E82NuMehpJ2NfqleCew/oLFqoHuj0xMcB5/TJSZztGrWVK7TMywg+24jqRXGHQpCHWbhfSVfgK1qt6FN8Vrg18upjt7beI/AtLhjtdNyZGbX36g/hRvI8RstUvw/k1BuclxFUAESCh1JlBe6XEKQPt02sVHY7EbHClc71oNsPa2KWjoJY2yQt/adj9Uj63UKvmpusVBbrTbVMY+qRTg4dDaN7f3lbElR/FthgDGR09gBkrVo6CKBp5Y/ldqTByrKpUSTJrzzL7zCHHWxpshRSCR+DirppASAE8I2W3VJlWbIouZac/Lacjqjy2nF9RJSQgKF+fa+L1jC+B7RvYpfcLWjWZqRQnEP1WqwIDS07mRISi/xc7/jHieECVtSNIuDuqkXCH6z4+5GhlaKb/Eay6DYmLH0pUfdSyNvgHHsmwOKgMJ3St8SvF2ZnKMmkuZZYpsJLgfS4t4uu6gCBY7AXCjtY4PHFoN7qzW2S9fhKTTKjmFqLrixClMl0KHkJIAFuSTccYuXgODepRA0kXUanM5ip9ZU61R3JDFbpak9H6pKHC2RqBFuDZPB5BI74FI+N+L5BWhT0M7oHVAHsjfPilhKalxZTsV1lIcZWW1C6uQSDh0OaRe6zyLJz10lK1POyk9KXZPTCdfbmx27Af0cAfGdRKU5uhoxtuiai5KylWID1ThOPwFBSh1CsOpKgLkkEcD2O9tuxwEFzMIzJBICQMfVCubKa9lysuUl6XGkrSlKtTQNgCLgEHg2sfzgzHahdXCqjLvcJv7+mLLlEpdfhwct5oodV6iVSwHoDqUFQ6qQu17bg7gA+2AyRkva5vvRWOAaQVBezj9HR6EqkuSjUIbWl8SEAtKPoT9x47Y7kanO17FPwV5hp3RNJ9rpi3v+1kKTKrNlTHpTrcfqPOKcXZHckk4YEbALLNJuV//Z"			  
			  }
			]
		  }
		]
	  }`
	return str
}

//dealAuthorized 处理授权
func (this *Sdk) dealAuthorized(buf []byte) string {
	// 获取
	// 	str := `
	// {
	// 	"Num": 1,
	// 	"PersonInfoList": [
	// 		{
	// 			"PersonID": %s,
	// 			"LastChange": %d,
	// 			"PersonCode": "%s",
	// 			"PersonName": "%s",
	// 			"Remarks": "%s",
	// 			"TimeTemplateNum": 0,
	// 			"IdentificationNum": 0,
	// 			"ImageNum": %d,
	// 			"ImageList": %s
	// 		}
	// 	]
	// }
	// `

	// fmt.Println("打印请求参数：", string(buf))
	persionID := gjson.GetBytes(buf, "persionID").String()
	fmt.Println("打印请persionID：", persionID)

	persionName := gjson.GetBytes(buf, "persionName").String()
	remark := gjson.GetBytes(buf, "remark").String()
	imageList := gjson.GetBytes(buf, "imageList").Array()
	timestamp := time.Now().Unix()

	var imageArr []interface{}
	var i = 0
	for _, val := range imageList {
		i++
		temp := map[string]interface{}{
			"FaceID": i,
			"Name":   val.Get("Name").String(),
			"Size":   val.Get("Size").Int(),
			"Data":   val.Get("Data").String(),
		}
		imageArr = append(imageArr, temp)
	}

	type Info struct {
		PersonID          string
		LastChange        int64
		PersonCode        string
		PersonName        string
		Remarks           string
		TimeTemplateNum   int
		IdentificationNum int
		ImageNum          int
		ImageList         []interface{}
	}
	var info = Info{
		PersonID:          persionID,
		LastChange:        timestamp,
		PersonCode:        persionID,
		PersonName:        persionName,
		Remarks:           remark,
		TimeTemplateNum:   0,
		IdentificationNum: 0,
		ImageNum:          i,
		ImageList:         imageArr,
	}
	var authorized = struct {
		Num            int
		PersonInfoList []Info
	}{
		Num:            1,
		PersonInfoList: []Info{info},
	}
	p, _ := json.Marshal(authorized)
	content := string(p)
	fmt.Println("组合json结果：", content)
	// content := fmt.Sprintf(str, persionID, timestamp, persionID, persionName, remark, i, string(p))

	return content
}
