# 车辆识别适配器说明

## 功能

- 设备状态上报
- 上报设备信息
- 上报通行记录
- 下发白名单
- 删除白名单
- 远程开门

## 配置文件说明

```
{
    "httpAddr": ":9082",  // 对外开放api端口
    "pushApiAddr" : "",
    "logLevel" : "debug", // 日志 输出等级
    "bridge" : "",          // 桥接模式
    "controlMode": "mqtt",  // 下发指令方式
    "mqtt" : {
        "mode" : "special",     // special: 针对服务下发，device: 针对设备下发
        "addr" : "127.0.0.1:1883", // mqtt 服务地址
        "clientID" : "carguardAdapter",  // 默认使用的 clientid
        "username": "name",     // mqtt 默认使用的 username
        "password": "password", // mqtt 默认使用的 password
        "topicPrefix" : "request_10000" // 默认topic
    }
}
```

## 数据上报

涉及到 接口调用的统一使用格式：

 HTTP接口 (Post): http://ip:port/接口路径

 Content-Type: application/json;charset=UTF-8

## 字段介绍

* 涉及图片的都是 jpg/ jpeg 格式的
* base64编码图片 数据头不需要增加 base64，以/9j/开头
* 涉及到的枚举类型看文档最后

### 上报 方式

1. http 请求业务开放的接口

> 使用同一个接口，上报内容，业务自行通过 topic 判断上报类型

2. kafka topic publish 

> kafka 上报内容只有 payload ，业务订阅不同的topic

### 上报设备状态 消息结构 (注意 kafka: 只有payload)

```
{
    "topic" : "deviceStatusCar",
    "key" : "1111",
    "payload" : {
        "parkId": "10000",
        "online" : 0, // 0: 离线 ， 1: 在线
        "timestamp" : 1564735558
    }
    
}
```


### 上报设备信息 消息结构 (注意 kafka: 只有payload)

```
{
    "topic" : "deviceInfoCar",
    "key" : "1111",
    "payload" : {
        "parkId" : "10000",
        "mac" : "xxxxxx",
        "ip" : "192.168.1.13",
        "deviceCode" : "1111",
        "name" : "设备1"
    }
    
}
```
### 通行识别记录 消息结构

```
{
    "topic" : "plateVerification",
    "key" : "1111",
    "payload" : {
        "parkId": "10000",
        "recordId": "ec7ede33-6c91-4aee-9e6b-a859046b8c91",  // 记录id
        "picTime": "2020-06-01T12:00:00",   // 拍照时间
        "plateNo": "浙A12345",  // 车牌
        "confidence": 99,   // 可信度
        "vehicleType": 0,   // 车辆类型
        "vehicleColor": 0,  // 车辆颜色
        "plateType": 0,     // 车牌类型
        "plateColor": 0,    // 车牌颜色
        "picInfo": [{   // 第一张为全景图
            "type": 1,
            "size": 1024,
            "data": "Y2guY29tFw==",
            "url": "http://aliyun.com/park.jpg"
            }, {        // 第二张为车牌特拍图
            "type": 1,
            "size": 1024,
            "data": "Y2guY29tFw==",
            "url": "http://aliyun.com/park.jpg"
            }]
        "timestamp" : 1564735558
    }
    
}
```

## api接口 业务请求此中间件

接口： http://ip:port/api/carguard

### 下发白名单

请求示例：
```
{
    "topic" : "plateAuthorized",
    "key" : "1111",
    "payload" : {
        "id" : "123", // 白名单id
        "plateNo" : "京xxxx",  // 车牌号 范围[1,16]
        "ownerName" : "住户1", // 住户姓名 范围[1,16]
        "phoneNo" : "13800000", // 手机号  范围[1,16]
        "beginTime" : "",  // 开始 时间戳
        "endTime" : "", // 结束 时间戳
        "remark" : "备注", // 范围[1,128]
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
    "topic" : "plateAuthorizedCancel",
    "key" : "1111",
    "payload" : {
        "id" : "111" // 白名单记录ID
        "plateNo" : "京xxxx",  // 车牌号 范围[1,16]
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
    "topic" : "deviceOpenCar",
    "key" : "1111",
    "payload" : {
        "parkId" : "10000",
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

#### 错误码



#### 车牌种类

| 车牌种类 | 说明 |
| --   | -- | 
| 0	 | 大型汽车 |
| 1	 | 小型汽车 |
| 2	 | 使馆汽车 |
| 3	 | 领馆汽车 |
| 4	 | 境外汽车 |
| 5	 | 外籍汽车 |
| 6	 | 普通摩托车号牌（两、三轮摩托车号牌） |
| 7	 | 轻便摩托车 |
| 8	 | 使馆摩托车 |
| 9	 | 领馆摩托车 |
| 10 | 	外摩托车 |
| 11 | 	外籍摩托车 |
| 12 | 	低速车号牌（农用运输车） |
| 13 | 	拖拉机 |
| 14 | 	挂车 |
| 15 | 	教练汽车 |
| 16 | 	教练摩托车 |
| 17 | 	临时入境汽车 |
| 18 | 	临时入境摩托车 |
| 19 | 	临时行驶车 |
| 20 | 	警用汽车 |
| 21 | 	警用摩托 |
| 22 | 	原农机号牌 |
| 23 | 	香港入出境车 |
| 24 | 	澳门入出境车 |
| 25 | 	武警号牌 |
| 26 | 	军队号牌 |
| 27 | 	行人 |
| 28 | 	非机动车 |
| 29 | 	大型新能源车牌汽车 |
| 30 | 	小型新能源车牌车 |
| 31 | 	中型车 |
| 32 | 	试验汽车 |
| 33 | 	试验摩托车 |
| 98 | 	其他 |
| 99 | 	未知 |

#### 车辆类型

| 车辆类型 | 说明 |
| --   | -- | 
| 0	 | 三轮车 |
| 1	 | 大客车 |
| 2	 | 中型车 |
| 3	 | 小型车 |
| 4	 | 大型车 |
| 5	 | 二轮车 |
| 6	 | 摩托车 |
| 7	 | 拖拉机 |
| 8	 | 农用货车 |
| 9	 | 轿车 |
| 10 | 	SUV |
| 11 | 	面包车 |
| 12 | 	小货车 |
| 13 | 	中巴车/中型客车 |
| 14 | 	大客车/大型客车 |
| 15 | 	大货车/大型货车 |
| 16 | 	皮卡车 |
| 17 | 	MPV 商务车 |
| 18 | 	跑车 |
| 19 | 	微型轿车 |
| 20 | 	两厢轿车 |
| 21 | 	三厢轿车 |
| 22 | 	轻型客车 |
| 23 | 	中型货车 |
| 24 | 	挂车 |
| 25 | 	槽罐车 |
| 26 | 	洒水车 |
| 998 |	其他 |
| 999 |	未知 |


#### 通用颜色

| 通用颜色 | 说明 |
| --   | -- | 
| 0	 | 黑色 | 
| 1	 | 白色 | 
| 2	 | 灰色 | 
| 3	 | 红色 | 
| 4	 | 蓝色 | 
| 5	 | 黄色 | 
| 6	 | 橙色 | 
| 7	 | 棕色 | 
| 8	 | 绿色 | 
| 9	 | 紫色 | 
| 10 | 	青色 | 
| 11 | 	粉色 | 
| 12 | 	透明 | 
| 13 | 	银白 | 
| 14 | 	深色 | 
| 15 | 	浅色 | 
| 16 | 	无色 | 
| 17 | 	黄绿双色 | 
| 18 | 	渐变绿色 | 
| 99 | 	其他 | 
| 100 |	未知 | 

#### 状态码

| 状态码 | 说明 |
| --   | -- | 
| 200 |	success |
| 101 |	common error |
| 301 |	invalid param |

## 设备模型示例

### 设备指令的定义

```
{
    url: "请求的完成连接",
    method: "请求的方法",
    body : "请求消息体的字符串",
    payload : "mqtt请求的消息体"
}
```

### 远程开门http示例

```
url : http://ip:port/LAPI/V1.0/ParkingLots/Entrances/Lanes/0/GateControl
method : POST
body: { "Command": 0}
```

### 远程开门mqtt示例

```
{
    "requestId": "202006171129000001",
    "version": "1.0",
    "parkId": "10000",
    "deviceId": "1001",
    "type": 5,
    "params": {}
}
```

### 白名单新增 http 示例

```
url : http://ip:port/LAPI/V1.0/ParkingLots/Vehicles/AllowList
method : POST
body: {
    "Num": 2,
    "AllowListInfo": [
        {
            "AllowID": 123,
            "PlateNo": "浙A12345",
            "OwnerName": "张三",
            "PhoneNo": "123456789",
            "BeginTime": "1592807453",
            "EndTime": "1592807453",
            "Remarks": "校内车",
        },
        {
            "AllowID": 234,
            "PlateNo": "京A12345",
            "OwnerName": "李四",
            "PhoneNo": "987654321",
            "BeginTime": "1592807453",
            "EndTime": "1592807453",
            "Remarks": "备注",
        }
    ]
}
```

### 白名单新增 MQTT示例

其中 listType 为 1 白名单, 2 为 黑名单
mode ：0 为新增 ， 1 为删除

```
{
    "requestId": "202006171129000001",
    "version": "1.0",
    "parkId": "10000",
    "deviceId": "1001",
    "type": 4,
    "params": {
    "listType": 1,
    "mode": 0,
    "num": 2,
    "listInfo": [{
        "plateNo": "浙A12345",
        "startTime": "1589951640",
        "endTime": "1589951740"
    },
    {
        "plateNo": "浙A12346",
        "startTime": "1589951640",
        "endTime": "1589951740"
    }]
    }
}
```

### 白名单删除 http 示例

```
url : http://ip:port/LAPI/V1.0/ParkingLots/Vehicles/AllowList/${ID}
method : DELETE
body: null
```
