# Document Contribution Guide

Thank you for your support in Layotto!

This document describes how to modify/add documents. Documentation for this repository is written in Markdown.

## 1. Document Path description

Documents are stored in the 'docs/' directory, where 'docs/en' stores English documents and 'docs/zh' stores Chinese documents.

![img_2.png](../../img/development/doc/img_2.png)

## 2. Documentation Site Description
Files under docs/ directory will be automatically deployed to github pages and rendered through [docsify](https://docsify.js.org/#/).

Generally speaking, after the .md file is merged into the main branch, you can see the new page on Layotto's documentation site, and all deployment and rendering processes are done automatically.

## 3. How to Add a New Document
### step 1. Write a new markdown file
To add a document, create a folder and a .md file based on the directory structure. For example, if you want to write a design document for the distributed lock API, just create a new directory:

![img_1.png](../../img/development/doc/img_1.png)

### step 2. Update the sidebar
Remember to update the sidebar after adding new documents or revising existing documents.

Chinese sidebar: 'docs/zh/_sidebar.md'

English sidebar: 'docs/_sidebar.md'

### step 3. Submit a Pull request
After writing the above Markdown files, submitting pr, and merging it into the main branch, new documents are now available on the official website.

## 4. Tips on Hyperlinks

One annoying problem with Docsify is that the use of hyperlinks is confusing.

The hyperlink mentioned here is the kind of links that will jump to other documents once clicked, such as the following:

![image](https://user-images.githubusercontent.com/26001097/132220354-db2b6ad0-58e4-46ed-b005-71d8134f725b.png)

### Incorrect Syntax
If you try to create a hyperlink with a relative path, then a 404 page will appear once you clicked it:

![img_6.png](../../img/development/doc/img_6.png)

![img_7.png](../../img/development/doc/img_7.png)

### Correct Syntax

There are two suggested ways to write hyperlinks:

a. Use a path relative to the 'docs/' directory. Such as:

![img_5.png](../../img/development/doc/img_5.png)

b. Use the full Url. Such as:

```markdown
see [runtime_config.json](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json):
```

## 5. Tips on image links
Images are stored under docs/img/ directory for the purpose that the Docsify site can access it

![img.png](../../img/development/doc/img.png)

It is recommended to use the full path when referencing images in documents, to avoid a bunch of messy path problems.

For example, when referencing the images under the main branch, the prefix of the image url is `raw.githubusercontent.com/mosn/layotto/main/docs/img/xxx`

and the Markdown phrase referring to an image will be ：

```markdown
![Architecture](https://raw.githubusercontent.com/mosn/layotto/main/docs/img/runtime-architecture.png)
```

Note: Relative paths can also be used, but you may encounter many problems. For example, the relative path logic of the `<img>` tag and `![xxx](url)` tag are different; for example, users may access the README through different paths, so it's hard for you to define the relative path. To avoid these problems, it's recommended to use a full url.
