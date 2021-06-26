package main

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type ImageInfo struct {
	ImageID   string `json:"imageID"`   // 图像标识
	EventSort int64  `json:"eventSort"` // 事件分类：自动分析事件类型，设备采集必选
	// DeviceID    string `json:"deviceID"`// 设备编码
	// StoragePath string `json:"faceID"`// 存储路径：图像文件的存储路径，采用URI命名规则
	// Type       string `json:"type"`       // 不清楚字段含义
	FileFormat string `json:"fileFormat"` // 图像文件格式
	ShotTime   string `json:"shotTime"`   // 拍摄时间
	Width      int64  `json:"width"`      // 宽度
	Height     int64  `json:"height"`     // 高度
	Data       string `json:"data"`       // 图片数据，使用BASE64加密
}

type FaceObject struct {
	FaceID            string `json:"faceID"`            // 人脸标识
	InfoKind          int64  `json:"infoKind"`          // 信息分类：人工采集/自动采集
	SourceID          string `json:"sourceID"`          // 来源标识：来源图像信息标识
	DeviceID          string `json:"deviceID"`          // 设备编码，自动采集必选
	ShotTime          string `json:"shotTime"`          // 拍摄时间
	LeftTopX          int64  `json:"leftTopX"`          // 左上角X坐标
	LeftTopY          int64  `json:"leftTopY"`          // 左上角Y坐标
	RightBtmX         int64  `json:"rightBtmX"`         // 右下角X坐标
	RightBtmY         int64  `json:"rightBtmY"`         // 右下角Y坐标
	FaceAppearTime    string `json:"faceAppearTime"`    // 人脸出现时间
	FaceDisAppearTime string `json:"faceDisAppearTime"` // 人脸消失时间
	Similaritydegree  int64  `json:"similaritydegree"`  // 相似度：人脸相似度[0,1]
	// LocationMarkTime  string `json:"locationMarkTime"` // 位置标记时间
	// GenderCode           string // 性别代码
	// AgeUpLimit           int    // 年龄上限
	// AgeLowerLimit        int    // 年龄下限
	// GlassStyle           string // 眼镜款式
	// Emotion              string // 不清楚字段含义
	// IsDriver             int    // 是否驾驶员：0-否；1-是；2-不确定
	// IsForeigner          int    // 是否涉外人员：0-否；1-是；2-不确定
	// IsSuspectedTerrorist int    // 是否涉恐人员：0-否；1-是；2-不确定
	// IsCriminalInvolved   int    // 是否涉案人员：0-否；1-是；2-不确定
	// IsDetainees          int    // 是否在押人员：：0-否；1-是；2-不确定，人工采集必填
	// IsVictim             int    // 是否被害人：0-否；1-是；2-不确定
	// IsSuspiciousPerson   int    // 是否可疑人：0-否；1-是；2-不确定

	// SubImageList         interface{}
	// RelatedType          string
	// RelatedList          interface{}
	ImageList []ImageInfo `json:"data"`
}

func main() {
	// 读取 faceList的列表
	id := gjson.Get(str, "FaceListObject.FaceObject.0.DeviceID").String()
	fmt.Println(id)
	return
	faceArr := gjson.Get(str, "FaceListObject.FaceObject").Array()
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
	// 处理发送或者打印
	buf, _ := json.Marshal(objectList)
	fmt.Println(string(buf))
}

var str = `{
    "FaceListObject": {
        "FaceObject": [
            {
                "FaceID": "111110220200710143217001770600178",
                "InfoKind": 1,
                "SourceID": "11111022020071014321700177",
                "DeviceID": "11111",
                "ShotTime": "20200710143217",
                "LeftTopX": 512,
                "LeftTopY": 369,
                "RightBtmX": 749,
                "RightBtmY": 707,
                "LocationMarkTime": "20200710143217",
                "FaceAppearTime": "20200710143217",
                "FaceDisAppearTime": "20200710143217",
                "GenderCode": "1",
                "AgeUpLimit": 28,
                "AgeLowerLimit": 28,
                "GlassStyle": "99",
                "Emotion": "1",
                "IsDriver": 2,
                "IsForeigner": 2,
                "IsSuspectedTerrorist": 2,
                "IsCriminalInvolved": 2,
                "IsDetainees": 2,
                "IsVictim": 2,
                "IsSuspiciousPerson": 2,
                "Similaritydegree": 0,
                "SubImageList": {
                    "SubImageInfoObject": [
                        {
                            "ImageID": "11111022020071014321700177",
                            "EventSort": 10,
                            "DeviceID": "11111",
                            "StoragePath": "",
                            "Type": "14",
                            "FileFormat": "Jpeg",
                            "ShotTime": "20200710143217",
                            "Width": 1920,
                            "Height": 1264,
                            "Data": "图片数据"
                        },
                        {
                            "ImageID": "11111022020071014321700180",
                            "EventSort": 10,
                            "DeviceID": "11111",
                            "StoragePath": "",
                            "Type": "11",
                            "FileFormat": "Jpeg",
                            "ShotTime": "20200710143217",
                            "Width": 896,
                            "Height": 700,
                            "Data": "图片数据"
                        }
                    ]
                },
                "RelatedType": "01",
                "RelatedList": {
                    "RelatedObject": [
                        {
                            "RelatedType": "01",
                            "RelatedID": "111110220200710143217001770100179"
                        }
                    ]
                }
            }
        ]
    }
}`
