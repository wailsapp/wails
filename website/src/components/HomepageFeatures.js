import React from 'react';
import clsx from 'clsx';
import styles from './HomepageFeatures.module.css';

import ThemedImage from '@theme/ThemedImage';
import useBaseUrl from '@docusaurus/useBaseUrl';
import {translate} from '@docusaurus/Translate'; // i18n component

const FeatureList = [
  {
    title: translate({id: "homepage.Features.Title1"}),
    lightImage: '/img/desktop-application-app-dark.svg',
    darkImage: '/img/desktop-application-app.svg',
    description: translate({id: "homepage.Features.Description1"}),
  },
  {
    title: translate({id: "homepage.Features.Title2"}),
    lightImage: '/img/go_js.svg',
    darkImage: '/img/go_js_dark.svg',
    description: translate({id: "homepage.Features.Description2"}),
  },
  {
    title: translate({id: "homepage.Features.Title3"}),
    lightImage: '/img/Terminal-icon.svg',
    darkImage: '/img/Terminal-icon-dark.svg',
    description: translate({id: "homepage.Features.Description3"}),
  },
];

function Feature({lightImage, darkImage, title, description}) {
  const imageStyle = {
    width: '25%',
    margin: 'auto',
    minHeight: '100px',
  }

  return (
      <div className={clsx('col col--4')}>
        <div className="text--center">
          <ThemedImage
              style={imageStyle}
              alt="Docusaurus themed image"
              sources={{
                light: useBaseUrl(lightImage),
                dark: useBaseUrl(darkImage),
              }}
          />
        </div>
        <div className="text--center padding-horiz--md">
          <h3>{title}</h3>
          <p>{description}</p>
        </div>
      </div>
  );
}

export default function HomepageFeatures() {
  return (
      <section className={styles.features}>
        <div className="container">
          <div className="row">
            {FeatureList.map((props, idx) => (
                <Feature key={idx} {...props} />
            ))}
          </div>
        </div>
      </section>
  );
}
