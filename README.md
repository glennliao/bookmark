# bookmark 在线书签 
最新版本

![](https://img.shields.io/github/v/tag/glennliao/bookmark)


## 特性
### bookmark 书签
- [x] 分类新增、编辑、删除
- [x] 书签新增、编辑、删除
- [x] 书签搜索
- [x] html 书签导入
### notes 笔记
- [ ] 新增note
- [ ] note搜索

## 截图

![](./screenshot/bookmark1.png)
![](./screenshot/bookmark2.png)
![](./screenshot/bookmark3.png)


## 使用
### 1. 二进制文件执行
从[下载页面](https://github.com/glennliao/bookmark/releases)下载对应平台

####  使用
1. 执行 ./bookmark init 创建配置文件(默认为sqlite数据库)
2. 执行 ./bookmark createUser 创建用户
3. 执行 ./bookmark 启动

### 2. docker-compose 部署

```yaml
version: "3"
services:
  bookmark:
    image: glennliao/bookmark:latest
    container_name: bookmark
    restart: always
    # 使用mysql可外部挂载配置文件 config.toml , 默认使用sqlite， 需将数据库文件挂载到 /app/bookmark.db
    #volumes:
    #  - ./config.toml:/app/config.toml  
    ports:
      - 8082:8082
```


### 3 .源码编译 部署
1. 安装 goframe cli工具
2. 编译 前端 
```bash
cd ui && pnpm i && pnpm run build:prod
```
> 生成的静态文件会打包到packed目录中

3. 编译后端
```bash
# linux/amd64
gf build -s linux -a amd64 main.go
```




## changelog

### todo
- 分类拖动排序
- 书签排序
- 用户注册
- notes

### 0.7.0 (2023-07-04)
- feat 可自定义上传书签图标


### 0.6.0 (2023-07-02)
- feat 书签搜索时显示分类信息
- feat github action 打包
- fix 界面显示bug

### 0.5.0 (2023-06-23)
- feat 书签导入
- fix 书签多级目录显示
- perf 书签搜索优化

### 0.4.3 (2023-06-17)
- fix 分类无法删除
- fix 新增子分类导致父级消失 (编辑后未清空数据)

### 0.4.1 (2023-06-10)
- 增加github action 发布docker


### 0.4 (2023-06-10)
- 分类编辑，删除
- 书签编辑，删除
- 书签搜索
- 添加从浏览器收藏夹快速添加的支持
- 使用apijson-go v0.2.0-beta4 以上版本

