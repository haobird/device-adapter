# logger

## 使用说明

### 简单使用

```
logger.Debug("dddd")
```
### 输出到文件简单配置

```
# 初始化默认配置
cfg := logger.DefaultDefinition()
logger.InitWithConfig(cfg)
logger.Debug("dddd")
```

### 自定义输出到文件配置

```
cfg := logger.Cfg{
		Writers:       "stdout,file",
		Level:         "DEBUG",
		File:          "log/chassis.log",
		FormatText:    false,
		Color:         false,
		RollingPolicy: RollingPolicySize,
		RotateDate:    1,
		RotateSize:    10,
		BackupCount:   7,
	}

logger.InitWithConfig(&cfg)
logger.Debug("dddd")

```
