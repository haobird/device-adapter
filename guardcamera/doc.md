# 视频接口

## 视频播放方式

步骤：
1. 预约推流
2. 播放视频流

示例url：（摄像头数据流）
预约url：http://82.157.107.86:8080/api/gb28181/invite?id=34020000001320000005&channel=34020000001320000005
播放url：http://82.157.107.86:2020/34020000001320000005/34020000001320000005.flv

测试数据流：（可以直接播放，无需预约的）
播放url：http://82.157.107.86:2020/live/test.flv

华为摄像头url：
预约url：http://82.157.107.86:8080/api/gb28181/invite?id=34020000001320000002&channel=34020000001310000003
播放url：http://82.157.107.86:2020/34020000001320000002/34020000001310000003.flv

### 预约说明

```
预约请求参考：
Request URL: http://82.157.107.86:8080/api/gb28181/invite?id=34020000001320000005&channel=34020000001320000005
Request Method: GET
Status Code: 200 OK

返回码为 304 或者 200 即可 进行下一步播放
返回码为 404 则设备不在线
```

### 播放说明

```
随着播放持续时间越长，延迟越高
```

## 数据上报

涉及到 接口调用的统一使用格式：

 HTTP接口 (Post): http://ip:port/接口路径

 Content-Type: application/json;charset=UTF-8

## 字段介绍

* 涉及图片的都是 jpg/ jpeg 格式的
* base64编码图片 数据头不需要增加 base64，以/9j/开头
* 涉及到的枚举类型看文档最后

### 上报设备状态 消息结构 (注意 kafka: 只有payload)

```
{
    "topic" : "deviceStatusCamera",
    "key" : "34020000001320000002",
    "payload" : {
        "online" : 0, // 0: 离线 ， 1: 在线
        "timestamp" : 1564735558
    }
    
}
```

### 人脸抓拍上报 消息结构 (注意 kafka: 只有payload)

```
{
    "topic" : "faceShot",
    "key" : "34020000001320000002",
    "payload" : {
        "faceID": "111110220200710143217001770600178",  // 人脸标识
        "infoKind" : 1, // 信息分类：人工采集0/自动采集1
        "sourceID" : "11111022020071014321700177", // 来源标识：来源图像信息标识
        "shotTime" : "20200710143217", // 拍摄时间
        "leftTopX" : 512,   // 左上角X坐标
        "leftTopY" : 512,   // 左上角Y坐标
        "rightBtmX" : 512,  // 右下角X坐标
        "rightBtmY" : 512,  //  右下角Y坐标
        "faceAppearTime" : "",  // 人脸出现时间
        "faceDisAppearTime" : "", // 人脸消失时间
        "similaritydegree" : 0, // 相似度：人脸相似度[0,1]
        "imageList" : [
            {
                "imageID": "11111022020071014321700177",   // 图像标识
                "eventSort": 10,        // 事件分类：自动分析事件类型，设备采集必选
                "deviceID": "11111",      // 设备编码
                "fileFormat": "Jpeg",      // 图像文件格式
                "shotTime": "20200710143217",      // 拍摄时间
                "width": 1920,      // 宽度
                "height": 1264,      // 高度
                "data": "图片数据"      // 图片数据，使用BASE64加密
            },
            {
                "imageID": "11111022020071014321700180",
                "eventSort": 10,
                "deviceID": "11111",
                "fileFormat": "Jpeg",
                "shotTime": "20200710143217",
                "width": 896,
                "height": 700,
                "data": "图片数据"
            }
        ]
        "timestamp" : 1564735558
    }
    
}
```

实际数据参考对比
```
[
    {
        "faceID": "111110220200710143217001770600178",
        "infoKind": 1,
        "sourceID": "11111022020071014321700177",
        "deviceID": "11111",
        "shotTime": "512",
        "leftTopX": 369,
        "leftTopY": 369,
        "rightBtmX": 749,
        "rightBtmY": 707,
        "faceAppearTime": "20200710143217",
        "faceDisAppearTime": "20200710143217",
        "similaritydegree": 0,
        "data": [
            {
                "imageID": "11111022020071014321700177",
                "eventSort": 10,
                "fileFormat": "Jpeg",
                "shotTime": "20200710143217",
                "width": 1920,
                "height": 1264,
                "data": "图片数据"
            },
            {
                "imageID": "11111022020071014321700180",
                "eventSort": 10,
                "fileFormat": "Jpeg",
                "shotTime": "20200710143217",
                "width": 896,
                "height": 700,
                "data": "图片数据"
            }
        ]
    }
]
```