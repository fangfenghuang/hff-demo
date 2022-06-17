# harbor升级数据库迁移及回滚
## harbor安装
1. 下载docker-compose
1. 下载，解压Harbor离线包
2. 修改配置harbor.yml.tmpl>harbor.yml及生成证书等
3. 安装并启动：./install.sh

## harbor升级（大于v1.10.x到更新版本）
harbor不同阶段的升级方式不同，存在兼容性问题，有些版本无法直接跨版本升级，需要多阶段升级方式达到最终升级目的
自harborv2起，迁移工具全部转移到了 goharbor/prepare 这个镜像中

1. 到github下载2.4.2版本的离线包（包含镜像goharbor/prepare）
2. 准备迁移工具镜像goharbor/prepare
3. 停止harbor
4. 备份数据库和安装目录
5. 解压离线包到安装目录
6. 移动旧的配置文件到新的安装目录下
7. 数据迁移：docker run -it --rm -v /:/hostfs goharbor/prepare:v2.4.2 migrate -i /opt/harbor/harbor.yml
（此更新包括数据库模式（schema）和配置文件数据）
8. 安装并启动：./install.sh

## harbor回滚
1. 停止harbor
2. 删除安装目录和数据库文件
3. 将备份目录和备份数据库还原回安装目录和数据库目录
4. 安装并启动：./install.sh

## harbor数据库迁移及golang-migrate
数据库升级是自动完成的，用户手动执行升级配置文件的命令行工具包。此工具包与Harbor一同发布，被包含在goharbor/prepare镜像中。

每次启动Harbor实例时，它的数据库模式都是自动升级的，其原理为：Harbor在每次启动时都会调用第三方库 “golang-migrate”，它会检测当前数据库模式的版本，如果实例的版本比当前数据库的版本更高，则会自动升级。

```bash
$ docker run -v /:/hostfs goharbor/prepare:v2.0.0 migrate -i ${harbor.yml路径}
```
