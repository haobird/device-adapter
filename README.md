# device-adapter

## 项目启动

make start-boot
或者
docker-compose up --force-recreate --build 

## 门禁适配器

### 功能

- 上报通行记录
- 下发白名单
- 删除白名单
- 读取设备信息

### 需求

- 读取 yml 配置文件
- 支持 多个设备 同时在线
- 解析设备数据，提取标准信息
- 下发指令，根据标准信息 封装 对应的设备支持的数据调用
- 支持 三种 上报/下发 方式： MQTT/API/RabbitMQ 目前只需要调用抽象方法即可


### logger包

需求：
- 支持文件日志和控制台两种输出方式
- 支持日志输出等级
- 支持配置日志文件路径
- 支持文件定期切割(logrus配合lumberjack使用)

配置文件示例：
```
{
    "Writers":       "stdout,file", // 输出支持
    "Level":         "DEBUG", // 输出等级
    "File":          "log/chassis.log", // 日志文件路径
    "FormatText":    false, // 格式: json 或 string
    "Color":         false, //  是否彩色输出
    "RollingPolicy": "size", // 文件大小单位
    "RotateDate":    1, // 滚动天数
    "RotateSize":    10, // 文件大小
    "BackupCount":   7 // 分割文件数
}
```



## 目录文件结构


├── README.md           # 说明文件
├── go.mod              # go mod
├── main.go             # 服务启动文件
├── config.yml          # 项目配置文件
├── docker-compose.yml  # 容器启动文件
├── Makefile            # 常用命令
├── html                # 页面文件
│   ├── index.html      # 首页
│   ├── welcome.html    # 欢迎页面
│   ├──                 # 
├── faceguard           # 门禁适配
│   ├── common.go       # 通用模块
├── carguard            # 车辆识别适配
│   ├──                 # 
├── zone                # docker环境配置
├── utils               # 辅助函数
│   ├── time.go         # 时间函数
│   ├── trans.go        # 转换函数
