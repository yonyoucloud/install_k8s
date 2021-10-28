module.exports = {
  chainWebpack: config => {
    // 移除 preload 插件
    config.plugins.delete('preload-index')
  },
  pages: {
    index: {
      entry: 'src/main.js',
      env: process.env.NODE_ENV
    }
  },
  publicPath: './',
  assetsDir: './',
}