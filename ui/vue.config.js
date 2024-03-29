module.exports = {
  pages: {
    index: {
      entry: 'src/pages/index/main.js',
      title: '知识库'
    },
    login: {
      entry: 'src/pages/login/main.js',
      title: '登录 | 知识库'
    },
    signup: {
      entry: 'src/pages/signup/main.js',
      title: '注册账号 | 知识库'
    }
  },
  devServer: {
    port: 3000
  },
  publicPath: '/',
  assetsDir: 'static'
}