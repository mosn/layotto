project_path=$(pwd)
true=0
false=1

# don't build layotto/protoc if it already exists
CON=$(docker image ls 'layotto/protoc:latest' | wc -l)
if [ $CON -eq 1 ]; then
  # 1 means that the image does not exist
  docker build -t layotto/protoc docker/proto
fi

# check if protoc-gen-p6 has been installed
CON=$(which protoc-gen-p6 | wc -l)
if [ $CON -eq 0 ]; then
  # 0 means that protoc-gen-p6 does not exist
  go install github.com/seeflood/protoc-gen-p6@latest
fi

needGenerate() {
  file=$1

  # check no `@exclude` tag
  if [ $(grep "@exclude code generator" $file | wc -l) -eq 0 ]; then
    # check if there's a gRPC service in it
    if [ $(grep "service " $file | wc -l) -gt 0 ]; then
      return $true
    fi
  fi
  return $false
}

generateSdkAndSidecar() {
  protos=$*

  if test ${#protos[@]} -eq 0; then
    return 0
  fi

  mkdir _output/tmp

  protoc -I . \
    --go_out . --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative \
    --p6_out _output/tmp --p6_opt=paths=source_relative \
    ${protos}

  # move code to the right places
  # sdk
  mv _output/tmp/client/client_generated.go sdk/go-sdk/client/client_generated.go
  mv _output/tmp/grpc/context_generated.go pkg/grpc/context_generated.go
  # runtime related code
  mv _output/tmp/runtime/* pkg/runtime/

  # api plugin
  mv _output/tmp/grpc/* pkg/grpc/
  #  rm -rf _output/tmp/grpc

  # component
  mv _output/tmp/components/* components/
  #  rm -rf _output/tmp/components
  rm -rf _output/tmp

}

echo "===========> Generating .pb.go code for spec/proto/extension/v1/"
# generate .pb.go for extension/v1
res=$(ls -d spec/proto/extension/v1/*/)
toGenerate=()
idx=0
for r in $res; do
  #  echo $r
  docker run --rm \
    -v $project_path/$r:/api/proto \
    layotto/protoc

  protos=$(ls ${r}*.proto)
  for proto in ${protos}; do
    # check if it needs sdk & sidecar generation
    needGenerate "${proto}"
    if test $? -eq $true; then
      echo "${proto} needs code generation."
      toGenerate[${idx}]=${proto}
      idx=$((idx + 1))
    fi
  done
done
echo "${#toGenerate[*]} packages need code generate generation"

echo "===========> Generating sdk & sidecar code for spec/proto/extension/v1/"
generateSdkAndSidecar "${toGenerate[*]}"

# generate .pb.go for runtime/v1
echo "===========> Generating code for spec/proto/runtime/v1/"
docker run --rm \
  -v $project_path/spec/proto/runtime/v1:/api/proto \
  layotto/protoc
