#!/bin/bash
image_name=model_extend_hw
docker_name=model_ext
docker run -itd -u root \
--ipc=host \
--device=/dev/davinci7 \
--device=/dev/davinci_manager \
--device=/dev/devmm_svm \
--device=/dev/hisi_hdc \
-v /etc/localtime:/etc/localtime \
-v /usr/local/Ascend/driver:/usr/local/Ascend/driver \
-v /usr/local/Ascend/driver/tools/hccn_tool:/usr/local/bin/hccn_tool \
-v /etc/ascend_install.info:/etc/ascend_install.info \
-v /var/log/npu:/usr/slog \
-v /usr/local/bin/npu-smi:/usr/local/bin/npu-smi \
-v /etc/hccn.conf:/etc/hccn.conf \
-v /media/data0/model_extend:/model_extend \
-p 15000:15000 \
-p 8681:8681 \
-p 8682:8682 \
-p 8683:8683 \
-p 8684:8684 \
-p 8685:8685 \
-p 10891:10891 \
-p 10892:10892 \
-p 10893:10893 \
-p 10894:10894 \
-p 10895:10895 \
-p 20041:20041 \
-p 20042:20042 \
-p 20043:20043 \
-p 20044:20044 \
-p 20045:20045 \
-p 19200:9200 \
--name ${docker_name} \
${image_name} \
/model_extend/start.sh

