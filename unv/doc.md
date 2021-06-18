# 说明

## 功能

- 接收设备数据 中转到 mysql 的 messages库
- 响应设备心跳，维持在线状态
- 提供接口，实现白名单下发，人员库新增

## 功能聚焦

- 通行上报：人员ID/图片/时间
- 白名单下发：ID、姓名、头像、标签
- 白名单删除：ID

## 接口

统一的返回示例：
```
{
    "code":200,
    "message":"success",
    "data":{}
}
```

### 关闭设备连接

请求
```
{
  "actionCode": "00",
  "header": {
    "deviceid": "210235C3XT320B000818",
    "libid": 4,
    "requestid": "1234"
  }
}
```

### 白名单下发

请求举例
```
{
  "actionCode": "03",
  "header": {
    "deviceid": "210235C3XT320B000818",
    "libid": 4,
    "requestid": "1234"
  },
  "data": {
    "Num": 1,
    "PersonInfoList": [
      {
        "PersonID": 2222,
        "LastChange": 1602329484,
        "PersonCode": "2001",
        "PersonName": "介的陌生人",
        "Remarks": "陌生人的尝试",
        "TimeTemplateNum": 0,
        "IdentificationNum": 2,
        "IdentificationList": [
          {
            "Type": 0,
            "Number": "13022319890520222X"
          },
          {
            "Type": 99,
            "Number": "3214124"
          }
        ],
        "ImageNum": 1,
        "ImageList": [
          {
            "FaceID": 1,
            "Name": "1_1.jpg",
            "Size": 39516
          }
        ]
      }
    ]
  }
}
```

### 人员库查询

Warn: 此接口要求完全返回设备对应的响应内容

请求示例
```
{
  "actionCode": "01",
  "header": {
    "deviceid": "210235C3XT3204000391",
    "libid": 3,
    "requestid": 1234
  },
  "data": {
    "Num": 0,
    "Limit": 10,
    "Offset": 0
  }
}
```
返回示例
```
{
    "code":200,
    "message":"success",
    "data":{
        "Response": {
            "ResponseURL": "/LAPI/V1.0/PeopleLibraries/BasicInfo",
            "CreatedID": -1, 
            "ResponseCode": 0,
            "SubResponseCode": 0,
            "ResponseString": "Succeed",
            "StatusCode": 0,
            "StatusString": "Succeed",
            "Data": {
                "Num":	2,
                "LibList":	[{
                    "ID":	3,
                    "Type":	3,
                    "PersonNum":	1,
                    "MemberNum":	1,
                    "FaceNum":	1,
                    "LastChange":	1610030977,
                    "Name":	"默认员工库",
                    "BelongIndex":	""
                }, {
                    "ID":	4,
                    "Type":	4,
                    "PersonNum":	13,
                    "MemberNum":	13,
                    "FaceNum":	5,
                    "LastChange":	1610690522,
                    "Name":	"默认访客库",
                    "BelongIndex":	""
                }]
            }
	    }
    }

}

```


### 人员库新增

```
{
  "actionCode": "01",
  "header": {
    "deviceid": "210235C3XT3204000391",
    "libid": 3,
    "requestid": 1234
  },
  "data": {
    "Num": 1,
    "LibList": [
        {
        "ID": 5,
        "Type": 4,
        "LastChange": 1610681665,
        "Name": "访客库新增"
        }
    ]
    }
}
```

