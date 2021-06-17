# 题目1

​	使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。
​	10字节测试：

​	启一个redis服务器，往redis分别写入10 20 50 100 200 1k 5k 字节 value 大小，再使用redis-benchmark命令查看redis的get、set性能。此处默认用50个客户端测试。

`redis-benchmark -t set,get -n 100000 -q`命令测试，虚拟机配置：8G内存4核

结果：

```shell
10字节测试：
	SET: 160513.64 requests per second
	GET: 168350.17 requests per second
20字节测试
        SET: 180831.83 requests per second
        GET: 165837.48 requests per second
50字节测试	
        SET: 169204.73 requests per second
        GET: 167224.08 requests per second
100字节测试
        SET: 73260.07 requests per second
        GET: 162866.44 requests per second
200字节测试
        SET: 173913.05 requests per second
        GET: 169491.53 requests per second
1k字节测试
        SET: 178253.12 requests per second
        GET: 170940.17 requests per second
5k字节测试
        SET: 176056.33 requests per second
        GET: 174216.03 requests per second
```

结论：从结果中可以知道，不管是10字节还是5k大小的value对性能几乎没有影响，但是中间出现了一次“意外”100字节的时候，我猜测是网络出现抖动，再次测试就正常了。

# 题目2

写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息  , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

```shell
#未写入任何时info memory
used_memory:811576
used_memory_human:792.55K
used_memory_rss:7065600
used_memory_rss_human:6.74M
used_memory_peak:3671824
used_memory_peak_human:3.50M
total_system_memory:8185933824
total_system_memory_human:7.62G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:8.71
mem_allocator:jemalloc-3.6.0

#写入1w string类型，key为"v1~10000" 值为"a1~10000"
used_memory:1662744
used_memory_human:1.59M
used_memory_rss:7413760
used_memory_rss_human:7.07M
used_memory_peak:4514616
used_memory_peak_human:4.31M
total_system_memory:8185933824
total_system_memory_human:7.62G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:4.46
mem_allocator:jemalloc-3.6.0

#写入50w string类型，key为"v1~500000" 值为"a1~500000"
used_memory:44205984
used_memory_human:42.16M
used_memory_rss:54243328
used_memory_rss_human:51.73M
used_memory_peak:47056352
used_memory_peak_human:44.88M
total_system_memory:8185933824
total_system_memory_human:7.62G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:1.23
mem_allocator:jemalloc-3.6.0
```

结论：在零key的时候，可见redis分配了接近800k的内存，写入1w 个key的时候（value大小在5byte以内），redis分配了接近1,6M内存，写入50w个key的时候，redis分配了接近42.16M内存。大概平均每个key和value大小占88byte。由于这里的value也占了内存50w byte，但是几乎不影响，可忽略。

脚本：

```shell
#!/bin/bash
str="a"
ret=""
echo "" > readme.txt
echo "first num is $1"
for k in $(seq 1 $1)
do
        echo "set v${k} a" >> readme.txt
done
unix2dos readme.txt
cat readme.txt | redis-cli
redis-benchmark -t get,set -n 100000 -q
```

