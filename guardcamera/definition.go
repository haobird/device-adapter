package guardcamera

const (
	Unknown    = "unknown"
	Connect    = "connect"
	Heart      = "heart"
	Publish    = "publish"
	PubAck     = "puback"
	Command    = "command"
	Disconnect = "disconnect"
)

// 定义消息包
type Package struct {
	MessageType string // 消息类型
	RequestID   string // 请求id(时间戳)
	ClientID    string // 设备序列号
	Action      string // 业务类型
	Topic       string // 主题
	Data        []byte // 载荷
}

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
	ImageList []ImageInfo `json:"imageList"`
}
