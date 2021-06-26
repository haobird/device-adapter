package guardcamera

import (
	"encoding/json"
	"time"

	"github.com/tidwall/gjson"
)

//RegisterHandler 注册
func RegisterHandler(clientID string) {
	_, err := cache.Value(clientID)
	if err == nil {
		return
	}

	// 不存在，则添加缓存
	cache.Add(clientID, 2*time.Minute, clientID)

	// 上报上线的信息
	result := map[string]interface{}{
		// "parkId":    "10000",
		"online":    1,
		"timestamp": time.Now().Unix(),
	}
	p, _ := json.Marshal(result)

	packet := Package{
		MessageType: Connect,
		RequestID:   "",
		ClientID:    clientID,
		Topic:       "deviceStatusCamera",
		Data:        p,
	}
	ProcessPublsih(packet)
}

// HeartHandler 心跳
func HeartHandler(clientID string) {
	// 保持心跳的信息
	cache.Value(clientID)
}

// DisconnectHandler 断开
func DisconnectHandler(clientID string) {
	// cache.Delete(clientID)
	// 上报离线的信息
	result := map[string]interface{}{
		"online":    0,
		"timestamp": time.Now().Unix(),
	}
	p, _ := json.Marshal(result)

	packet := Package{
		MessageType: Disconnect,
		RequestID:   "",
		ClientID:    clientID,
		Topic:       "deviceStatusCamera",
		Data:        p,
	}
	ProcessPublsih(packet)
}

// 处理 人脸抓拍
func FacesHandler(clientID string, buf []byte) {
	// 读取 faceList的列表
	faceArr := gjson.GetBytes(buf, "FaceListObject.FaceObject").Array()
	var objectList = []FaceObject{}
	for _, item := range faceArr {
		object := FaceObject{
			FaceID:    item.Get("FaceID").String(),
			InfoKind:  item.Get("InfoKind").Int(),
			SourceID:  item.Get("SourceID").String(),
			DeviceID:  item.Get("DeviceID").String(),
			ShotTime:  item.Get("LeftTopX").String(),
			LeftTopX:  item.Get("LeftTopY").Int(),
			LeftTopY:  item.Get("LeftTopY").Int(),
			RightBtmX: item.Get("RightBtmX").Int(),
			RightBtmY: item.Get("RightBtmY").Int(),
			// LocationMarkTime:  item.Get("LocationMarkTime").String(),
			FaceAppearTime:    item.Get("FaceAppearTime").String(),
			FaceDisAppearTime: item.Get("FaceDisAppearTime").String(),
			Similaritydegree:  item.Get("Similaritydegree").Int(),
		}
		imageArr := item.Get("SubImageList.SubImageInfoObject").Array()
		var imageList = []ImageInfo{}
		for _, image := range imageArr {
			imageInfo := ImageInfo{
				ImageID:   image.Get("ImageID").String(),
				EventSort: image.Get("EventSort").Int(),
				// Type:       image.Get("Type").String(),
				FileFormat: image.Get("FileFormat").String(),
				ShotTime:   image.Get("ShotTime").String(),
				Width:      image.Get("Width").Int(),
				Height:     image.Get("Height").Int(),
				Data:       image.Get("Data").String(),
			}
			imageList = append(imageList, imageInfo)
		}
		object.ImageList = imageList
		objectList = append(objectList, object)
	}
	// 单个发送，循环处理
	for _, value := range objectList {
		// 处理发送或者打印
		p, _ := json.Marshal(value)

		packet := Package{
			MessageType: Publish,
			RequestID:   "",
			ClientID:    clientID,
			Topic:       "faceShot",
			Data:        p,
		}
		ProcessPublsih(packet)
	}

}
