import React from 'react';
import clsx from 'clsx';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import styles from './index.module.css';
import HomepageFeatures from '../components/HomepageFeatures';
import {translate} from '@docusaurus/Translate'; // i18n component
var Carousel = require('react-responsive-carousel').Carousel;

function HomepageHeader() {
    const {siteConfig} = useDocusaurusContext();
    return (
        <header className={clsx('hero', styles.heroBanner)}>
            <div className="container">
                <Carousel showArrows={false} width={"50%"}
                          showThumbs={false} stopOnHover={false}
                          showStatus={false} autoPlay={true}
                          showIndicators={false}
                          infiniteLoop={true} interval={3000} transitionTime={1000}>
                    <div>
                        <img src="/img/logo-dark.svg"/>
                    </div>
                    <div>
                        <img src="/img/showcase/ytd.png"/>
                    </div>
                    <div>
                        <img src="/img/showcase/wombat.png"/>
                    </div>
                    <div>
                        <img src="/img/showcase/wally.png"/>
                    </div>
                    <div>
                        <img src="/img/showcase/mollywallet.png"/>
                    </div>

                </Carousel>
                {/*<ThemedImage*/}
                {/*    alt="Wails Logo"*/}
                {/*    width="30%"*/}
                {/*    sources={{*/}
                {/*        light: useBaseUrl('/img/logo-light.svg'),*/}
                {/*        dark: useBaseUrl('/img/logo-dark.svg'),*/}
                {/*    }}*/}
                {/*/>*/}
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
