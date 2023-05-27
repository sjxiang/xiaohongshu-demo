

# xiaohongshu-demo


example 示范

Swagger 生成 API 文档
    

# GIN 入门

文档 https://gin-gonic.com/zh-cn/docs/



# GORM 入门

简介

    功能丰富

        支持多种数据库，包括 MySQL
        支持简单查询、支持事务，也支持关联关系（一对多，多对多，多对一）
            不建议使用 orm 框架的关联关系，关联关系应该是业务上维护的
            体现就是数据库上的外键，但互联网业务禁止使用
        支持钩子 hook
            换句话就是中间件
        支持自动迁移工具
            没有建表，会帮你创建
            建好表的话，如果表结构发生变化，会帮你同步
            开发环境可以使用，但生产环境，表结构发生变化，要经过 DBA 评审
            谨慎使用，黑盒


入门：增删改查

    安装 GROM 依赖
        安装本体
            go get -u gorm.io/gorm

        安装对应数据库的驱动，注意 GORM 做了二次封装
            go get -u gorm.io/driver/mysql

    基本使用步骤
        初始化 DB 实例
        （可选）初始化表结构
        发起查询

文档 
    https://gorm.io/zh_CN/docs/


坑
更新有个坑的地方

    "仅更新非零值字段"

    指针 nil
    整数 0


很多接口都是接收 interface{}，作为参数，根据你传入不同类似的参数，会执行不同的行为

    典型就是 Updates，传入 map、结构体，都可以
    语义过于模糊，容易写错。



# 推荐写法

    面向接口编程
    依赖注入


    调用流程

        web
            ->  service
                    ->  repo 
                            -> dao
            
                            <- dao 
                    <-  repo      
            <-  service
        web


    没有啥包方法，又快又简单，但无法写单测