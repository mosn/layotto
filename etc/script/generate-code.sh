project_path=$(pwd)

echo "===========> Generating code for spec/proto/extension/v1/"
# generate code for extension/v1
res=$(ls -d spec/proto/extension/v1/*/)
for r in $res; do
  echo $r
  docker run --rm \
    -v $project_path/$r:/api/proto \
    layotto/protoc
done

# generate code for runtime/v1
echo "===========> Generating code for spec/proto/runtime/v1/"
docker run --rm \
  -v $project_path/spec/proto/runtime/v1:/api/proto \
  layotto/protoc
