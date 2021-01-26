/* eslint-disable */

const path = require('path');

const platform = process.env.WAILSPLATFORM;
if (!platform) {
  console.error("FATAL: Environment variable WAILSPLATFORM not set!");
  process.exit(1);
}

module.exports = {
  entry: './core/desktop',
  mode: 'production',
  output: {
    path: path.resolve(__dirname, '..', 'assets'),
    filename: 'desktop_'+platform+'.js',
    library: 'Wails'
  },
  resolve: {
    alias: {
      ipc$: path.resolve(__dirname, 'desktop/ipc.js'),
      platform$: path.resolve(__dirname, `desktop/${platform}.js`)
    }
  },
  module: {
    rules: [
      {
        test: /\.m?js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: 'babel-loader',
          options: {
            plugins: ['@babel/plugin-transform-object-assign'],
            presets: [
              [
                '@babel/preset-env',
                {
                  'useBuiltIns': 'entry',
                  'corejs': {
                    'version': 3,
                    'proposals': true
                  }
                }
              ]
            ]
          }
        }
      }
    ]
  }
};
