package carguard

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/haobird/goutils"
	"github.com/tidwall/gjson"
)

// 处理设备模型
type Things struct{}

// ParsePublishData 解析请求数据
func (t *Things) ParsePublishData(action string, buf []byte) *Package {
	// 根据对应的功能，做相应的内容处理
	version := gjson.GetBytes(buf, "version").String()
	parkId := gjson.GetBytes(buf, "parkId").String()
	deviceId := gjson.GetBytes(buf, "deviceId").String()
	params := gjson.GetBytes(buf, "params").String()

	message_type := Publish
	topic := action

	var result = map[string]interface{}{}
	switch action {
	case "basicinfo":
		topic = "deviceInfoCar"
		result = t.basicinfoHandler(params)
	case "parkalarm":
		message_type = Unknown
		result = t.parkalarmHandler(params)
	case "keepalive":
		message_type = Heart
		result = map[string]interface{}{}
	case "capture":
		topic = "plateVerification"
		result = t.captureHandler(params)
	case "connect":
		topic = "deviceStatusCar"
		result = t.registerHandler(params)
	}

	result["parkId"] = parkId
	fmt.Println("version:", version)

	data, _ := json.Marshal(result)

	return &Package{
		MessageType: message_type,
		RequestID:   "",
		ClientID:    deviceId,
		Topic:       topic,
		Data:        data,
	}
}

// ParsePubackData 解析响应的数据
func (t *Things) ParsePubackData(action string, buf []byte) *Package {
	// 根据对应的功能，做相应的内容处理
	requestId := gjson.GetBytes(buf, "requestId").String()
	parkId := gjson.GetBytes(buf, "parkId").String()
	deviceId := gjson.GetBytes(buf, "deviceId").String()
	params := gjson.GetBytes(buf, "data").String()
	code := gjson.GetBytes(buf, "code").String()
	message := gjson.GetBytes(buf, "message").String()

	message_type := PubAck
	topic := action

	var result = map[string]interface{}{}
	var data = map[string]interface{}{}
	switch action {
	case "gatecontrol":
		topic = "gatecontrol"
		data = t.basicinfoHandler(params)
	}

	data["parkId"] = parkId
	result["data"] = data
	result["code"] = code
	result["message"] = message

	p, _ := json.Marshal(result)

	return &Package{
		MessageType: message_type,
		RequestID:   requestId,
		ClientID:    deviceId,
		Topic:       topic,
		Data:        p,
	}
}

func (t *Things) ParseCommanData(packet Package) *Package {
	topic := packet.Topic
	buf := packet.Data

	parkId := gjson.GetBytes(buf, "parkId").String()
	if parkId == "" {
		parkId = "10000"
	}

	var params = map[string]interface{}{}
	params["deviceId"] = packet.ClientID
	params["requestId"] = packet.RequestID
	params["parkId"] = parkId
	params["raw"] = buf

	var result = map[string]interface{}{}
	switch topic {
	case "deviceOpenCar":
		result = t.commandOpenHandler(params)
	case "plateAuthorized":
		topic = "plateAuthorized"
		result = t.commandAuthorizedHandler(params)
	case "plateAuthorizedCancel":
		topic = "plateAuthorizedCancel"
		result = t.commandAuthorizedCancelHandler(params)

	}

	p, _ := json.Marshal(result)

	return &Package{
		MessageType: Command,
		RequestID:   packet.RequestID,
		ClientID:    packet.ClientID,
		Topic:       topic,
		Data:        p,
	}
}

// Register 注册包
func (t *Things) registerHandler(str string) map[string]interface{} {
	// 解析json数据
	result := map[string]interface{}{
		// "parkId":    "10000",
		"online":    1,
		"timestamp": time.Now().Unix(),
	}
	return result
}

func (t *Things) basicinfoHandler(str string) map[string]interface{} {
	// 解析json数据
	result := map[string]interface{}{
		"name":       gjson.Get(str, "deviceName").String(),
		"ip":         gjson.Get(str, "ipAddress").String(),
		"deviceCode": gjson.Get(str, "serialNum").String(),
		"mac":        gjson.Get(str, "MAC").String(),
	}
	return result
}

func (t *Things) parkalarmHandler(str string) map[string]interface{} {
	result := map[string]interface{}{}
	return result
}

func (t *Things) captureHandler(str string) map[string]interface{} {
	picTimeStr := gjson.Get(str, "picTime").String()
	result := map[string]interface{}{
		"recordId":     gjson.Get(str, "recordId").String(),
		"picTime":      picTimeStr,
		"plateNo":      gjson.Get(str, "plateNo").String(),
		"confidence":   gjson.Get(str, "confidence").String(),
		"vehicleType":  gjson.Get(str, "vehicleType").String(),
		"vehicleColor": gjson.Get(str, "vehicleColor").String(),
		"plateType":    gjson.Get(str, "plateType").String(),
		"plateColor":   gjson.Get(str, "plateColor").String(),
		"picInfo":      gjson.Get(str, "picInfo").String(),
	}
	picTime, _ := goutils.GetTimeByString(picTimeStr)
	result["timestamp"] = goutils.GetTimeUnix(picTime)
	return result
}

// 手动抓拍结果（出入口）
func (t *Things) manualcaptureResultHandler(str string) map[string]interface{} {
	result := map[string]interface{}{}
	return result
}

// 开闸放行结果
func (t *Things) gatecontrolResultHandler(str string) map[string]interface{} {
	result := map[string]interface{}{}
	return result
}

// [指令] 定义为

// [指令] 远程开门
func (t *Things) commandOpenHandler(params map[string]interface{}) map[string]interface{} {
	var payload = map[string]interface{}{
		"version":   "1.0",
		"requestId": params["requestId"],
		"parkId":    params["parkId"],
		"deviceId":  params["deviceId"],
		"type":      5,
		"params":    map[string]interface{}{},
	}
	var result = map[string]interface{}{
		"url":     "http://ip:port/LAPI/V1.0/ParkingLots/Entrances/Lanes/0/GateControl",
		"method":  "POST",
		"body":    `{ "Command": 0}`,
		"payload": payload,
	}
	return result
}

// [指令] 下发白名单
func (t *Things) commandAuthorizedHandler(params map[string]interface{}) map[string]interface{} {

	raw := params["raw"].([]byte)
	var body = map[string]interface{}{
		"Num": 1,
		"AllowListInfo": []map[string]string{
			{
				"AllowID":   gjson.GetBytes(raw, "id").String(),
				"PlateNo":   gjson.GetBytes(raw, "plateNo").String(),
				"OwnerName": gjson.GetBytes(raw, "ownerName").String(),
				"PhoneNo":   gjson.GetBytes(raw, "phoneNo").String(),
				"BeginTime": gjson.GetBytes(raw, "beginTime").String(),
				"EndTime":   gjson.GetBytes(raw, "endTime").String(),
				"Remarks":   gjson.GetBytes(raw, "remark").String(),
			},
		},
	}

	var payload = map[string]interface{}{
		"requestId": params["requestId"],
		"version":   "1.0",
		"parkId":    params["parkId"],
		"deviceId":  params["deviceId"],
		"type":      4,
		"params": map[string]interface{}{
			"listType": 1,
			"mode":     0,
			"num":      1,
			"listInfo": []map[string]string{
				{
					"plateNo":   gjson.GetBytes(raw, "plateNo").String(),
					"startTime": gjson.GetBytes(raw, "beginTime").String(),
					"endTime":   gjson.GetBytes(raw, "endTime").String(),
				},
			},
		},
	}

	var result = map[string]interface{}{
		"url":     "/LAPI/V1.0/ParkingLots/Vehicles/AllowList",
		"method":  "POST",
		"body":    body,
		"payload": payload,
	}
	return result
}

// [指令] 删除白名单
func (t *Things) commandAuthorizedCancelHandler(params map[string]interface{}) map[string]interface{} {

	raw := params["raw"].([]byte)
	fmt.Println("rwa::", string(raw))
	id := gjson.GetBytes(raw, "id").String()
	plateNo := gjson.GetBytes(raw, "plateNo").String()

	var payload = map[string]interface{}{
		"requestId": params["requestId"],
		"version":   "1.0",
		"parkId":    params["parkId"],
		"deviceId":  params["deviceId"],
		"type":      4,
		"params": map[string]interface{}{
			"listType": 1,
			"mode":     1,
			"num":      1,
			"listInfo": []map[string]string{
				{
					"plateNo": plateNo,
				},
			},
		},
	}

	var result = map[string]interface{}{
		"url":     "/LAPI/V1.0/ParkingLots/Vehicles/AllowList/" + id,
		"method":  "POST",
		"body":    nil,
		"payload": payload,
	}
	return result
}
