### 简介
daily-practice的后端golang实现
### 部署
1. go build
2. scp dailypractice user@ip:~/build
### 测试
1. 获取所有tips `curl http://localhost:9000/api/tips`
2. 登录 `curl -d '{"email": "11@ekohe.com", "password": "andrew123"}' -H "Content-Type: application/json" --request POST  localhost:9000/api/login`
