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

# update the sidebar
cd $project_path
sidebar_zh=docs/zh/api_reference/README.md
sidebar=docs/en/api_reference/README.md
echo "===========> Updating the sidebar"
# delete existing lines
# -i "" is for compatibility with MacOS. See https://blog.csdn.net/dawn_moon/article/details/8547408
sed -i "" '/.*: \[.*\]\(.*\)/d' $sidebar_zh
sed -i "" '/.*: \[.*\]\(.*\)/d' $sidebar
# reinsert the reference lines
for r in $res; do
  echo "$r: [spec/proto/extension/v1/$r](https://mosn.io/layotto/api/v1/$r.html) \n" >> $sidebar_zh
  echo "$r: [spec/proto/extension/v1/$r](https://mosn.io/layotto/api/v1/$r.html) \n" >> $sidebar
done
# delete last line
sed -i "" '$d' $sidebar_zh
sed -i "" '$d' $sidebar


cd $project_path
# generate index for api references
#idx=$(cd docs && ls api/v1/*)
#echo $idx > docs/api/extensions.txt
