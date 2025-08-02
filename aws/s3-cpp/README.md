# S3 C++ Stream

C++ での S3 オブジェクトのストリーム処理を行うプログラム．


## ビルド
```sh
docker compose build s3-cpp
```

## 実行
```sh
export BUCKET_NAME=xxxxxx
export OBJECT_KEY=yyyyyy

# list objects
docker compose run --rm s3-cpp ./build/s3_stream_processor $BUCKET_NAME list

# get object
docker compose run --rm s3-cpp ./build/s3_stream_processor $BUCKET_NAME get $OBJECT_KEY output/sample.mp4
```
