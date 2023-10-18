FROM ysicing/debian

COPY ./_output/qcadmin_linux_amd64 /usr/bin/qcadmin

RUN chmod +x /usr/bin/qcadmin
