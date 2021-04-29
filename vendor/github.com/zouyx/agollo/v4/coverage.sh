#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

workdir=.cover
profile="$workdir/cover.out"
mode=count

generate_cover_data() {
    rm -rf "$workdir"
    mkdir "$workdir"

    for pkg in "$@"; do
        f="$workdir/$(echo $pkg | tr / -).cover"
        go test -covermode="$mode" -coverprofile="$f" "$pkg"
    done

    echo "mode: $mode" >"$profile"
    grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"
}

show_html_report() {
    go tool cover -html="$profile" -o="$workdir"/coverage.html
}

show_csv_report() {
    go tool cover -func="$profile" -o="$workdir"/coverage.csv
}

push_to_coveralls() {
    echo "Pushing coverage statistics to coveralls.io"
    # ignore failure to push - it happens
    $HOME/gopath/bin/goveralls -coverprofile="$profile" \
                               -service=travis-ci  || true
}

generate_cover_data $(go list ./...)
show_csv_report

case "$1" in
"")
    ;;
--html)
    show_html_report ;;
--coveralls)
    push_to_coveralls ;;
*)
    echo >&2 "error: invalid option: $1"; exit 1 ;;
esac