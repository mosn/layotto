project_path=$(pwd)
true=0
false=1

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

  protoc -I . \
    --go_out spec/proto/extension/v1 --go_opt=paths=source_relative \
    --go-grpc_out=spec/proto/extension/v1 \
    --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative \
    --p6_out spec/proto/extension/v1 --p6_opt=paths=source_relative \
    ${protos}

  # move code to the right places
  # sdk
  mv spec/proto/extension/v1/client/client_generated.go sdk/go-sdk/client/client_generated.go
  mv spec/proto/extension/v1/grpc/context_generated.go pkg/grpc/context_generated.go
  # runtime related code
  mv spec/proto/extension/v1/runtime/* pkg/runtime/

  # api plugin
  mv spec/proto/extension/v1/grpc/* pkg/grpc/
  rm -rf spec/proto/extension/v1/grpc

  # component
  mv spec/proto/extension/v1/components/* components/
  rm -rf spec/proto/extension/v1/components

}

echo "===========> Generating .pb.go code for spec/proto/extension/v1/"
# generate .pb.go for extension/v1
res=$(ls -d spec/proto/extension/v1/*/)
toGenerate=()
idx=0
for r in $res; do
  echo $r
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

echo "===========> Generating sdk & sidecar code for spec/proto/extension/v1/"
echo "${#toGenerate[*]} packages need code generate generation"
generateSdkAndSidecar "${toGenerate[*]}"

# generate .pb.go for runtime/v1
echo "===========> Generating code for spec/proto/runtime/v1/"
docker run --rm \
  -v $project_path/spec/proto/runtime/v1:/api/proto \
  layotto/protoc
