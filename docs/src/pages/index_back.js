import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';
import styles from './index.module.css';
//
// function HomepageHeader() {
//   const {siteConfig} = useDocusaurusContext();
//   return (
//     <header className={clsx('hero hero--primary', styles.heroBanner)}>
//       <div className="container">
//         <Heading as="h1" className="hero__title">
//           {siteConfig.title}
//         </Heading>
//         <p className="hero__subtitle">{siteConfig.tagline}</p>
//         <div className={styles.buttons}>
//           <Link
//             className="button button--secondary button--lg"
//             to="/docs/README.md">
//             Layotto Tutorial - 5min ⏱️
//           </Link>
//         </div>
//       </div>
//     </header>
//   );
// }
//
// export default function Home() {
//   const {siteConfig} = useDocusaurusContext();
//   return (
//     <Layout
//       title={`Hello from ${siteConfig.title}`}
//       description="Description will go into a meta tag in <head />">
//       <HomepageHeader />
//       <main>
//         <HomepageFeatures />
//       </main>
//     </Layout>
//   );
// }


import { useEffect } from 'react';
import { useHistory } from 'react-router-dom';

function RedirectToDocs() {
    const history = useHistory();
    const {siteConfig} = useDocusaurusContext();

    useEffect(() => {
        console.log("本地化:"+siteConfig.i18n.locales)
        history.replace('/layotto-docusaurus/');
    }, []);

    return null;
}

export default RedirectToDocs;

