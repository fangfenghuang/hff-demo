1、  apk –no-cache
2、  rm –rf /var/lib/apk/*
3、  pip install –no-cache-dir
4、  run rm –rf `find / -name ‘*.pyc’`
5、  禁止使用docker commit，使用dockerfile
6、  镜像layer不得超过6层
7、  &&\ 可以把run、copy等合成一个
8、  基础镜像可以处理成scratch,使用2个阶段构建镜像
