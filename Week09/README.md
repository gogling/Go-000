
### 作业：
Q : 用 Go 实现一个 tcp server ，用两个 goroutine 读写 conn，两个 goroutine 通过 chan 可以传递 message，能够正确退出

### 回答：
```
➜  Week09 git:(main) ✗ go run server.go

starting the server ...
all server stoped, program quit
biz : 2 : start do biz ... 
biz : 1 : start do biz ... 
listener : 1 : start listen ... 
listener : 2 : start listen ... 

listener : 1 : read conn cnt : hello 
bizer : 2 : read msg from chan : hello read by listener: 1
 
listener : 2 : read conn cnt : world 
bizer : 1 : read msg from chan : world read by listener: 2 

read error or finished EOF
read error or finished read tcp 127.0.0.1:3000->127.0.0.1:54494: use of closed network connection
^C
got signal : interrupt 
stoping listen... 
biz : 1 : stop do biz ... 
biz : 2 : stop do biz ... 
2021/01/28 23:46:29 accept error:accept tcp 127.0.0.1:3000: use of closed network connection

```