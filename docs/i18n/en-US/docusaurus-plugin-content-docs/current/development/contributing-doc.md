# Docs contributing guide

Thank you for your support of Layotto!

This document describes how to modify/add documents.The documents in this repository are written using Markdown syntax.

## 1. Document Path Description

Documents are uniformly placed in the docs/ directory, with docs/docs storing documents, docs/blog storing blogs, and docs/i18n storing translated multilingual documents.

![img\_14.png](/img/development/doc/img_14.png)

## 2. Documentation Site Description

Files under the docs/ directory will be automatically deployed to github pages and rendered through [docusaurus](https://docusaurus.io/).

Generally speaking, after the .md file is merged into the main branch, you will be able to see the new page on the Layotto documentation site, and the deployment and rendering processes are automatic.

### Local startup documentation site

After writing the document locally, in order to quickly preview the effect, you can also refer to [docusaurus_installation](https://docusaurus.io/docs/installation) to start the documentation site locally.

Here is a summary of the steps:

step 1. Install Docusaurus, make sure the Node.js version is 18.0 or above

```shell
npm install
```

step 2. Compile, compile the documents under 'docs' into static HTML files, and place the compiled files under 'docs/build'

```shell
npm run build --config docs
```

Step 3. Start the documentation site

```shell
# Run npm run serve --config docs in the layotto project root directory
npm run serve --config docs
```

step 3. Open http://localhost:3000/ to view the documentation site.

## 3. What needs to be done to add a document

### step 1. Create a markdown document

When you need to add a new document, you can create a new folder according to the directory structure, and create a .md file.For example, if you want to write the design document for distributed lock API, create directories under both Chinese and English directories:

![img\_8.png](/img/development/doc/img_8.png)
![img\_9.png](/img/development/doc/img_9.png)

You can use Crowdin for auxiliary translation. After writing the Chinese documentation, use npm run crowdin:upload to upload the Chinese documentation to the Crowdin platform. Find the Layotto project on the platform, translate it, and download it to the local response directory.[crowdin reference documentation](https://docusaurus.io/docs/i18n/crowdin)

### step 2. Add the document to the sidebar

After adding new documents and finishing the content, remember to update the sidebar.

The Chinese sidebar is in docs/sidebars.js

The English sidebar needs to perform the following:

```shell
npm run write-translations -- --locale en-US
```

Then modify the corresponding sidebar content in docs/i18n/en-US/docusaurus-plugin-content-docs/current.json

![img\_10.png](/img/development/doc/img_10.png)

### step 3. (optional) Start the local documentation site, validate

You can use docusaurus to start a local documentation server and view the editing effects

### Step 4. Submit PR, merge into the code repository

After finishing the above markdown file, submitting a PR, merging into the master branch, opening the official website will show the new document.

## 4. Common pitfall: hyperlinks in documents

There is an annoying issue with using Docusaurus to build a website: hyperlinks look very strange.

The hyperlinks mentioned here are links that, when clicked, will redirect to other documents, such as the following:

![img\_4.png](/img/development/doc/img_4.png)

### 4.1. Incorrect Writing

If you try to use a relative path to write a hyperlink URL, you will find that clicking on it in the website will lead to a 404 error:

![img\_6.png](/img/development/doc/img_6.png)

![img\_7.png](/img/development/doc/img_7.png)

### 4.2. Correct Writing

There are two correct ways to use hyperlinks:

a. Use absolute path relative to the docs/ directory.For example:

![img\_11.png](/img/development/doc/img_11.png)

Write in this way, warning in local compilers will be alerted, not ignored, and automatically processed into available hyperlinks when compiled.

b. Use completed Url。For example:

```markdown
see [runtime_config.json](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json):
```

## 5. Image directory to image link

Image is placed in the docs/static/img/ directory.This is in order to allow the docusaurus site to reach and docusaurus compiles the folders below static's top：

![img12.png](/img/development/doc/img_12.png)

Local images with statics in the document recommend an absolute path /img if they are a network image, and a full image URL is sufficient.

For example, in the case of images with the mains branch, the prefix for the image Url is `raw.githubusercontent.com/mosn/layotto/main/docs/img/xxx`

The Markdown syntax is as follows:

```markdown
![Architecture](https://raw.githubusercontent.com/mosn/layotto/main/docs/img/runtime-architecture.png)
```
