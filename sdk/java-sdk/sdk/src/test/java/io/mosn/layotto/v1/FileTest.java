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
package io.mosn.layotto.v1;

import io.grpc.ManagedChannel;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.testing.GrpcCleanupRule;
import io.mosn.layotto.v1.mock.MyFileService;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.sdk.runtime.v1.client.RuntimeClient;
import spec.sdk.runtime.v1.domain.file.GetFileRequest;
import spec.sdk.runtime.v1.domain.file.PutFileRequest;
import spec.sdk.runtime.v1.domain.file.DelFileRequest;
import spec.sdk.runtime.v1.domain.file.ListFileResponse;
import spec.sdk.runtime.v1.domain.file.ListFileRequest;
import spec.sdk.runtime.v1.domain.file.GetMetaRequest;
import spec.sdk.runtime.v1.domain.file.GetMeteResponse;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.util.HashMap;

import static org.mockito.AdditionalAnswers.delegatesTo;
import static org.mockito.Mockito.mock;

@RunWith(JUnit4.class)
public class FileTest {
    @Rule
    public final GrpcCleanupRule              grpcCleanup = new GrpcCleanupRule();

    RuntimeGrpc.RuntimeImplBase               fileService = new MyFileService();

    private final RuntimeGrpc.RuntimeImplBase serviceImpl = mock(RuntimeGrpc.RuntimeImplBase.class,
                                                              delegatesTo(fileService));

    private RuntimeClient                     client;

    @Before
    public void setUp() throws Exception {
        String serverName = InProcessServerBuilder.generateName();
        grpcCleanup.register(InProcessServerBuilder
            .forName(serverName).directExecutor()
            .addService(serviceImpl)
            .build().start());
        ManagedChannel channel = grpcCleanup.register(
            InProcessChannelBuilder.forName(serverName).directExecutor().build());
        client = new RuntimeClientBuilder()
            .buildGrpcWithExistingChannel(channel);
    }

    // normal case
    @Test
    public void testPutFile1() throws Exception {

        PutFileRequest req = new PutFileRequest();
        req.in = new ByteArrayInputStream("hello world".getBytes());
        req.storeName = "oss";
        req.fileName = "test.log";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.putFile(req, 10000);
    }

    // miss request
    @Test(expected = IllegalArgumentException.class)
    public void testPutFile2() throws Exception {
        client.putFile(null, 10000);
    }

    // miss in stream
    @Test(expected = IllegalArgumentException.class)
    public void testPutFile3() throws Exception {

        PutFileRequest req = new PutFileRequest();
        req.storeName = "oss";
        req.fileName = "test.log";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.putFile(req, 10000);
    }

    // miss file name
    @Test(expected = IllegalArgumentException.class)
    public void testPutFile4() throws Exception {

        PutFileRequest req = new PutFileRequest();
        req.in = new ByteArrayInputStream("hello world".getBytes());
        req.storeName = "oss";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.putFile(req, 10000);
    }

    // miss store name
    @Test(expected = IllegalArgumentException.class)
    public void testPutFile5() throws Exception {

        PutFileRequest req = new PutFileRequest();
        req.in = new ByteArrayInputStream("hello world".getBytes());
        req.fileName = "test.log";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.putFile(req, 10000);
    }

    // normal case
    @Test
    public void testGetFile1() throws Exception {

        GetFileRequest req = new GetFileRequest();
        req.out = new ByteArrayOutputStream();
        req.storeName = "oss";
        req.fileName = "test.log";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.getFile(req, 10000);

        String echo = req.out.toString();
        Assert.assertEquals("get file store name oss, meta 2, file name test.log", echo);
    }

    // miss request
    @Test(expected = IllegalArgumentException.class)
    public void testGetFile2() throws Exception {
        client.getFile(null, 10000);
    }

    // miss out stream
    @Test(expected = IllegalArgumentException.class)
    public void testGetFile3() throws Exception {

        GetFileRequest req = new GetFileRequest();
        req.storeName = "oss";
        req.fileName = "test.log";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.getFile(req, 10000);

        String echo = req.out.toString();
        Assert.assertEquals("put file store name oss, meta 2, file name test.log", echo);
    }

    // miss store name
    @Test(expected = IllegalArgumentException.class)
    public void testGetFile4() throws Exception {

        GetFileRequest req = new GetFileRequest();
        req.out = new ByteArrayOutputStream();
        req.fileName = "test.log";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.getFile(req, 10000);

        String echo = req.out.toString();
        Assert.assertEquals("put file store name oss, meta 2, file name test.log", echo);
    }

    // miss file name
    @Test(expected = IllegalArgumentException.class)
    public void testGetFile5() throws Exception {

        GetFileRequest req = new GetFileRequest();
        req.out = new ByteArrayOutputStream();
        req.storeName = "oss";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.getFile(req, 10000);

        String echo = req.out.toString();
        Assert.assertEquals("put file store name oss, meta 2, file name test.log", echo);
    }

    // normal case
    @Test
    public void testDelFile1() throws Exception {

        DelFileRequest req = new DelFileRequest();
        req.storeName = "oss";
        req.fileName = "test.log";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.delFile(req, 10000);
    }

    // miss request
    @Test(expected = IllegalArgumentException.class)
    public void testDelFile2() throws Exception {
        client.delFile(null, 10000);
    }

    // miss store name
    @Test(expected = IllegalArgumentException.class)
    public void testDelFile3() throws Exception {

        DelFileRequest req = new DelFileRequest();
        req.fileName = "test.log";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.delFile(req, 10000);
    }

    // miss file name
    @Test(expected = IllegalArgumentException.class)
    public void testDelFile4() throws Exception {

        DelFileRequest req = new DelFileRequest();
        req.storeName = "oss";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.delFile(req, 10000);
    }

    // normal
    @Test
    public void testListFile1() throws Exception {

        ListFileRequest req = new ListFileRequest();
        req.storeName = "oss";

        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        ListFileResponse resp = client.listFile(req, 10000);

        Assert.assertEquals(1, resp.fileNames.length);
        Assert.assertEquals("put file store name oss, meta 2", resp.fileNames[0]);
    }

    // miss request
    @Test(expected = IllegalArgumentException.class)
    public void testListFile2() throws Exception {
        client.listFile(null, 10000);
    }

    // miss store name
    @Test(expected = IllegalArgumentException.class)
    public void testListFile3() throws Exception {

        ListFileRequest req = new ListFileRequest();
        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        ListFileResponse resp = client.listFile(req, 10000);

        Assert.assertEquals(1, resp.fileNames.length);
        Assert.assertEquals("put file store name oss, meta 2", resp.fileNames[0]);
    }

    // normal
    @Test
    public void testGetFileMeta1() throws Exception {

        GetMetaRequest req = new GetMetaRequest();
        req.storeName = "oss";
        req.fileName = "test.log";
        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        GetMeteResponse res = client.getMeta(req, 10000);

        Assert.assertEquals(100L,res.size);
        Assert.assertEquals("2021-11-22 10:24:11",res.lastModified);
        Assert.assertArrayEquals(new String[]{"v1","v2"},res.meta.get("k1"));
    }

    // miss request
    @Test(expected = IllegalArgumentException.class)
    public void testGetFileMeta2() throws Exception {
        client.getMeta(null, 10000);
    }

    // miss store name
    @Test(expected = IllegalArgumentException.class)
    public void testGetFileMeta3() throws Exception {

        GetMetaRequest req = new GetMetaRequest();
        req.fileName = "test.log";
        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        GetMeteResponse res = client.getMeta(req, 10000);

        Assert.assertEquals(100L,res.size);
        Assert.assertEquals("2021-11-22 10:24:11",res.lastModified);
        Assert.assertArrayEquals(new String[]{"v1","v2"},res.meta.get("k1"));
    }

    // miss file name
    @Test(expected = IllegalArgumentException.class)
    public void testGetFileMeta4() throws Exception {

        GetMetaRequest req = new GetMetaRequest();
        req.storeName = "oss";
        req.metaData = new HashMap<>();
        req.metaData.put("k1", "v1");
        req.metaData.put("k2", "v2");

        client.getMeta(req, 10000);
    }
}
