## genx

用于 gin HandlerFunc 代码生成, 给你的 handler 添加一个 generate 注释 , 可以自动生成代码绑定请求参数, 错误处理, 响应绑定等代码.

## 功能

- [x] 生成 api handler func 绑定函数
- [x] 生成绑定 json 请求参数到结构体
- [x] 生成绑定 json 响应
- [ ] 支持定多种参数类型(Query, Form 等)
- [ ] 支持 gin 外的 web framework
- [ ] 自定义校验器支持
- [ ] 生成公共响应包装
- [ ] 参数校验
- [ ] 错误处理及自定义处理过程

## 使用方式

### 安装

```shell
go install github.com/dengzii/genx
```

### 定义处理函数

```go
//go:generate genx handler
func TestHandler(ctx *gin.Context, request *param.TestRequest) (*param.TestResponse, error) {
// ...
return &param.LoginResponse{Token: "token"}, nil
}
```

其中 `request` 为绑定的结构体，`response` 为返回的结构体，`error` 为错误类型

所有参数均为可选, 但顺序不可变

### 生成代码

```shell
go generate
```

或者点击函数名称左侧运行按钮

### 添加路由

生成的函数以 Genx 开头

```go
gin.Default().GET("/test", GenxTestHandler)
```
