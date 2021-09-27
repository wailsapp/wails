import React from 'react';
import clsx from 'clsx';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import styles from './index.module.css';
import HomepageFeatures from '../components/HomepageFeatures';

import ThemedImage from '@theme/ThemedImage';
import useBaseUrl from "@docusaurus/core/lib/client/exports/useBaseUrl";
import {translate} from '@docusaurus/Translate'; // i18n component

function HomepageHeader() {
    const {siteConfig} = useDocusaurusContext();
    return (
        <header className={clsx('hero', styles.heroBanner)}>
            <div className="container">
                <ThemedImage
                    alt="Wails Logo"
                    width="30%"
                    sources={{
                        light: useBaseUrl('/img/logo-light.svg'),
                        dark: useBaseUrl('/img/logo-dark.svg'),
                    }}
                />
                <p className="hero__subtitle">{translate({id: "homepage.Tagline"})}</p>
                <div className={styles.buttons}>
                    <Link
                        className="button button--primary button--lg"
                        to="/docs/gettingstarted/installation">
                        {translate({id: "homepage.ButtonText"})}
                    </Link>
                </div>
            </div>
        </header>
    );
}

export default function Home() {
    const {siteConfig} = useDocusaurusContext();
    return (
        <Layout
            title={`The ${siteConfig.title} Project`}
            description={translate({id: "homepage.Tagline"})}>
            <HomepageHeader/>
            <main>
                <HomepageFeatures/>
            </main>
        </Layout>
    );
}
