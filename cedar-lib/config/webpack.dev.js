const webpack = require('webpack');
const merge = require('webpack-merge');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

const common = require('./webpack.common.js')();

module.exports = merge(common, {
    mode: 'development',
    plugins:[
        new webpack.DefinePlugin({
            PRODUCTION: JSON.stringify(false),
            VERSION: JSON.stringify('VERSION')
        }),
        new CleanWebpackPlugin({ cleanStaleWebpackAssets: false })
    ],
    devtool: 'inline-source-map',
    devServer: {
        contentBase: './../assets/static/assets/js/'
    }
});
