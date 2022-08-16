project_path=$(pwd)

echo "===========> Generating docs for spec/proto/extension/v1/"
# generate docs for extension/v1
res=$(cd spec/proto/extension/v1/ && ls -d *)
for r in $res; do
  docker run --rm \
    -v $project_path/docs/api/v1:/out \
    -v $project_path/spec/proto/extension/v1/$r:/protos \
    -v $project_path/spec/proto/runtime/v1:/protos/tpl \
    pseudomuto/protoc-gen-doc --doc_opt=/protos/tpl/html.tmpl,$r.html
done

# generate docs for runtime/v1
echo "===========> Generating docs for spec/proto/runtime/v1/"
docker run --rm \
  -v $project_path/docs/api/v1:/out \
  -v $project_path/spec/proto/runtime/v1:/protos \
  pseudomuto/protoc-gen-doc --doc_opt=/protos/html.tmpl,runtime.html

cd $project_path
