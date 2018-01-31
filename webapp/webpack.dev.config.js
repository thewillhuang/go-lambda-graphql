const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const baseConfig = require('./webpack.config');

const port = 4000;

module.exports = Object.assign(baseConfig, {
  devServer: {
    host: '0.0.0.0',
    headers: { 'Access-Control-Allow-Origin': '*' },
    port,
    compress: true,
    proxy: {
      '/': 'http://localhost:3000',
    },
    overlay: {
      warnings: true,
      errors: true,
    },
  },
  output: {
    path: '/',
    publicPath: `http://localhost:${port}/`,
    filename: '[name].bundle.js',
    libraryTarget: 'umd',
    devtoolModuleFilenameTemplate: '/[absolute-resource-path]',
  },
  devtool: 'eval-source-map',
  plugins: [
    new ExtractTextPlugin({
      filename: '[name].[contenthash].css',
      allChunks: true,
    }),
    new webpack.ContextReplacementPlugin(/graphql-language-service-interface[/\\]dist/, /\.js$/),
    new HtmlWebpackPlugin({
      template: '../template.html',
    }),
    new webpack.NamedModulesPlugin(),
  ],
});
