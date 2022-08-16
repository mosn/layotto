module mosn.io/layotto/components

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/alicebob/miniredis/v2 v2.16.0
	github.com/aliyun/aliyun-oss-go-sdk v2.2.0+incompatible
	github.com/apache/dubbo-go-hessian2 v1.10.2
	github.com/apolloconfig/agollo/v4 v4.2.0
	github.com/aws/aws-sdk-go-v2 v1.16.4
	github.com/aws/aws-sdk-go-v2/config v1.15.9
	github.com/aws/aws-sdk-go-v2/credentials v1.12.4
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.14
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.10
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/deckarep/golang-set v1.8.0
	github.com/go-redis/redis/v8 v8.8.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-zookeeper/zk v1.0.2
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/hashicorp/consul/api v1.3.0
	github.com/jarcoal/httpmock v1.2.0
	github.com/jinzhu/copier v0.3.6-0.20220506024824-3e39b055319a
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/minio/minio-go/v7 v7.0.15
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/qiniu/go-sdk/v7 v7.11.1
	github.com/spf13/afero v1.2.2 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tencentyun/cos-go-sdk-v5 v0.7.33
	github.com/valyala/fasthttp v1.26.0
	go.beyondstorage.io/services/hdfs v0.3.0
	go.beyondstorage.io/v5 v5.0.0
	go.etcd.io/etcd/api/v3 v3.5.0
	go.etcd.io/etcd/client/v3 v3.5.0
	go.etcd.io/etcd/server/v3 v3.5.0
	go.mongodb.org/mongo-driver v1.8.0
	go.uber.org/atomic v1.7.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5 // indirect
	google.golang.org/grpc v1.38.0
	mosn.io/api v1.0.0
	mosn.io/mosn v1.0.2-0.20220624093205-6d37b64a8f06
	mosn.io/pkg v1.0.0
)

replace github.com/klauspost/compress => github.com/klauspost/compress v1.13.1
