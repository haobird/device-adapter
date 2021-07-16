package guardface

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/haobird/goutils"
	"github.com/haobird/logger"
	"github.com/tidwall/gjson"
)

// 针对设备具体功能处理物模型

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

// 定义处理结构体
type Things struct{}

// 解析上行数据包 ,返回 action / body / error
func (t *Things) _parseReq(buf []byte) (string, string, error) {
	action := "unknown"
	// 分割字符串
	str := string(buf)
	pos := strings.Index(str, "{")
	if pos < 1 {
		return action, str, errors.New("未读取到数据")
	}

	var input map[string]interface{}
	content := str[pos:]
	err := json.Unmarshal([]byte(content), &input)
	if err != nil {
		return action, content, errors.New("解析数据错误")
	}

	if strings.Contains(str, "HeartReportInfo") {
		action = "HeartReportInfo"
	} else if strings.Contains(str, "PersonVerification") && !strings.Contains(str, "Notifications") {
		action = "PersonVerification"
	} else if strings.Contains(str, "DeviceBasicInfo") {
		action = "DeviceBasicInfo"
	} else if strings.Contains(str, "Response") {
		action = "reply"
	}

	return action, content, nil
}

// ProcessDataUp 处理上行消息
func (t *Things) ProcessDataUp(buf []byte) *Package {
	var packet *Package
	var message_type = Unknown
	var requestID = ""
	var topic = ""

	// 解析读取内容
	action, content, err := t._parseReq(buf)
	if err != nil {
		logger.Error("[things] ProcessDataUp error: ", err)
		return packet
	}

	// 提取 requestid
	requestID = between(string(buf), "RequestID: ", "\r\n")

	// 进行数据的格式化转换
	var result map[string]interface{}
	switch action {
	case "HeartReportInfo":
		message_type = Heart
	case "PersonVerification":
		message_type = Publish
		topic = "personVerification"
		result = t.business_handler_personVerification(content)
	case "DeviceBasicInfo":
		message_type = Publish
		topic = "deviceInfoFace"
		result = t.business_handler_deviceInfo(content)
	case "reply":
		message_type = PubAck
	}

	packet = &Package{
		MessageType: message_type,
		ClientID:    "",
		RequestID:   requestID,
		Topic:       topic,
	}

	var p []byte
	if result != nil {
		p, _ = json.Marshal(result)
	} else {
		p = []byte(content)
	}
	packet.Data = p

	return packet
}

// ProcessDataDown 处理下行消息 (接收 Package消息，返回 发送给设备的包)
func (t *Things) ProcessDataDown(packet *Package) []byte {
	topic := packet.Topic
	// 设备白名单
	if topic == "personAuthorized" {
		body := t.business_command_personAuthorized(packet.Data)
		// body := string(pack.Data)
		header := map[string]string{"RequestID": packet.RequestID}
		content := t.PackReq("POST", "/LAPI/V1.0/PeopleLibraries/3/People", body, header)
		// fmt.Println(content)
		return []byte(content)
	}

	// 删除白名单
	if topic == "personAuthorizedCancel" {
		header := map[string]string{"RequestID": packet.RequestID}
		persionID := gjson.GetBytes(packet.Data, "persionID").String()
		timestamp := time.Now().Unix()
		url := fmt.Sprintf("/LAPI/V1.0/PeopleLibraries/3/People/%s?LastChange=%d", persionID, timestamp)
		content := t.PackReq("DELETE", url, "", header)
		// fmt.Println(content)
		return []byte(content)
	}

	if topic == "personSearch" {
		body := t.business_command_personSearch(packet.Data)
		header := map[string]string{"RequestID": packet.RequestID}
		content := t.PackReq("POST", "/LAPI/V1.0/PeopleLibraries/3/People/Info", body, header)
		return []byte(content)
	}

	// 远程开门
	if topic == "deviceOpenFace" {
		header := map[string]string{"RequestID": packet.RequestID}
		content := t.PackReq("PUT", "/LAPI/V1.0/PACS/Controller/RemoteOpened", "", header)
		return []byte(content)
	}

	// 设备信息
	if topic == "deviceInfoFace" {
		header := map[string]string{"RequestID": packet.RequestID}
		content := t.PackReq("GET", "/LAPI/V1.0/System/DeviceBasicInfo", "", header)
		return []byte(content)
	}

	return nil
}

func (this *Things) PackResp(body string, headers map[string]string) string {
	respHeader := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain;charset=ISO-8859-1\r\n"
	for key, val := range headers {
		respHeader = respHeader + key + ": " + val + "\r\n"
	}
	//拼装http返回的消息
	resp := respHeader + fmt.Sprintf("Content-Length: %d\r\n", len(body)) + "\r\n" + body
	return resp
}

// 封装 请求
func (t *Things) PackReq(method, url, body string, headers map[string]string) string {
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

//ParseConnect 解析连接信息
func (t *Things) ParseConnect(buf []byte) (string, error) {
	action, content, err := t._parseReq(buf)
	if err != nil {
		return "", err
	}
	if action != "HeartReportInfo" {
		return "", errors.New("不正确的数据包:" + action)
	}
	var input map[string]interface{}
	json.Unmarshal([]byte(content), &input)
	if deviceCode, ok := input["DeviceCode"]; ok {
		clientID := goutils.String(deviceCode)
		return clientID, nil
	}
	return "", errors.New("设备编号不存在")
}

// 设备状态封包
func (t *Things) Packet_deviceStatus(flag int) *Package {
	result := t.business_handler_deviceStatus(flag)
	buf, _ := json.Marshal(result)
	packet := &Package{
		MessageType: Publish,
		RequestID:   "",
		ClientID:    "",
		Topic:       "deviceStatusFace",
		Data:        buf,
	}
	return packet
}

// 心跳封包
func (t *Things) Packet_heart() *Package {
	str := `{"ResponseURL": "/LAPI/V1.0/PACS/Controller/HeartReportInfo", "Code": 0,"Data": {"Time": "%s" }}`
	str = fmt.Sprintf(str, goutils.GetNormalTimeString(time.Now()))
	resp := t.PackResp(str, nil)

	packet := &Package{
		MessageType: PubAck,
		RequestID:   "",
		ClientID:    "",
		Topic:       "pong",
		Data:        []byte(resp),
	}
	return packet
}

// 获取设备信息的请求
func (t *Things) Packet_DeviceInfoReq() *Package {
	packet := &Package{
		MessageType: Publish,
		RequestID:   goutils.String(time.Now().Unix()),
		ClientID:    "",
		Topic:       "deviceInfoFace",
	}
	data := t.ProcessDataDown(packet)
	packet.Data = data
	return packet
}

// 通行记录上报
func (t *Things) business_handler_personVerification(body string) map[string]interface{} {
	var info PersonVerifyInfo
	var respInfo map[string]interface{}

	err := json.Unmarshal([]byte(body), &info)
	if err == nil {
		// libMatInfoNum := info.LibMatInfoNum
		faceInfoNum := info.FaceInfoNum
		panoImage := Image{}
		faceImage := Image{}

		if faceInfoNum == 1 || faceInfoNum == 2 {
			panoImage = info.FaceInfoList[0].PanoImage
			faceImage = info.FaceInfoList[0].FaceImage
		}

		respInfo = map[string]interface{}{
			"persionID":   goutils.String(info.LibMatInfoList[0].MatchPersonID),
			"personCode":  info.LibMatInfoList[0].MatchPersonInfo.PersonCode,
			"persionName": info.LibMatInfoList[0].MatchPersonInfo.PersonName,
			"identityNo":  info.LibMatInfoList[0].MatchPersonInfo.IdentityNo,
			"openType":    "1",
			"panoImage":   panoImage,
			"faceImage":   faceImage,
			"timestamp":   info.Timestamp,
		}
	}
	return respInfo
}

// 设备信息上报
func (t *Things) business_handler_deviceInfo(body string) map[string]interface{} {
	respInfo := map[string]interface{}{
		"mac":        gjson.Get(body, "Response.Data.MAC").String(),
		"ip":         gjson.Get(body, "Response.Data.Address").String(),
		"deviceCode": gjson.Get(body, "Response.Data.SerialNumber").String(),
		"name":       gjson.Get(body, "Response.Data.DeviceModel").String(),
	}
	return respInfo
}

// 设备上线、离线 上报
func (t *Things) business_handler_deviceStatus(flag int) map[string]interface{} {
	respInfo := map[string]interface{}{
		"online":    flag,
		"timestamp": time.Now().Unix(),
	}
	return respInfo
}

// 增加白名单
func (t *Things) business_command_personAuthorized(buf []byte) string {
	/*
	   数据结构
	   {
	   	"Num": 1,
	   	"PersonInfoList": [
	   		{
	   			"PersonID": %s,
	   			"LastChange": %d,
	   			"PersonCode": "%s",
	   			"PersonName": "%s",
	   			"Remarks": "%s",
	   			"TimeTemplateNum": 0,
	   			"IdentificationNum": 0,
	   			"ImageNum": %d,
	   			"ImageList": %s
	   		}
	   	]
	   }
	*/
	// fmt.Println("打印请求参数：", string(buf))
	persionID := gjson.GetBytes(buf, "persionID").String()
	fmt.Println("打印请persionID：", persionID)

	persionName := gjson.GetBytes(buf, "persionName").String()
	persionCode := gjson.GetBytes(buf, "persionCode").String()
	// identityNo := gjson.GetBytes(buf, "IdentityNo").String()
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
		PersonID   string
		LastChange int64
		PersonCode string
		PersonName string
		// IdentityNo        string
		Remarks           string
		TimeTemplateNum   int
		IdentificationNum int
		ImageNum          int
		ImageList         []interface{}
	}
	var info = Info{
		PersonID:   persionID,
		LastChange: timestamp,
		PersonCode: persionCode,
		PersonName: persionName,
		// IdentityNo:        identityNo,
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

// 删除白名单
func (t *Things) business_command_personAuthorizedCancel(data []byte) string {
	return ""
}

// 查询人员信息
func (t *Things) business_command_personSearch(data []byte) string {
	str := `{ "Num": 0, "QueryInfos": [ {"QryType": 27, "QryCondition": 0, "QryData": "1001" },{"QryType": 55, "QryCondition": 0, "QryData": "Uniview" }],"Limit": 10, "Offset": 0 }`
	return str
}

// 远程开门
func (t *Things) business_command_deviceOpen(data []byte) string {
	return ""
}

// 获取设备信息
func (t *Things) business_command_deviceInfo(data []byte) string {
	return ""
}
