const path = require('path');
const VueLoaderPlugin = require('vue-loader/lib/plugin');

module.exports = {
    entry: './src/index.js',
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'dist'),
    },

    module: {
        rules: [
            {
                test: /\.vue$/,
                loader: 'vue-loader',
            },
            
            {
                test: /\.m?js$/,
                exclude: /(node_modules|bower_components)/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env']
                    },
                },
            },

             {
                 test: /\.css$/i,
                 use: ['style-loader', 'css-loader'],
             },
        ],
    },

    resolve: {
        alias: {
            'vue$': 'vue/dist/vue.esm.js,',
        },
    },
    
    plugins: [
        new VueLoaderPlugin(),
    ]
};
