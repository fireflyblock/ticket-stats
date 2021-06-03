# ticket-stats



#### 	创建

```shell
go build -o ticket main.go
```



#### 	环境变量

```shell
//指定结果写入的文件位置
export OUTPUT_FILE="./stats.out"
//lotus                                                                               
export LOTUS_HOST="http://ip:port/rpc/v0"
//lotus api sign权限的token，用于计算VRF
export LOTUS_SIGN_TOKEN="eyJhbGciOiJIUzI1NiIsInR..........."
```



#### 	运行

1. 根据时间

   ```shell
    ./ticket stats time --miner f0419945 --start 2021-05-24T00:00:00 --end 2021-05-24T16:21:00
   ```

2. 根据高度

   ```shell
    ./ticket stats epoch --miner=f0419945 --start 780000 --end 784000
   ```

   

