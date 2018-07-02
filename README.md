# b-bit MinHash(b=1)の実験
大量の乱数を生成し，それらをb-bit MinHashにより生成されたスケッチとみなして総当たりによる類似度推定にどのくらい時間がかかるかベンチマークをとった

## 各種ファイル説明
- gen-rands.sh
`$ openssl rand -hex 32`によって100万件の乱数を生成するスクリプト

- randsXX.txt
`gen-rands.sh`により生成された128bit長の乱数がそれぞれXX件分記述されているテキストファイル

- hdsearch_test.go
ベンチマークのコード.
`BenchmarkHdAllSearch1()`は，単純に全ての乱数に対するハミング距離を計算していくアルゴリズム
`BenchmarkHdAllSearch2()`は，Jaccard係数の閾値を元にハミング距離の閾値を逆算し，ハミング距離の計算過程で閾値を超えた場合に処理を打ち切ることで高速化を目論んだアルゴリズム

- bench.log
ベンチマークの結果

- benchstat.log
`bench.log`を元に計算した各アルゴリズムの結果の平均実行時間と誤差
