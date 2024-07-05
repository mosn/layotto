# Automate testing of Quickstart documentation with tools

Quickstart is the face of the project, so if a new user enters the repository and finds that the Quickstart documentation doesn't work, they may be disappointed and walk away.

So we need to test Quickstart regularly to make sure it works.

But the process of manually testing Quickstart periodically and fixing exceptions in the documentation is too time-consuming.

<img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*fTI5RbfAK3gAAAAAAAAAAAAAARQnAQ" width="30%" height="30%">

It's a pain in the ass!

Let's use the tool to test the documentation automatically!

## Principle

Use the tool to execute all shell scripts in a markdown file sequentially, i.e. all scripts wrapped in:

~~~markdown
```shell
```
~~~

Note: The script wrapped in `bash` blocks will NOT be run.

~~~markdown
```bash
```
~~~

## step 1. Install `mdx`
see https://github.com/seeflood/mdx#installation

## step 2. Close local software that may cause conflicts
Close the local Layotto, to avoid port conflicts when running the document.

Similarly, if the documentation will start containers like Redis with Docker, you need to shut down and remove containers that may cause port conflicts and container name conflicts first.

## step 3. Running documentation

As an example, run the Quickstart documentation for the state API:

```shell
mdx docs/en/start/state/start.md 
```

## step 4. Reported an error? Test-driven development, optimize your documentation!
You can think of each document as a UT, which should have 4 phases: prepare, execute, verify, and release resources.

If the document runs with an error, it means that the case needs to be optimized.

This is also the idea of "test-driven development", optimizing the documentation to make it "testable", right?

For example, I ran the Quickstart documentation for the state API and found an error:

```bash
SaveState succeeded.key:key1 , value: hello world 
GetState succeeded.[key:key1 etag:1]: hello world
SaveBulkState succeeded.[key:key1 etag:2]: hello world
SaveBulkState succeeded.[key:key2 etag:2]: hello world
GetBulkState succeeded.key:key2 ,value:hello world ,etag:1 ,metadata:map[] 
GetBulkState succeeded.key:key1 ,value:hello world ,etag:2 ,metadata:map[] 
GetBulkState succeeded.key:key3 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key4 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key5 ,value: ,etag: ,metadata:map[] 
panic: error deleting state: rpc error: code = Aborted desc = failed deleting state with key key1: possible etag mismatch. error from state store: ERR Error running script (call to f_9b5da7354cb61e2ca9faff50f6c43b81c73c0b94): @user_script:1: user_script:1: failed to delete key1

goroutine 1 [running]:
main.testDelete(0x16bc760, 0xc0000ac000, 0x16c56a0, 0xc0000b90e0, 0x15f30e1, 0x5, 0x15f2539, 0x4)
        /Users/qunli/projects/layotto/demo/state/redis/client.go:73 +0x13d
main.main()
        /Users/qunli/projects/layotto/demo/state/redis/client.go:57 +0x2f4
exit status 2
```

After some troubleshooting, we found that the demo client did not pass the `etag` field when deleting the specified key, which caused the demo to run abnormally.

See, through the automated testing documentation, we found a Quickstart bug :)

### Q: How to write "testable" documents
Note: You can refer to the sample documentation that runs the test: docs/en/start/state/start.md

Some common details are explained below.

#### The demo code should panic when it doesn't meet expectations
For example, in the demo `demo/state/redis/client.go`, if you get an error when calling Layotto, you should just panic:

```go
if err := cli.SaveBulkState(ctx, store, item, &item2); err != nil {
	panic(err)
}
```

In addition to judging errors, the demo should also verify the test results, and panic directly if it does not meet expectations. This is equivalent to UT, after calling a method, the result of the call needs to be verified.

The advantage of this is that once the Quickstart does not meet expectations, the demo will exit abnormally, allowing automated tools to find "the test failed! Find someone to fix it!"

#### It is best to delete the container and free resources at the end of the document

When writing UT, we do things like release resources, restore mocks, etc. in the final stages; to make the document "testable", do something similar.

For example to delete the redis Docker container at the end of the document:

```shell
docker rm -f redis-test
```

Note: Layotto's github workflow deletes all containers and closes applications such as layotto, etcd after each md is executed.
So even if the container is not deleted in the document, it will not affect the github workflow to run the test.

#### What should I do if I don't want a certain command to be executed?
`mdx` by default will only execute shell code blocks, i.e. code blocks written like this:

```shell
```shell
```

If you don't want a block of code to be executed, you can change the shell to something else, for example:

```bash
```bash
```

#### A certain shell command will hang and affect the test, what should I do?
Again, take docs/en/start/state/start.md as an example.

One of the scripts will run Layotto, but if you run it it will hang, preventing the test tool from continuing to run the next command:

```bash
./layotto start -c ../../configs/config_redis.json
```

How to do it?

##### Solution 1:

Annotated with @background, see https://github.com/seeflood/mdx#background

~~~
```shell @background
./layotto start -c ../../configs/config_standalone.json
```
~~~

##### Solution 2:

Don't run this script, add a "hidden script" that "runs Layotto in the background", this hidden script is wrapped in a comment, so it won't be seen by people reading the documentation, but `mdx` will still run it:

```bash
    ```bash
    ./layotto start -c ../../configs/config_redis.json
    ```
    
    <!-- The command below will be run when testing this file 
    ```shell
    nohup ./layotto start -c ../../configs/config_redis.json &
    ```
    -->
```

#### How to deal with the command to switch directories?

We can assume that the current directory is the project's root path.


Then the switch path can be written like this:

```bash
# change directory to ${your project path}/demo/state/redis/
 cd demo/state/redis/
 go run .
```

What if you want to go back to the root path after running this command?

##### Solution 1:


Use the `${project_path}` variable to represent the project root path, see https://github.com/seeflood/mdx#cd-project_path

```shell 
cd ${project_path}/demo/state/redis/
```

##### Solution 2:

Add a hidden script to switch directories. For example, write:

    <!-- The command below will be run when testing this file 
    ```shell
    cd ../../
    # if we should wait for layotto to start, we can:
    # sleep 1s 
    ```
    -->
    
    ```shell
    # open a new terminal tab
    # change directory to ${your project path}/demo/state/redis/
    cd demo/state/redis/
    go run .
    ```

### Other markdown annotations

The mdx tool provides many "markdown annotations" to help you write "runnable markdown files". If you are interested, you can check the [mdx documentation](https://github.com/seeflood/mdx#usage)

### Fix the error and see the effect!

After a fix, I ran the document again:

```shell
mdx docs/en/start/state/start.md
```


The document does not report an error, it can run normally and exit:

```bash
admindeMacBook-Pro-2:layotto qunli$ mdx docs/en/start/state/start.md
latest: Pulling from library/redis
Digest: sha256:69a3ab2516b560690e37197b71bc61ba245aafe4525ebdece1d8a0bc5669e3e2
Status: Image is up to date for redis:latest
docker.io/library/redis:latest
REPOSITORY                     TAG         IMAGE ID       CREATED         SIZE
redis                          latest      bba24acba395   3 days ago      113MB
pseudomuto/protoc-gen-doc      latest      35472df9ecbb   6 weeks ago     39.5MB
apache/skywalking-oap-server   8.0.1-es7   887769fd3bf6   21 months ago   191MB
apache/skywalking-ui           8.0.1       42b3b496616e   21 months ago   127MB
5835d4652c057ce7ea109564c3e36351335ec53c3dedb02650f2056d3cc3edd5
appending output to nohup.out
runtime client initializing for: 127.0.0.1:34904
SaveState succeeded.key:key1 , value: hello world 
GetState succeeded.[key:key1 etag:1]: hello world
SaveBulkState succeeded.[key:key1 etag:2]: hello world
SaveBulkState succeeded.[key:key2 etag:2]: hello world
GetBulkState succeeded.key:key1 ,value:hello world ,etag:2 ,metadata:map[] 
GetBulkState succeeded.key:key4 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key3 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key5 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key2 ,value:hello world ,etag:1 ,metadata:map[] 
DeleteState succeeded.key:key1
DeleteState succeeded.key:key2
redis-test
```

## step 5. Modify CI to automatically test the newly written quickstart document

If you have written a new quickstart document, and the self-test can run normally, the next step is to modify the CI to achieve "every time someone submits a pull request, the tool automatically tests that this quickstart document can run through".

The modification method is:

1. Modify the script `etc/script/test-quickstart.sh` to add your documentation to it:

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*ZPRlRa7a-0QAAAAAAAAAAAAAARQnAQ)

2. If you need to automatically release some resources before and after the document runs (such as automatically killing the process, deleting the docker container), you can add the resources to be released in the script. For example, if you want to implement "automatically kill the etcd process every time a document is run", you can add to the script:

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*0th0Q7yn5MIAAAAAAAAAAAAAARQnAQ)


3. After making the above changes, it is time to test the new CI.

Run in the project root directory

```shell
make style.quickstart
```

These documents will be tested:

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*I7LRSryXwWYAAAAAAAAAAAAAARQnAQ)


> [!TIP|label: run locally with caution, this script will delete some docker containers]
> This command will delete the Docker containers that contain the keywords in the image. If you don't want to delete these containers, don't run them locally:
> ![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*N3CIRb0883kAAAAAAAAAAAAAARQnAQ)


whereas if you run:

```shell
make style.quickstart QUICKSTART_VERSION=1.17
```

The following documents will be tested (these documents can only run successfully in golang 1.17 and above):

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*X3F9QJSKq3QAAAAAAAAAAAAAARQnAQ)
