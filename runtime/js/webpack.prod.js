/* eslint-disable */

const path = require('path');

module.exports = {
  entry: './core/main',
  mode: 'production',
  output: {
    path: path.resolve(__dirname, '..', 'assets'),
    filename: 'wails.js'
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
