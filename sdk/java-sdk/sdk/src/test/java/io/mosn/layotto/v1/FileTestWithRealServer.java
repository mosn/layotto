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

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.mosn.layotto.v1.grpc.ExceptionHandler;
import io.mosn.layotto.v1.grpc.GrpcRuntimeClient;
import io.mosn.layotto.v1.mock.MyFileService;
import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.sdk.runtime.v1.domain.file.GetFileRequest;
import spec.sdk.runtime.v1.domain.file.PutFileRequest;
import spec.sdk.runtime.v1.domain.file.DelFileRequest;
import spec.sdk.runtime.v1.domain.file.ListFileRequest;
import spec.sdk.runtime.v1.domain.file.GetMetaRequest;
import spec.sdk.runtime.v1.domain.file.GetMeteResponse;
import spec.sdk.runtime.v1.domain.file.GetFileResponse;
import spec.sdk.runtime.v1.domain.file.ListFileResponse;
import spec.sdk.runtime.v1.domain.file.FileInfo;

import java.io.ByteArrayInputStream;
import java.util.HashMap;
import java.util.Map;

@RunWith(JUnit4.class)
public class FileTestWithRealServer {

    private final RuntimeGrpc.RuntimeImplBase fileService = new MyFileService();

    private Server                            srv;
    private GrpcRuntimeClient                 client;

    int                                       port        = 9999;
    String                                    ip          = "127.0.0.1";

    @Before
    public void setUp() throws Exception {
        // start grpc server
        /* The port on which the server should run */
        srv = ServerBuilder.forPort(port)
            .addService(fileService)
            .intercept(new ExceptionHandler())
            .build()
            .start();

        // build a client
        client = new RuntimeClientBuilder()
            .withIp(ip)
            .withPort(port)
            .withConnectionPoolSize(4)
            .withTimeout(1000)
            .buildGrpc();
    }

    @After
    public void shutdown() {
        client.shutdown();
        srv.shutdownNow();
    }

    // normal case
    @Test
    public void testPutFile1() throws Exception {

        PutFileRequest req = new PutFileRequest();
        req.setIn( new ByteArrayInputStream("hello world".getBytes()));
        req.setStoreName("oss");
        req.setFileName("test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

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
        req.setStoreName( "oss");
        req.setFileName("test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        client.putFile(req, 10000);
    }

    // miss file name
    @Test(expected = IllegalArgumentException.class)
    public void testPutFile4() throws Exception {

        PutFileRequest req = new PutFileRequest();
        req.setIn(  new ByteArrayInputStream("hello world".getBytes()));
        req.setStoreName( "oss");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        client.putFile(req, 10000);
    }

    // miss store name
    @Test(expected = IllegalArgumentException.class)
    public void testPutFile5() throws Exception {

        PutFileRequest req = new PutFileRequest();
        req.setIn(  new ByteArrayInputStream("hello world".getBytes()));
        req.setFileName( "test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        client.putFile(req, 10000);
    }

    // normal case
    @Test
    public void testGetFile1() throws Exception {

        GetFileRequest req = new GetFileRequest();
        req.setStoreName( "oss");
        req.setFileName( "test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        GetFileResponse resp = client.getFile(req, 10000);

        byte[] buf = new byte[126];
        int len = resp.getIn().read(buf);

        String echo = new String(buf,0,len);
        Assert.assertEquals("get file store name oss, meta 2, file name test.log", echo);
    }

    // miss request
    @Test(expected = IllegalArgumentException.class)
    public void testGetFile2() throws Exception {
        client.getFile(null, 10000);
    }

    // miss store name
    @Test(expected = IllegalArgumentException.class)
    public void testGetFile3() throws Exception {

        GetFileRequest req = new GetFileRequest();
        req.setFileName( "test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        client.getFile(req, 10000);
    }

    // miss file name
    @Test(expected = IllegalArgumentException.class)
    public void testGetFile4() throws Exception {

        GetFileRequest req = new GetFileRequest();
        req.setStoreName( "oss");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        client.getFile(req, 10000);
    }

    // normal case
    @Test
    public void testDelFile1() throws Exception {

        DelFileRequest req = new DelFileRequest();
        req.setStoreName( "oss");
        req.setFileName( "test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

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
        req.setFileName( "test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        client.delFile(req, 10000);
    }

    // miss file name
    @Test(expected = IllegalArgumentException.class)
    public void testDelFile4() throws Exception {

        DelFileRequest req = new DelFileRequest();
        req.setStoreName( "oss");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        client.delFile(req, 10000);
    }

    // normal
    @Test
    public void testListFile1() throws Exception {

        ListFileRequest req = new ListFileRequest();
        req.setStoreName( "oss");
        req.setName("dir");
        req.setMarker("test.log");
        req.setPageSize(10);

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        ListFileResponse resp = client.listFile(req, 10000);

        Assert.assertTrue(resp.isTruncated());
        Assert.assertEquals("marker", resp.getMarker());
        Assert.assertEquals(1, resp.getFiles().length);

        FileInfo f = resp.getFiles()[0];
        Assert.assertEquals("put file store name oss, meta 2", f.getFileName());
        Assert.assertEquals(100L, f.getSize());
        Assert.assertEquals("2021-11-23 10:24:11", f.getLastModified());
        Assert.assertEquals("v1", f.getMetaData().get("k1"));
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
        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        ListFileResponse resp = client.listFile(req, 10000);
    }

    // normal
    @Test
    public void testGetFileMeta1() throws Exception {

        GetMetaRequest req = new GetMetaRequest();
        req.setStoreName( "oss");
        req.setFileName( "test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        GetMeteResponse res = client.getFileMeta(req, 10000);

        Assert.assertEquals(100L,res.getSize());
        Assert.assertEquals("2021-11-22 10:24:11",res.getLastModified());
        Assert.assertArrayEquals(new String[]{"v1","v2"},res.getMeta().get("k1"));
    }

    // miss request
    @Test(expected = IllegalArgumentException.class)
    public void testGetFileMeta2() throws Exception {
        client.getFileMeta(null, 10000);
    }

    // miss store name
    @Test(expected = IllegalArgumentException.class)
    public void testGetFileMeta3() throws Exception {

        GetMetaRequest req = new GetMetaRequest();
        req.setFileName( "test.log");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        GetMeteResponse res = client.getFileMeta(req, 10000);

        Assert.assertEquals(100L,res.getSize());
        Assert.assertEquals("2021-11-22 10:24:11",res.getLastModified());
        Assert.assertArrayEquals(new String[]{"v1","v2"},res.getMeta().get("k1"));
    }

    // miss file name
    @Test(expected = IllegalArgumentException.class)
    public void testGetFileMeta4() throws Exception {

        GetMetaRequest req = new GetMetaRequest();
        req.setStoreName( "oss");

        Map<String,String> metaData = new HashMap<>();
        metaData.put("k1", "v1");
        metaData.put("k2", "v2");
        req.setMetaData(metaData);

        client.getFileMeta(req, 10000);
    }
}
