# Layotto GitHub Workflows

This document explains Layotto's four workflows in Github:

+ Layotto Env Pipeline 🌊
+ Layotto Dev Pipeline 🌊 (Before Merged)
+ Layotto Dev Pipeline 🌊 (After Merged)
+ Layotto Release Pipeline 🌊

The workflow contains one or more tasks, It improves the standardization and security of the code in layotto, simplifies repetitive steps of development / build / release. The following is a detailed explanation of the above four workflows.

## Layotto Env Pipeline 🌊

### Job Task Content

Layotto Env Pipeline is mainly responsible for the project of layotto and the specification of relevant environment,it current contains the following tasks：

+ Title Validation (Check the specification of PR title based on semantic style)
+ Quickstart Validation (Verification of QuickStart documents)
+ Update Stale Status (Update of issue / PR status)
+ License Validation (Verification of license)
+ DeadLink Validation (Check the deadLink in document)
+ CodeQL (Analysis of CodeQL)

### Job Trigger Event

Layotto Env Pipeline Task Trigger Events:

+ Title Validation: 
  
  ```
  pull_request:
      types:
      - opened open PR 
      - edited edit PR
      - synchronize synchronize PR
      - labeled PR add Label
      - unlabeled PR cancel Label
  ```

+ Quickstart Validation: 
  
  ```
  push:
      branches:
      - main merge PR
  pull_request:
      branches:
      - main commit PR
  ```

+ Update Stale Status: 
  
  ```
  on:
  schedule:
      - cron: '30 1 * * *' timed tasks
  ```

+ License Validation: 
  
  ```
  push:
      branches:
      - main merge PR
  pull_request:
      branches:
      - main commit PR
  ```

+ DeadLink Validation: 
  
  ```
  pull_request:
      branches:
      - main commit PR
  ```

+ CodeQL: 
  
  ```
  schedule:
      - cron: '0 4 * * 5' timed tasks
  ```

## Layotto Dev Pipeline 🌊 (Before Merged)

![release.png](../../img/development/workflow/workflow-dev.png)

### Job Task Content

The layotto dev pipeline (before merged)  is mainly responsible for verifying the code after submitting the PR, which currently includes the following tasks:

+ Go Style Check : Check the style of the code
+ Go CI Linter : Perform linter specification of verification on the code
+ Go Unit Test : Unit test the code
+ Coverage Analysis : Coverage analysis of the code
+ Integrate with WASM : WASM integration test on the code
+ Integrate with Runtime : Run time integration test on the code
+ Darwin AMD64 Artifact : Build Darwin AMD64 binary verification for code
+ Darwin ARM64 Artifact : Build Darwin arm64 binary verification for code
+ Linux AMD64 Artifact : Build linux amd64 binary verification for code
+ Linux ARM64 Artifact : Build linux arm64 binary verification for code
+ Linux AMD64 WASM Artifact : Build linux AMD64 binary verification for layotto wasm

### Job Trigger Event

```
    on:
    push:
        branches: [main] merge PR
        paths-ignore: ignore the following changes: docs directory files，markdown files
        - 'docs/**'
        - '**/*.md'
    pull_request:
        branches: "*" merge PR
        paths-ignore: ignore the following changes: docs directory files，markdown files
        - 'docs/**'
        - '**/*.md'
```

## Layotto Dev Pipeline 🌊 (After Merged)

![release.png](../../img/development/workflow/workflow-merge.png)

### Job Task Content

The layotto dev pipeline (after merged)  is mainly responsible for the verification and release of the combined layotto code, which currently includes the following tasks：

+ Go Style Check : Check the style of the code
+ Go CI Linter : Perform linter specification of verification on the code
+ Go Unit Test : Unit test the code
+ Coverage Analysis : Coverage analysis of the code
+ Integrate with WASM : WASM integration test on the code
+ Integrate with Runtime : Run time integration test on the code
+ Darwin AMD64 Artifact : Build Darwin AMD64 binary verification for code
+ Darwin ARM64 Artifact : Build Darwin arm64 binary verification for code
+ Linux AMD64 Artifact : Build linux amd64 binary verification for code
+ Linux ARM64 Artifact : Build linux arm64 binary verification for code
+ Linux AMD64 WASM Artifact : Build linux AMD64 binary verification for layotto wasm
+ Linux AMD64 WASM Image : Release the latest version of layotto wasm image. The image specification is layotto/faas-amd64:latest
+ Linux AMD64 Image : Release the latest version of layotto wasm image. The image specification is layotto/layotto:latest
+ Linux ARMD64 Image : Release the latest version of layotto wasm image. The image specification is layotto/layotto.arm64:latest

### Job Trigger Event

```
    on:
    push:
        branches: [main] merge PR
        paths-ignore: ignore the following changes： docs directory files，markdown files
        - 'docs/**'
        - '**/*.md'
    pull_request:
        branches: "*" create a PR
        paths-ignore: ignore the following changes： docs directory files，markdown files
        - 'docs/**'
        - '**/*.md'
```

## Layotto Release Pipeline 🌊

![release.png](../../img/development/workflow/workflow-release.png)

### Job Task Content

The layotto release pipeline  is mainly responsible for the release and verification of the new version of layotto, which currently includes the following tasks :

+ Go Style Check : Check the style of the code
+ Go CI Linter : Perform linter specification of verification on the code
+ Go Unit Test : Unit test the code
+ Coverage Analysis : Coverage analysis of the code
+ Integrate with WASM : WASM integration test on the code
+ Integrate with Runtime : Run time integration test on the code
+ Darwin AMD64 Artifact : Build Darwin AMD64 binary verification for code
+ Darwin ARM64 Artifact : Build Darwin arm64 binary verification for code
+ Linux AMD64 Artifact : Build linux amd64 binary verification for code
+ Linux ARM64 Artifact : Build linux arm64 binary verification for code
+ Linux AMD64 WASM Artifact : Build linux AMD64 binary verification for layotto wasm
+ Linux AMD64 WASM Image : Release the latest version of layotto wasm image. The image specification is layotto/faas-amd64:{latest_tagname}
+ Linux AMD64 Image : Release the latest version of layotto wasm image. The image specification is layotto/layotto:{latest_tagname}
+ Linux ARMD64 Image : Release the latest version of layotto wasm image. The image specification is layotto/layotto.arm64:{latest_tagname}

### Job Trigger Event

```
    on:
    create  Tag or Branch,combined with the following conditions

    if: ${{ startsWith(github.ref, 'refs/tags/') }} changes to Tag(Ignore creation of new branch)
```

> The configuration file of Layotto's GitHub workflow is in [here](https://github.com/mosn/layotto/tree/main/.github/workflows)