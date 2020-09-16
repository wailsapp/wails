const path = require('path');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
    entry: './src/main.js',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'dist'),
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader']
            },
            {
                test: /\.(jpe?g|png|ttf|eot|svg|woff(2)?)(\?[a-z0-9=&.]+)?$/,
                use: 'base64-inline-loader'
            }

        ]
    },
    plugins: [
        new CopyWebpackPlugin([
            {
                from: 'src/index.html',
                to: ''
            },
            {
                from: 'src/assets/css/main.css',
                to: ''
            }
        ])
    ]
};
