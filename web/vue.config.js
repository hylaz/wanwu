"use strict";

const path = require("path");
const { name } = require('./package.json')
const webpack = require('webpack')
const VersionInfoPlugin = require('./build/version-plugin');
function resolve(dir) {
  return path.join(__dirname, dir);
}

const CompressionWebpackPlugin = require('compression-webpack-plugin')
const isProdOrTest = process.env.NODE_ENV !== 'development'

module.exports = {
  // 基础配置 详情看文档
  publicPath: process.env.VUE_APP_BASE_PATH + "/aibase",
  outputDir: "dist",
  assetsDir: "static",
  lintOnSave: process.env.NODE_ENV === "development",
  productionSourceMap: false,//源码映射
  chainWebpack(config){
    config.module
        .rule('md')
        .test(/\.md$/)
        .use('html-loader')
        .loader('html-loader')
        .end()
        .use('markdown-loader')
        .loader('markdown-loader')
        .end()

    config.plugins.delete('prefetch')
    if (isProdOrTest) {
      // 对超过10kb的文件gzip压缩
      config.plugin('compressionPlugin').use(
        new CompressionWebpackPlugin({
          test: /\.(css|html)$/,
          threshold: 10240,
        })
      )
    }
    config.module
      .rule('svg')
      .exclude.add(resolve('src/assets/icons'))	//svg文件位置
      .end()
    config.module
      .rule('icons')
      .test(/\.svg$/)
      .include.add(resolve('src/assets/icons'))	//svg文件位置
      .end()
      .use('svg-sprite-loader')
      .loader('svg-sprite-loader')
      .options({
          symbolId: 'icon-[name]'
      })
      .end()
  },
  devServer: {
    port: 8080,
    open: false,
    hot:true,
    overlay: {
      warnings: false,
      errors: true,
    },
    headers: {
        'Access-Control-Allow-Origin': '*',
    },
    proxy: {
      "/openAi":{
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/workflow/api":{
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/user/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/service/url/openurl/v1":{
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/service/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/training/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/resource/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/datacenter/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/modelprocess/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/expand/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/record/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/img": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/konwledgeServe": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/proxyupload": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/use/model/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      },
      "/prompt/api": {
        target: "http://192.168.0.64:8081",
        changeOrigin: true,
        secure: false,
      }
    },
  },
  css: {
    sourceMap: false,
    loaderOptions: {
      sass: {
        prependData: `@import "~@/style/theme/vars_blue.scss";@import "~@/style/theme/common.scss";` // 假设variables.scss位于src/styles目录下
      }
    }
  },
  configureWebpack: {
    //    @路径走src文件夹
    module: {
      rules: [
        {
          test: /\.swf$/,
          loader: "url-loader",
          options: {
            limit: 10000,
            name: "static/media/[name].[hash:7].[ext]",
          },
        },
      ],
    },
    resolve: {
      alias: {
        'vue$': 'vue/dist/vue.esm.js',
        "@": resolve("src"),
          "@common": resolve("common"),
      },
    },
      output: {
          // 把子应用打包成 umd 库格式(必须)
          library: `${name}-[name]`,
          libraryTarget: 'umd',
          jsonpFunction: `webpackJsonp_${name}`,
      },
      plugins:[
        new webpack.optimize.LimitChunkCountPlugin({
            maxChunks: 10, // 来限制 chunk 的最大数量
          }),
          new webpack.optimize.MinChunkSizePlugin({
            minChunkSize: 50000 // Minimum number of characters
          }),
          new VersionInfoPlugin()
      ]
  },
};
