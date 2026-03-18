# OpenAPI 下 declarative comment 格式说明

## 端点 handler 注释格式

例子：

```go
// List godoc
// @Summary 	获取指定汉化组的成员列表
// @Description 获取指定汉化组的成员列表，注意当列表为空，会返回 null 而不是空数组
//
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		team_id query int true "汉化组 ID"
//
// @Tags 		member
// @Produce 	json
// @Success 	200 {object} []*value.MemberInfo
//
// @Router 		/api/members [get]
func ListMembers(appState *state.AppState) {
    // ...
}
```

上述的代码将会被 swag init 命令转换为 Swagger 格式的文档，务必严格遵守。

注意，不需要给出 @Failure 的结果。
