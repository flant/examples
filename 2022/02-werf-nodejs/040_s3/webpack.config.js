const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');

module.exports = {
  mode: 'development',
  devServer: { static: '/dist' },
  devtool: 'inline-source-map',

  entry: {
    index: path.resolve(__dirname, 'public/assets/javascripts/index.js'),
    image: path.resolve(__dirname, 'public/assets/javascripts/image.js'),
  },

  output: {
    filename: 'js/[name].[contenthash].js',
    path: path.resolve(__dirname, 'dist'),
    clean: true,
    publicPath: '/static/',
  },

  plugins: [

    new MiniCssExtractPlugin({
      filename: 'css/[name].[contenthash].css',
    }),

    new HtmlWebpackPlugin({
      title: 'werf first app',
      template: 'public/pages/index.html',
      chunks: ['index'],
    }),

    new HtmlWebpackPlugin({
      title: 'werf logo',
      filename: 'image.html',
      template: 'public/pages/image.html',
      chunks: ['image'],
    }),
  ],
  module: {
    rules: [
      {
        test: /\.css$/i,
        use: [MiniCssExtractPlugin.loader, 'css-loader'],
      },
      {
        test: /\.(png|svg|jpg|jpeg|gif)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'media/[contenthash][ext]',
        },
      },
    ],
  },

  optimization: {
    moduleIds: 'deterministic',
    runtimeChunk: 'single',
    splitChunks: {
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          name: 'vendors',
          chunks: 'all',
        },
      },
    },
    minimizer: [new CssMinimizerPlugin()],
  },
};
