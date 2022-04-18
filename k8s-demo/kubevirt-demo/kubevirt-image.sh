#!/bin/bash

set -xe

# kubevirt组件镜像打包，包括kubevurt+cdi+yaml
# 参数:修改kubevirt及cdi版本参数

KUBEVIRT_VERSION=v0.35.0
CDI_VERSION=v1.26.0
# 以kubevirt的版本号为准
SAVE_TAR_NAME=kubevirt-base-image-${KUBEVIRT_VERSION}.tar


# 清空所有virt及cdi镜像
docker images | grep -E 'virt|cdi' | awk '{print $1":"$2}' | xargs -I {} docker rmi {}

docker pull kubevirt/virt-operator:${KUBEVIRT_VERSION}
docker pull kubevirt/virt-api:${KUBEVIRT_VERSION}
docker pull kubevirt/virt-controller:${KUBEVIRT_VERSION}
docker pull kubevirt/virt-handler:${KUBEVIRT_VERSION}
docker pull kubevirt/virt-launcher:${KUBEVIRT_VERSION}

docker pull kubevirt/cdi-controller:${CDI_VERSION}
docker pull kubevirt/cdi-importer:${CDI_VERSION}
docker pull kubevirt/cdi-cloner:${CDI_VERSION}
docker pull kubevirt/cdi-uploadproxy:${CDI_VERSION}
docker pull kubevirt/cdi-apiserver:${CDI_VERSION}
docker pull kubevirt/cdi-uploadserver:${CDI_VERSION}
docker pull kubevirt/cdi-operator:${CDI_VERSION}

# 重新打tag:
docker images | grep virt-api | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/virt-api:${KUBEVIRT_VERSION}
docker images | grep virt-launcher | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/virt-launcher:${KUBEVIRT_VERSION}
docker images | grep virt-controller | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/virt-controller:${KUBEVIRT_VERSION}
docker images | grep virt-handler | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/virt-handler:${KUBEVIRT_VERSION}
docker images | grep virt-operator | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/virt-operator:${KUBEVIRT_VERSION}

docker images | grep cdi-controller | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/cdi-controller:${CDI_VERSION}
docker images | grep cdi-importer | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/cdi-importer:${CDI_VERSION}
docker images | grep cdi-cloner | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/cdi-cloner:${CDI_VERSION}
docker images | grep cdi-uploadproxy | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/cdi-uploadproxy:${CDI_VERSION}
docker images | grep cdi-apiserver | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/cdi-apiserver:${CDI_VERSION}
docker images | grep cdi-uploadserver | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/cdi-uploadserver:${CDI_VERSION}
docker images | grep cdi-operator | awk '{print $3}' | xargs -I {} docker tag {} kubevirt/cdi-operator:${CDI_VERSION}

# 删除原始tag镜像
docker images | grep docker.io | awk '{print $1":"$2}' | xargs -I {} docker rmi {}

# 全部镜像打完tag后统一打成一个包：
docker images | grep -E 'virt|cdi' | awk '{print $3}' | sed ':a;N;$!ba;s/\n/ /g' | xargs -I {} echo {} docker save -o ${SAVE_TAR_NAME} {} 

# lz4压缩
lz4 -9 ${SAVE_TAR_NAME}

# 下载yaml
wget https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-operator.yaml
wget https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-cr.yaml

wget https://github.com/kubevirt/containerized-data-importer/releases/download/${CDI_VERSION}/cdi-cr.yaml
wget https://github.com/kubevirt/containerized-data-importer/releases/download/${CDI_VERSION}/cdi-operator.yaml
