module.exports = {
    runtimeCompiler: true,
    productionSourceMap: false,
    publicPath: process.env.NODE_ENV === 'production'
        ? './'
        : '/',
    devServer: {
        disableHostCheck: true,
        proxy: {
            '/dep/*': {
                target:  'http://localhost:3000',
                secure: false,
            },
            '/list/*': {
                target:  'http://localhost:3000',
                secure: false,
            },
            // '/ws*': {
            //     target:  'http://localhost:3000',
            //     secure: false,
            // },
        }
    },
}
