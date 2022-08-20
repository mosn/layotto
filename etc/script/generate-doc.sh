project_path=$(pwd)
tpl_path="${project_path}/docs/template"
quickstart_path="${project_path}/docs/en/start"
quickstart_path_zh="${project_path}/docs/zh/start"
true=0
false=1

needGenerateQuickstart() {
  file=$1
  name=$2
  # check no `@exclude` tag
  if [ $(grep "@exclude quickstart generator" $file | wc -l) -eq 0 ]; then
    # check if there's a gRPC service in it
    if [ $(grep "service " $file | wc -l) -gt 0 ]; then
      return $true
    fi
  fi
  return $false
}

# 1. generate docs for extension/v1
proto_path="spec/proto/extension/v1"
echo "===========> Generating docs for ${proto_path}"
res=$(cd $proto_path && ls -d *)
for directory in $res; do
  echo "===========> Generating API reference for ${proto_path}/${directory}"
  # 1.1. generate the API reference
  docker run --rm \
    -v ${project_path}/docs/api/v1:/out \
    -v ${project_path}/${proto_path}/${directory}:/protos \
    -v ${tpl_path}:/protos/tpl \
    pseudomuto/protoc-gen-doc --doc_opt=/protos/tpl/api_ref_html.tmpl,${directory}.html

  # 1.2. generate the quickstart document
  # find all protos
  protos=$(cd ${proto_path}/${directory} && ls *.proto)
  for p in ${protos}; do
    # check if it needs our help
    needGenerateQuickstart "${proto_path}/${directory}/${p}" ${directory}
    if [ $? -eq $true ]; then
      # copy and rename the proto file
      cp "${proto_path}/${directory}/${p}" "${proto_path}/${directory}/${directory}"

      # check if the quickstart doc already exists
      if ! test -e "${quickstart_path}/${directory}/start.md"; then
        # generate the english quickstart document
        echo "===========> Generating the english quickstart document for ${proto_path}/${directory}/${p}"
        docker run --rm \
          -v ${quickstart_path}/${directory}:/out \
          -v ${project_path}/${proto_path}/${directory}:/protos \
          -v ${tpl_path}:/protos/tpl \
          pseudomuto/protoc-gen-doc --doc_opt=/protos/tpl/quickstart.tmpl,start.md ${directory}

        # TODO remove design doc if not exist.
      fi

      # check if the chinese quickstart doc already exists
      if ! test -e "${quickstart_path_zh}/${directory}/start.md"; then
        # generate the chinese quickstart document
        echo "===========> Generating the chinese quickstart document for ${proto_path}/${directory}/${p}"
        docker run --rm \
          -v ${quickstart_path_zh}/${directory}:/out \
          -v ${project_path}/${proto_path}/${directory}:/protos \
          -v ${tpl_path}:/protos/tpl \
          pseudomuto/protoc-gen-doc --doc_opt=/protos/tpl/quickstart_zh.tmpl,start.md ${directory}

        # TODO remove design doc if not exist.
      fi
      # clean up
      rm "${proto_path}/${directory}/${directory}"
    fi
  done
done

# 2. generate docs for runtime/v1
echo "===========> Generating docs for spec/proto/runtime/v1/"
docker run --rm \
  -v $project_path/docs/api/v1:/out \
  -v $project_path/spec/proto/runtime/v1:/protos \
  -v ${tpl_path}:/protos/tpl \
  pseudomuto/protoc-gen-doc --doc_opt=/protos/tpl/api_ref_html.tmpl,runtime.html

# 3. update the api reference index
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
  echo "$r: [spec/proto/extension/v1/$r](https://mosn.io/layotto/api/v1/$r.html) \n" >>$sidebar_zh
  echo "$r: [spec/proto/extension/v1/$r](https://mosn.io/layotto/api/v1/$r.html) \n" >>$sidebar
done
# delete last line
sed -i "" '$d' $sidebar_zh
sed -i "" '$d' $sidebar

# 4. update the sidebar
cd $project_path
# TODO

cd $project_path
# generate index for api references
#idx=$(cd docs && ls api/v1/*)
#echo $idx > docs/api/extensions.txt
