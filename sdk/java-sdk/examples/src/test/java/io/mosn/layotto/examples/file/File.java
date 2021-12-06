/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package io.mosn.layotto.examples.file;

import io.mosn.layotto.v1.RuntimeClientBuilder;
import io.mosn.layotto.v1.config.RuntimeProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import spec.sdk.runtime.v1.client.RuntimeClient;
import spec.sdk.runtime.v1.domain.file.DelFileRequest;
import spec.sdk.runtime.v1.domain.file.FileInfo;
import spec.sdk.runtime.v1.domain.file.GetFileRequest;
import spec.sdk.runtime.v1.domain.file.GetFileResponse;
import spec.sdk.runtime.v1.domain.file.GetMetaRequest;
import spec.sdk.runtime.v1.domain.file.GetMeteResponse;
import spec.sdk.runtime.v1.domain.file.ListFileRequest;
import spec.sdk.runtime.v1.domain.file.ListFileResponse;
import spec.sdk.runtime.v1.domain.file.PutFileRequest;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

/**
 * Specially
 * <p>
 * 1. add `"local:{}` to "files" node in layotto/configs/config_file.json
 * 2. start server by `./layotto start -c ../../configs/config_file.json`
 */
public class File {

    private static final Logger logger    = LoggerFactory.getLogger(File.class.getName());

    static String               storeName = "local";
    static String               fileName  = "/tmp/test.log";

    public static void main(String[] args) throws Exception {

        RuntimeClient client = new RuntimeClientBuilder()
            .withPort(RuntimeProperties.DEFAULT_PORT)
            .build();

        putFile(client);
        getFile(client);
        listFile(client);
        getFileMeta(client);
        delFile(client);
    }

    public static void putFile(RuntimeClient client) throws Exception {

        PutFileRequest request = new PutFileRequest();
        request.setStoreName(storeName);
        request.setFileName(fileName);

        Map<String, String> meta = new HashMap<>();
        meta.put("FileMode", "521");
        meta.put("FileFlag", "0777");
        request.setMetaData(meta);

        request.setIn(new ByteArrayInputStream("hello world".getBytes()));

        client.putFile(request, 3000);
    }

    public static void getFile(RuntimeClient client) throws Exception {

        GetFileRequest request = new GetFileRequest();
        request.setStoreName(storeName);
        request.setFileName(fileName);

        Map<String, String> meta = new HashMap<>();
        meta.put("k1", "v1");
        request.setMetaData(meta);

        GetFileResponse resp = client.getFile(request, 3000);

        InputStream reader = resp.getIn();

        byte[] buf = new byte[128];
        for (int len = reader.read(buf); len > 0; len = reader.read(buf)) {
            logger.info(new String(buf, 0, len));
        }
    }

    public static void delFile(RuntimeClient client) throws Exception {

        DelFileRequest request = new DelFileRequest();
        request.setStoreName(storeName);
        request.setFileName(fileName);

        client.delFile(request, 3000);
    }

    public static void listFile(RuntimeClient client) throws Exception {

        ListFileRequest request = new ListFileRequest();
        request.setStoreName(storeName);
        request.setMarker("test.log");
        request.setName("/tmp");
        request.setPageSize(10);

        ListFileResponse resp = client.listFile(request, 3000);

        for (FileInfo f : resp.getFiles()) {
            logger.info(f.getFileName());
        }
    }

    public static void getFileMeta(RuntimeClient client) throws Exception {

        GetMetaRequest request = new GetMetaRequest();
        request.setStoreName(storeName);
        request.setFileName(fileName);

        GetMeteResponse response = client.getFileMeta(request, 3000);
        logger.info(response.getLastModified());
        logger.info("" + response.getMeta().size());
    }
}
