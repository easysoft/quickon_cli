# Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

FROM ysicing/k8s as fetch

WORKDIR /root

RUN wget https://ghproxy.com/https://github.com/easysoft/quickon_cli/releases/download/v1.2.11/qcadmin_linux_amd64 && \
  chmod +x qcadmin_linux_amd64

FROM ysicing/k8s

COPY --from=fetch /root/qcadmin_linux_amd64 /qcadmin_linux_amd64

COPY entrypoint.sh /entrypoint.sh

CMD [ "/entrypoint.sh" ]
