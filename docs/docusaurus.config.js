
import {themes as prismThemes} from 'prism-react-renderer';

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Layotto',
  tagline: 'Layotto (L8): To be the next layer of OSI layer 7',
  favicon: 'https://gw.alipayobjects.com/zos/bmw-prod/65518bfc-8ba5-4234-a5c5-2bc065e3a5f0.svg',

  url: 'https://layotto.github.io',
  baseUrl: '/layotto/',

  organizationName: 'mosn',
  projectName: 'layotto',

  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  i18n: {
    defaultLocale: 'zh-Hans',
    locales: ['en-US','zh-Hans'],
    localeConfigs: {
      "en-US": {
        label:"English"
      }
    }
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl:({  docPath, locale }) => {
            //把docPath 拆分，中间加上对应的路径。
            let newDocPath;

            if (locale !== 'en-US') {
              const pathSegments = docPath.split('/');
              newDocPath = ['docs', ...pathSegments].join('/');
              return `https://github.com/mosn/layotto/edit/main/docs/`+newDocPath;
            }else{
              const pathSegments = docPath.split('/');
              newDocPath = ['i18n/en-US/docusaurus-plugin-content-docs/current', ...pathSegments].join('/');
              return `https://github.com/mosn/layotto/edit/main/docs/`+newDocPath;
            }

          },
        },
        blog: {
          blogSidebarTitle: '全部博客',
          blogSidebarCount: 'ALL',
          showReadingTime: true,

          editUrl:({  locale,blogDirPath, blogPath }) => {
            return `https://github.com/mosn/layotto/edit/main/docs/${blogDirPath}/${blogPath}`;
          }
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      },
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      docs: {
        sidebar: {
          hideable: true,
          autoCollapseCategories: true,
        },
      },
      // Replace with your project's social card
      image: 'img/docusaurus-social-card.jpg',
      navbar: {
        title: '',
        logo: {
          alt: 'Layotto Logo',
          src: 'https://gw.alipayobjects.com/zos/bmw-prod/65518bfc-8ba5-4234-a5c5-2bc065e3a5f0.svg',
        },
        items: [
          {
            type: 'doc',
            docId: 'README',
            position: 'left',
            label: '文档',
          },
          {
            type: 'localeDropdown',
            position: 'right',
          },
          {to: '/blog', label: '博客', position: 'left'},
          {
            href: 'https://github.com/mosn/layotto',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },

      prism: {
        theme: prismThemes.github,
        darkTheme: prismThemes.dracula,
      },
      algolia: {
          appId: 'B0I4Q5CLN8',
          apiKey: '79f410a7e620927c50c3ae6c8c9af5bb',
          indexName: 'layotto',
          contextualSearch: true,
      },
    }),
};

export default config;


