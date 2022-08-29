import React from "react";
import clsx from "clsx";
import Layout from "@theme/Layout";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import styles from "./index.module.css";
import HomepageFeatures from "../components/HomepageFeatures";
import {translate} from "@docusaurus/Translate"; // i18n component
import useBaseUrl from "@docusaurus/useBaseUrl";

var Carousel = require("react-responsive-carousel").Carousel;

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
      <header className={clsx("hero", styles.heroBanner)}>
        <div className="container">
          <Carousel
          showArrows={false}
          width={"100%"}
          showThumbs={false}
          stopOnHover={false}
          showStatus={false}
          autoPlay={true}
          showIndicators={false}
          infiniteLoop={true}
          interval={4000}
          transitionTime={1000}
        >
          <div className="slide-item-box">
            <img src={useBaseUrl("/img/showcase/mac-app.png")}/>
          </div>
          <div className="slide-item-box">
            <img src={useBaseUrl("/img/showcase/ytd.webp")}/>
          </div>
          <div className="slide-item-box">
            <img src={useBaseUrl("/img/showcase/wombat.webp")}/>
          </div>
          <div className="slide-item-box">
            <img src={useBaseUrl("/img/showcase/restic-browser-2.png")} />
          </div>
          <div className="slide-item-box">
            <img src={useBaseUrl("/img/showcase/wally.webp")}/>
          </div>
          <div className="slide-item-box">
            <img src={useBaseUrl("/img/showcase/october.webp")}/>
          </div>
          <div className="slide-item-box">
            <img
                className="img"
                src={useBaseUrl("/img/showcase/mollywallet.webp")}
            />
          </div>
        </Carousel>

        <p className="hero__subtitle">
          {translate({id: "homepage.Tagline"})}
        </p>
          <div className={styles.buttons}>
            <Link
                className="button button--secondary button--outline button--lg"
                to="/docs/introduction"
            >
              {translate({id: "homepage.LearnMoreButtonText"})}
            </Link>
            <Link
                className="button button--primary button--lg"
                to="/docs/gettingstarted/installation"
            >
              {translate({id: "homepage.ButtonText"})}
            </Link>
          </div>
      </div>
    </header>
  );
}

export default function Home() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`The ${siteConfig.title} Project`}
      description={translate({ id: "homepage.Tagline" })}
    >
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
