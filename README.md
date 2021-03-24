# dashboard

## docker build
1. clone project

    ```bash
    cd go/src
    git clone https://github.com/JouYin-Chen/ip-rate-limit.git
    cd ip-rate-limit
    ```

2. build docker image
    ```bash
    docker build . -t ip-rate-limit
    ```

3. run Redis
    ```bash
    docker run --name redis -p 6379:6379 -d redis
    ```
4. run service
    ```bash
    docker run --name ip-rate-limit -p 3000:3000 --link redis -e REDIS_HOST=redis -e REDIS_PORT=6379 -d ip-rate-limit
    ```
## API

### CreateUser API
***Description***
  - Create a user with name, return a UNIQUE user ID

***Method***
  - POST

***Router***
```
http://localhost:3000/user
```
***Request***
```
{
    "name": <user_name>
}
```

### GetUser API
***description***
  - Get the user name by user ID

***Method***
  - GET

***Router***
```
http://localhost:3000/user?id=<user_ID>
```

## Test

***description***

執行 script 進行 create user 以及 get user name的測試, script 中會進行
1. 進行 curl POST 指令 call `CreateUser API`, 並收到response: user ID
2. 使用取得的 user ID, 進行 curl GET 指令 call `GetUser API` 獲得 userName, 並call api 40次 來驗證 IP rate limit的功能

***執行方法***
```
./test.sh
```
