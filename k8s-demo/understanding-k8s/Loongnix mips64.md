[TOC]

# 龙芯mips64

# 社区适配情况：
龙芯平台已适配了MIPS下的loongnix-Server以及Debian10对应的版本：nodejs-v12.16.3，LoongArch下的Loongnix-20.loongarch64桌面系统以及Loongnix-server-20.loongarch64服务器系统对应的版本: nodejs-v14.16.1，并将持续维护，力争为用户提供好用的开发环境。
- mips:
  loongnix-Server
  Debian10
- LoongArch:
  Loongnix-20.loongarch64桌面系统
  Loongnix-server-20.loongarch64服务器系统


# 源：
## yum源修改：
```bash
yum install epel-release -y

vi  /etc/yum.repos.d/Loongnix-Base.repo
org>cn
sed -i 's/org/cn/g' /etc/yum.repos.d/Loongnix-Base.repo 
yum clean all 
yum makecache

yum install epel-release -y && cp /etc/yum.repos.d/epel.repo /etc/yum.repos.d/epel-test.repo && sed -i s/ftp.loongnix.org/10.2.5.28/g /etc/yum.repos.d/epel-test.repo && yum makecache && yum install git which gcc g++ libatomic gpg tar openssl11 -y

yum install epel-release -y && yum makecache
```



## apt源：
```bash
echo "deb http://os.loongnix.org/mirrors/debian/debian/ buster main" >> /etc/apt/sources.list

echo "deb http://mirrors.163.com/debian/ buster main contrib non-free" > /etc/apt/sources.list && echo "deb http://mirrors.163.com/debian/ buster-updates main contrib non-free" >> /etc/apt/sources.list && echo "deb http://mirrors.163.com/debian/ buster-backports main contrib non-free" >> /etc/apt/sources.list && echo "deb http://mirrors.163.com/debian-security buster/updates main contrib non-free" >> /etc/apt/sources.list
```

## Loongnix-20.mips64el系统源地址
```
deb [trusted=yes] http://ftp.loongnix.org/os/loongnix/20/mips64el/ DaoXiangHu-testing main contrib non-free 
deb-src [trusted=yes] http://ftp.loongnix.org/os/loongnix/20/mips64el/ DaoXiangHu-testing main contrib non-free
Loongnix-server-1.7.2007 服务器系统源地址
http://ftp.loongnix.org/os/loongnix-server/1.7/
```

## 龙芯NPM源
源地址1：http://npm.loongnix.cn:4873
源地址2：http://registry.loongnix.cn:4873


# harbor
用户名: loongsoncloud
密码: loongson@SYS3

执行以下命令编辑/etc/docker/daemon.json，增加insecure-registries的配置，重新加载并重启docker使配置生效
```bash
mkdir -p /etc/docker/
tee /etc/docker/daemon.json <<-‘EOF’
{
“insecure-registries”:[“harbor.loongnix.cn”]
}
EOF
sudo systemctl daemon-reload
sudo systemctl enable docker
sudo systemctl restart docker
```


# 参考资料
http://doc.loongnix.cn/web/#/50?page_id=146 龙芯Docker安装手册
http://ftp.loongnix.cn/ 龙芯开源社区ftp下载站点
http://www.loongnix.cn/index.php/Container-Registry Loongarch架构的软件仓库站点
http://ask.loongnix.cn/?/search/q-bm9kZWpz#all
