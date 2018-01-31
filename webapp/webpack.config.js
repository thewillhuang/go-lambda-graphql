const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HashOutput = require('webpack-plugin-hash-output');

const inlineSizeLimit = 1000;

module.exports = {
  entry: {
    index: [
      path.resolve('src/index'),
    ],
  },

  module: {
    rules: [
      {
        test: /\.jsx?$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        options: {
          babelrc: false,
          plugins: ['ramda', 'lodash', '@babel/transform-runtime'],
          presets: [
            '@babel/flow',
            ['@babel/env', {
              targets: {
                forceAllTransforms: true,
              },
              modules: false,
            }],
            '@babel/react',
            '@babel/stage-0',
          ],
        },
      },
      {
        test: /\.s?css$/,
        exclude: /node_modules/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: [{
            loader: 'css-loader',
            options: {
              modules: true,
              importLoaders: 1,
              localIdentName: '[sha512:hash:base64:8]',
            },
          },
          { loader: 'postcss-loader' }],
        }),
      },
      {
        test: /\.s?css$/,
        include: /node_modules/,
        use: ExtractTextPlugin.extract({
          fallback: 'style-loader',
          use: [{
            loader: 'css-loader',
          },
          { loader: 'postcss-loader' }],
        }),
      },
      {
        test: /\.(eot|ttf|woff|woff2|jpe?g|png|gif|svg)$/,
        loader: 'url-loader',
        options: {
          limit: inlineSizeLimit,
          name: '[sha512:hash:base64:7].[ext]',
        },
      },
    ],
  },

  output: {
    path: `${__dirname}/build/`,
    filename: '[name].[chunkhash].js',
    libraryTarget: 'umd',
  },

  resolve: {
    extensions: ['.js', '.jsx', '.json', '.css', '.scss'],
  },

  plugins: [
    new ExtractTextPlugin({
      filename: '[name].[contenthash].css',
      allChunks: true,
    }),
    new webpack.ContextReplacementPlugin(/graphql-language-service-interface[/\\]dist/, /\.js$/),
    new HtmlWebpackPlugin({
      template: '../template.html',
      minify: {
        removeComments: true,
        collapseWhitespace: true,
        conservativeCollapse: true,
      },
    }),
    new webpack.optimize.ModuleConcatenationPlugin(),
    new webpack.HashedModuleIdsPlugin({
      hashFunction: 'sha256',
    }),
    new HashOutput({}),
  ],
};
