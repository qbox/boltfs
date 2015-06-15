QBS API
========

# 协议

## 创建卷

请求包：

```
POST /v1/volumes
Content-Type: application/json
Authorization: Qiniu <MacToken>

{
	"title": <VolumeTitle>
}
```

返回包：

```
200 OK
Content-Type: application/json

{
	"id": <VolumeId>
}
```

## 删除卷

请求包：

```
DELETE /v1/volumes/<VolumeId>
Authorization: Qiniu <MacToken>
```

返回包：

```
200 OK
```
