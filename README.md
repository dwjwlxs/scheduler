# Introduction
一个用go写的分布式任务调度软件(a distributed job scheduler using golang)

# Features
- 多实例协同工作
- 延时型、间隔型和定时型三种作业
- 兼容crontab的"分 时 天 月 周"的时间定义方式
- 添加新的woker很容易
- 部署很方便
- golang的代码易读性好

# Dependency
- centos 6+
- go 1.8+
- redis
- mysql
