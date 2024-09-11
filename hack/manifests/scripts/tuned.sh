#!/bin/sh

apt update -qq
apt install -y -qq tuned

# 系统优化参数

# 查看当前操作系统的tuned策略
tuned-adm list

mkdir /etc/tuned/balanced-tidb-optimal

# Current active profile: balanced 表示当前操作系统的 tuned 策略使用 balanced，建议在当前策略的基础上添加操作系统优化配置
# 在现有的 balanced 策略基础上添加操作系统优化配置。
cat > /etc/tuned/balanced-quickon-optimal/tuned.conf <<EOF
[main]
include=balanced

[cpu]
governor=performance

[vm]
transparent_hugepages=never
EOF

# 应用新的 tuned 策略
tuned-adm profile balanced-quickon-optimal
