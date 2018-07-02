#!/bin/bash

# 100万件の乱数を生成する
# 1つの乱数は256bit長
# めっちゃ時間かかるので注意
for i in `seq 1 1000000`
do
  openssl rand -hex 16
done