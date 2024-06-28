# Publish Handbook

What should be done by publishing a new version under this section

## Publish Period

The Layotto publication cycle is tentatively scheduled to be issued on a quarterly basis.

## Publish checklist

### Step1: Check [Roadmap]for current iterations (https://github.com/mosn/layotto/projects)

1. Check ongoing tasks
2. Check unfinished tasks
3. Confirm Task Status and Publish Content with Owner

### Step2: Create Release tag, push to github, and check workflow

1. 规范：请按照 `v{majorVersion}.{subVersion}.{latestVersion}` 格式创建 tag。

2. Waiting for CI to end, confirm the following：
   - CI Test Jobs through：
     - Go code style verification
     - Go code specification validation
     - Go Unit Test
     - Go Integration Test
   - CI Multiplatform Artifacts build Jobs through：
     - Linux/AMD64 Artifacts successfully built and uploaded
     - Linux/ARM64 Artifacts successfully built and uploaded
     - Darwin/AMD64 Artifacts successfully built and uploaded
     - Darwin/ARM64 Artifacts successfully built and uploaded
   - CI Multiplatform Image build/publish Jobs through：
     - Linux/AMD64 Image successfully built and Push DockerHub
     - Linux/ARM64 Image successfully built and Push DockerHub
       - Image Tag specification：
         - AMD64/X86 架构的镜像：`layotto/layotto:{tag}`
         - ARM64 架构的镜像：`layotto/layotto.arm64:{tag}`

![release.png](/img/development/workflow/release.png)

### Step3: Draft a new release and prepare a release report

> Publish reports can be automatically generated using the github's functionality before making changes based on the generated content.

> Reference can be made to previous [发版报告](https://github.com/mosn/layotto/releases)

![img\_1.png](/img/development/release/img_1.png)

### Step4: Upload Binars of Multi-Platform Architecture

> Update：by 2022/05/04 is negligible.Release Pipeline from Layotto will automatically upload binary files without having to upload manually.PR See https://github.com/mosn/layotto/pull/566

> If you do not upload automatically, you can manually download and upload the multi-platform Artifacts built in `Step 2`

![img.png](/img/development/release/img.png)

### Step5: Confirm Publish

1. Click to publish
2. General information
3. Check [Roadmap](https://github.com/mosn/layotto/projects), modify unfinished tasks in previous version, change milestone to next version
4. If there is a SDK release, it needs to be published in the SDK repository and upload the central repository (e.g. Java SDK needs to be uploaded to Maven central repository).
