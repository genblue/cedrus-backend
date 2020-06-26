const webpack = require('webpack');
const merge = require('webpack-merge');
const common = require('./webpack.common.js')();

module.exports = merge(common, {
    mode: 'production',
    plugins:[
        new webpack.DefinePlugin({
            PRODUCTION: JSON.stringify(true),
            VERSION: JSON.stringify('VERSION')
        })
    ],
});
