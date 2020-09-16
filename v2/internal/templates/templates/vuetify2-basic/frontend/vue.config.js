let cssConfig = {};

if (process.env.NODE_ENV == 'production') {
	cssConfig = {
		extract: {
			filename: '[name].css',
			chunkFilename: '[name].css'
		}
	};
}

module.exports = {
    indexPath: 'index.html',
    publicPath: 'public',
    // disable hashes in filenames
    filenameHashing: false,
    // delete HTML related webpack plugins
    chainWebpack: config => {
      config.plugins.delete('preload');
      config.plugins.delete('prefetch');
      config
            .plugin('html')
            .tap(args => {
                args[0].template = 'public/index.html';
                args[0].inject = false;
                args[0].cache = false;
                args[0].minify = false;
                args[0].filename = 'index.html';
                return args;
            });

		let limit = 9999999999999999;
		config.module
			.rule('images')
			.test(/\.(png|gif|jpg)(\?.*)?$/i)
			.use('url-loader')
			.loader('url-loader')
			.tap(options => Object.assign(options, { limit: limit }));
		config.module
			.rule('fonts')
			.test(/\.(woff2?|eot|ttf|otf|svg)(\?.*)?$/i)
			.use('url-loader')
			.loader('url-loader')
			.options({
				limit: limit
			});
	},
	css: cssConfig,
	configureWebpack: {
		output: {
			filename: '[name].js'
		},
		optimization: {
			splitChunks: false
		}
	},
	devServer: {
		disableHostCheck: true
	}
};
