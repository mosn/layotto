project_path=$(pwd)
tpl_path="${project_path}/docs/template"
quickstart_path="${project_path}/docs/en/start"
quickstart_path_zh="${project_path}/docs/zh/start"
sidebar_path="${project_path}/docs/_sidebar.md"
sidebar_path_zh="${project_path}/docs/zh/_sidebar.md"
qs_ci_path="${project_path}/etc/script/test-quickstart.sh"
proto_path_extension="spec/proto/extension/v1"

true=0
false=1

needGenerateQuickstart() {
  file=$1

  # check no `@exclude` tag
  if [ $(grep "@exclude skip quickstart_generator" $file | wc -l) -eq 0 ]; then
    # check if there's a gRPC service in it
    if [ $(grep "service " $file | wc -l) -gt 0 ]; then
      return $true
    fi
  fi
  return $false
}

addQuickstartIntoCI() {
  doc=$1

  if [ $(grep $doc ${qs_ci_path} | wc -l) -eq 0 ]; then
    sed -i "" '/quickstarts_in_default="/a \
 '${doc}'\
 ' ${qs_ci_path}
  fi
}

generateQuickstart() {
  proto_path=$1
  proto_name=$2
  nickname=$3
  api_reference_url=$4

  # 0. check if it needs our help
  needGenerateQuickstart "${proto_path}/${proto_name}"
  if [ $? -eq $false ]; then
    return 0
  fi

  # 1. copy and rename the proto file
  cp "${proto_path}/${proto_name}" "${proto_path}/${nickname}"

  # 2. generate the english quickstart document
  # check if the quickstart doc already exists
  if ! test -e "${quickstart_path}/${nickname}/start.md"; then
    echo "===========> Generating the english quickstart document for ${proto_path}/${proto_name}"
    docker run --rm \
      -v ${quickstart_path}/${nickname}:/out \
      -v ${project_path}/${proto_path}:/protos \
      -v ${tpl_path}:/tpl \
      pseudomuto/protoc-gen-doc --doc_opt=/tpl/quickstart.tmpl,start.md ${nickname}

    # modify api reference url
    if ! test -z ${api_reference_url}; then
      sed -i "" "s#<!--api_reference_url-->#[API Reference](${api_reference_url})#" ${quickstart_path}/${nickname}/start.md
    fi

    # modify design doc url
    if test -e "docs/zh/design/${nickname}/design.md"; then
      sed -i "" "s#<!--design_doc_url-->#[Design doc](zh/design/${nickname}/design)#" ${quickstart_path}/${nickname}/start.md
    fi
  fi

  # 3. add the quickstart into the sidebar
  if [ $(grep "en/start/${nickname}/start" ${sidebar_path} | wc -l) -eq 0 ]; then
    sed -i "" '/quickstart_generator/a \
'"\  "'- [(Under construction) Use '${nickname}' API](en/start/'${nickname}'/start) \
' "${sidebar_path}"
  fi

  # 4. generate the chinese quickstart document
  # check if the chinese quickstart doc already exists
  if ! test -e "${quickstart_path_zh}/${nickname}/start.md"; then
    echo "===========> Generating the chinese quickstart document for ${proto_path}/${proto_name}"
    docker run --rm \
      -v ${quickstart_path_zh}/${nickname}:/out \
      -v ${project_path}/${proto_path}:/protos \
      -v ${tpl_path}:/tpl \
      pseudomuto/protoc-gen-doc --doc_opt=/tpl/quickstart_zh.tmpl,start.md ${nickname}

    # modify api reference url
    if ! test -z ${api_reference_url}; then
      sed -i "" "s#<!--api_reference_url-->#[API Reference](${api_reference_url})#" ${quickstart_path_zh}/${nickname}/start.md
    fi

    # modify design doc url
    if test -e "docs/zh/design/${nickname}/design.md"; then
      sed -i "" "s#<!--design_doc_url-->#[Design doc](zh/design/${nickname}/design)#" ${quickstart_path_zh}/${nickname}/start.md
    fi
  fi

  # 5. add the chinese quickstart into the sidebar
  if [ $(grep "zh/start/${nickname}/start" ${sidebar_path_zh} | wc -l) -eq 0 ]; then
    sed -i "" '/quickstart_generator/a \
'"\    "'- [(建设中)使用 '${nickname}' API](zh/start/'${nickname}'/start) \
' "${sidebar_path_zh}"
  fi

  # 6. add the quickstart doc into the test-quickstart.sh
  # check no `@exclude` tag
  if [ $(grep "@exclude skip ci_generator" "${proto_path}/${proto_name}" | wc -l) -eq 0 ]; then
    addQuickstartIntoCI "docs/en/start/${nickname}/start.md"
    addQuickstartIntoCI "docs/zh/start/${nickname}/start.md"
  fi

  # 7. clean up
  rm "${proto_path}/${nickname}"
}

# 1. generate docs for extension/v1
echo "===========> Generating docs for ${proto_path_extension}"
res=$(cd $proto_path_extension && ls -d *)
for directory in $res; do
  echo "===========> Generating the API reference for ${proto_path_extension}/${directory}"

  # 1.1. ignore empty directory
  if test $(ls ${proto_path_extension}/${directory}/*.proto |wc -l) -eq 0; then
    echo "[Warn] Directory ${directory} is empty. Ignore it."
    continue
  fi

  # 1.2. generate the API reference
  docker run --rm \
    -v ${project_path}/docs/api/v1:/out \
    -v ${project_path}/${proto_path_extension}/${directory}:/protos \
    -v ${tpl_path}:/tpl \
    pseudomuto/protoc-gen-doc --doc_opt=/tpl/api_ref_html.tmpl,${directory}.html

  # 1.3. generate the quickstart document
  # find all protos
  protos=$(cd ${proto_path_extension}/${directory} && ls *.proto)
  for p in ${protos}; do
    generateQuickstart "${proto_path_extension}/${directory}" "${p}" "${directory}" "https://mosn.io/layotto/api/v1/${directory}.html"
  done
done

# 2. generate docs for runtime/v1
echo "===========> Generating docs for spec/proto/runtime/v1/"
proto_path="spec/proto/runtime/v1"
# 2.1. generate the API reference
echo "===========> Generating the API reference for spec/proto/runtime/v1/"
docker run --rm \
  -v ${project_path}/docs/api/v1:/out \
  -v ${project_path}/${proto_path}:/protos \
  -v ${tpl_path}:/tpl \
  pseudomuto/protoc-gen-doc --doc_opt=/tpl/api_ref_html.tmpl,runtime.html
# 2.2. generate the quickstart doc
# find all protos
protos=$(cd ${proto_path} && ls *.proto)
for p in ${protos}; do
  nickname=$(basename ${p} | cut -d . -f1)
  generateQuickstart "${proto_path}" "${p}" "${nickname}" "https://mosn.io/layotto/api/v1/runtime.html"
done

# 3. update the api reference index
cd $project_path
sidebar_zh=docs/zh/api_reference/README.md
sidebar=docs/en/api_reference/README.md
echo "===========> Updating the API reference index"
# delete existing lines
# -i "" is for compatibility with MacOS. See https://blog.csdn.net/dawn_moon/article/details/8547408
sed -i "" '/.*: \[.*\]\(.*\)/d' $sidebar_zh
sed -i "" '/.*: \[.*\]\(.*\)/d' $sidebar
# reinsert the reference lines
for r in $res; do
  # ignore empty directory
  if test $(ls ${proto_path_extension}/${r}/*.proto |wc -l) -eq 0; then
    echo "[Warn] Directory ${r} is empty. Ignore it."
    continue
  fi
  # insert
  echo "$r: [spec/proto/extension/v1/$r](https://mosn.io/layotto/api/v1/$r.html) \n" >>$sidebar_zh
  echo "$r: [spec/proto/extension/v1/$r](https://mosn.io/layotto/api/v1/$r.html) \n" >>$sidebar
done
# delete last line
sed -i "" '$d' $sidebar_zh
sed -i "" '$d' $sidebar

# 4. update the sidebar
cd $project_path
# TODO
