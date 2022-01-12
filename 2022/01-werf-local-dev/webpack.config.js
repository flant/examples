const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');

module.exports = {
  mode: 'development',
  devServer: { static: '/dist' },
  devtool: 'inline-source-map',

  // [<en>] Two entrypoints to collect bundles to. Thus we split CSS and JS between two pages.
  // [<ru>] Две точки входа, для которых мы собираем бандлы. Таким образом мы разделяем стили и
  // [<ru>] скрипты для двух страниц.
  entry: {
    index: path.resolve(__dirname, 'public/assets/javascripts/index.js'),
    image: path.resolve(__dirname, 'public/assets/javascripts/image.js'),
  },

  // [<ru>] Определяем, как называются переработанные скрипты и с каким префиксом ожидать
  // [<ru>] статические файлы.
  // [<en>] Define script paths and the prefix for static files.
  output: {
    filename: 'js/[name].[contenthash].js',
    path: path.resolve(__dirname, 'dist'),
    clean: true,
    publicPath: '/static/',
  },

  plugins: [
    // [<ru>] Подключаем плагин для выделения CSS в отдельные файлы
    // [<en>] App plugin to extract CSS files
    new MiniCssExtractPlugin({
      filename: 'css/[name].[contenthash].css',
    }),
    // [<ru>] Генерируем главную страницу из ее шаблона и указываем, что подключить нужно только
    // [<ru>] точку входа index, включая CSS и JS. Так мы избежим подключения скриптов
    // [<ru>] и стилей от страницы с изображением.
    // [<en>] Generate home page from its template and include only index entry to it including
    // [<en>] JS and CSS. Thus we avoid including JS and CSS from image page.
    new HtmlWebpackPlugin({
      title: 'Werf Guide App',
      template: 'public/pages/index.html',
      chunks: ['index'],
    }),
    // [<ru>] То же самое для страницы c изображением, но еще явно указываем конечное имя HTML-файла.
    // [<en>] The same for the image page, except explicit destination name for the HTML file.
    new HtmlWebpackPlugin({
      title: 'Werf Logo',
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
      // [<ru>] Изображения перенесем в каталог /static/media, причем часть "static" мы определили
      // [<ru>] в настройках "output".
      // [<en>] Put images to /static/media, the "static" part being defined in "output" settings.
      {
        test: /\.(png|svg|jpg|jpeg|gif)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'media/[contenthash][ext]',
        },
      },
    ],
  },
  // [<ru>] С помощью оптимизаций мы выделили скрипт с райнтаймом, чтобы уменьшить общий размер статики,
  // [<ru>] который нужно будет загрузить для разных страниц. Также добавили минимизацию CSS для production.
  // [<en>] Using optimizations, we extracted runtime JS to separate file to decrese the size of scripts
  // [<en>] that needs to be loaded in clents browser. Also, we added the minimization of CSS for production.
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
