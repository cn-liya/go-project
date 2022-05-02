### 临时头像链接转存CDN
- Type: Quorum
- Name: avatar_to_cdn
- Arguments:
  + x-delivery-limit = 4 (Number)
```shell
curl -i -u guest:guest -XPUT http://localhost:15672/api/queues/%2F/avatar_to_cdn \
-d'{"durable":true,"arguments":{"x-queue-type":"quorum","x-delivery-limit":4}}'
```

