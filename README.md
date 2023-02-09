## 目录结构
```
.
├── cache             缓存
├── controllers       控制层      
│    ├── api          客户端或前端api接口
│    ├── render       客户端或前端获取到处理过的结果
│    └── router.go    路由定义
├── docs              文档目录     
│    ├── doc          api接口、说明文档等
│    └── sql          sql相关的初始脚本
├── middleware        middleware目录
├── model             数据模型目录
├── pkg               外部服务相关目录
├── repositories      数据处理层     
├── services          服务逻辑
├── build.sh          编译脚本
├── xxx.yaml          配置文件
└── main.go           程序入口 
```

## 环境初始化
```
go版本 >= 1.16
export GOPROXY=https://goproxy.io,direct
go mod download
tools/gofmt.sh
```

## 编译及运行
```
参考脚本：
sh build.sh
或：
go run main.go
```