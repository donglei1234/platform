## Mail API说明文档

### 发送邮件
```http request
POST  /v1/mail/send
```

### 请求参数
 * Headers: token: 管理员token
### body: 
```json
{
  "mail": {
    "theme": "13",
    "body": "好吧",
    "expire": 0 ,
    "rewards": [
      {
        "id": 1,
        "num": 21
      },
      {
        "id": 2,
        "num": 21
      }
    ]
  },
  "targets": [
  ]
}
```
### 字段介绍：
 *  theme: 主题  
 *  body: 内容  
 *  expire: 过期时间(小时),0表示永不过期
 *  rewards: 奖励列表
 *  targets: 目标玩家列表 为空默认是全服邮件