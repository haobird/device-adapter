# 门禁适配器说明

## 功能

- 设备状态上报
- 上报设备信息
- 上报通行记录
- 下发白名单
- 删除白名单
- 远程开门

## 数据上报

涉及到 接口调用的统一使用格式：

 HTTP接口 (Post): http://ip:port/接口路径

 Content-Type: application/json;charset=UTF-8


## 字段介绍

* 涉及图片的都是 jpg/ jpeg 格式的
* base64编码图片 数据头不需要增加 base64，以/9j/开头


### 上报 方式

1. http 请求业务开放的接口

> 使用同一个接口，上报内容，业务自行通过 topic 判断上报类型

2. kafka topic publish 

> kafka 上报内容只有 payload ，业务订阅不同的topic

### 上报设备状态 消息结构 (注意 kafka: 只有payload)

```
{
    "topic" : "deviceStatusFace",
    "key" : "210235C3R0320B000985",   // deviceCode 设备序列号
    "payload" : {
        "online" : 0, // 0: 离线 ， 1: 在线
        "timestamp" : 1564735558
    }
    
}
```

### 上报设备信息 消息结构 (注意 kafka: 只有payload)

```
{
    "topic" : "deviceInfoFace",
    "key" : "210235C3R0320B000985",  // deviceCode 设备序列号
    "payload" : {
        "mac" : "xxxxxx",
        "ip" : "192.168.1.13",
        "deviceCode" : "210235C3R0320B000985",
        "name" : "设备1"
    }
    
}
```

### 通行识别记录 消息结构

```
{
    "topic" : "personVerification",
    "key" : "210235C3R0320B000985",  // deviceCode 设备序列号
    "payload" : {
        "persionID" : "111",  // 住户id（注意设备支持 20位长度的字符）
        "persionName" : "住户1", // 住户姓名
        "openType": "1",  // 开门类型 1: 人脸开门
        "panoImage" : {
            "Name": "1564707615_1_86.jpg", 
            "Size": 101780, 
            "Data": "…"  // base64
        },
        "faceImage" : {
            "Name": "1564707615_1_86.jpg", 
            "Size": 101780, 
            "Data": "…"  // base64
        },
        "timestamp" : 1564735558
    }
    
}
```

## api接口 业务请求此中间件

接口： http://ip:port/api/faceguard

### 返回字段介绍 

| 字段 | 类型 | 说明 |
| -- | -- | -- | -- | 
| Code | int | 错误码（0为正确，其它为错误） | 
| Msg | string | 响应消息（错误消息） | 
| Data | array | 数据体 | 


### 下发白名单

请求示例：
```
{
    "topic" : "personAuthorized",
    "key" : "210235C3R0320B000985",  // deviceCode 设备序列号
    "payload" : {
        "persionID" : "111",  // 住户id（注意设备支持20位长度的字符串）
        "persionName" : "住户1", // 住户姓名（注意设备支持20位长度的字符串）
        "remark" : "备注", //（注意设备支持20位长度的字符串）
        "imageList" : [
            {
                "Name": "1564707615_1_86.jpg", 
                "Size": 101780, 
                "Data": "…"  // base64
            },
            {
                "Name": "1564707615_1_86.jpg", 
                "Size": 101780, 
                "Data": "…"  // base64
            }
        ]
    }
}
```

响应示例：
```
{
    "Code": 0, # 错误码
    "Msg": "", # 响应消息
    "Data": null
}
```
### 删除白名单

请求示例：
```
{
    "topic" : "personAuthorizedCancel",
    "key" : "210235C3R0320B000985",  // deviceCode 设备序列号
    "payload" : {
        "persionID" : "111"
    }
}
```

响应示例：
```
{
    "Code": 0, # 错误码
    "Msg": "", # 响应消息
    "Data": null
}
```

### 远程开门

请求示例：
```
{
    "topic" : "deviceOpenFace",
    "key" : "210235C3R0320B000985",  // deviceCode 设备序列号
    "payload" : {
        "timestamp" : 1564735558
    }
}
```

响应示例：
```
{
    "Code": 0, # 错误码
    "Msg": "", # 响应消息
    "Data": null
}
```

### 枚举字典


```
消息定义：
type Package struct {
	MessageType string `json:"message_type"` // 消息类型
	RequestID   string `json:"request_id"`   // 请求id(时间戳)
	ClientID    string `json:"client_id"`    // 设备序列号
	Topic       string `json:"topic"`        // 主题
	Data        []byte `json:"data"`         // 载荷
}


消息类型 ： 表示 设备连接状态
| message_type | 说明 |
| --   | -- | 
| unknown | 未知 |
| connect | 连接或注册 |
| heart | 心跳 |
| publish | 发布或请求 |
| puback | 确认或响应 |
| command | 指令 |
| disconnect | 断开连接 |

业务类型
| topic | 关键字 | 说明 | 方向 |
| --   | -- | 
| heart | HeartReportInfo | 心跳 | 上报 |
| PersonVerification | PersonVerification | 人脸识别上报 | 上报 |
| deviceInfoFace | DeviceBasicInfo | 设备信息 | 上报 |
| reply | Response | 响应 | 上报 |
| deviceStatut | HeartReportInfo | 在线状态（第一次心跳） | 上报 |
| personAuthorized | personAuthorized | 白名单下发 | 指令 |
| personAuthorized | personAuthorizedCancel | 白名单取消 | 指令 |
| deviceOpenFace | RemoteOpened | 远程开门 | 指令 |
| deviceInfoFace | DeviceBasicInfo | 设备信息 | 指令 |
| personSearch | personSearch | 人员信息 | 指令 |

```


#### 统一错误码

1-9 为特殊占用

| 错误码 | 说明 |
| --   | -- | 
| 0   | success | 
| 1   | 设备不在线 | 
| 2   | 参数错误 | 

#### 人脸识别错误码定义

| 错误码 | 说明 |
| --   | -- | 
| 1000 | 算法初始化失败 |
| 1001  | 人脸检测失败 | 
| 1002  | 图片未检测到人脸 | 
| 1003  | jpeg照片解码失败 | 
| 1004  | 图片质量分数不满足 | 
| 1005  | 图片缩放失败 | 
| 1006  | 未启用智能 | 
| 1007  | 导入图片过小 | 
| 1008  | 导入图片过大 | 
| 1009  | 导入图片分辨率超过1920*1080 
| 1010  | 导入图片不存在 |
| 1011 | 人脸元素个数已达到上限 |  
| 1012 | 智能棒算法模型不匹配 |  
| 1013 | 人脸导入库成员证件号非法 |  
| 1014 | 人脸导入库成员图片格式错误 |  
| 1015 | 通道布控已达设备能力上限 |  
| 1016 | 其它客户端正在进行操作人脸 库 | 
| 1017 | 人脸库文件正在更新中 | 
| 1018 | Json反序列化失败 | 
| 1019 | Base64解码失败 | 
| 1020 | 人脸照片，编码后的大小和实 际接收到的长度不一致 |








