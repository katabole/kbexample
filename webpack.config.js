const Webpack = require("webpack");
const Glob = require("glob");
const CopyWebpackPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const { WebpackManifestPlugin}  = require("webpack-manifest-plugin");
const TerserPlugin = require("terser-webpack-plugin");
const LiveReloadPlugin = require('webpack-livereload-plugin');

const configurator = {
  entries: function(){
    return {
      application: [
        './node_modules/jquery-ujs/src/rails.js',
        './public/assets/css/application.scss',
        './public/assets/js/application.js',
      ],
    }
  },

  plugins() {
    var plugins = [
      new Webpack.ProvidePlugin({$: "jquery",jQuery: "jquery"}),
      new MiniCssExtractPlugin({filename: "[name].[contenthash].css"}),
      new CopyWebpackPlugin([{from: "./public/assets",to: ""}], {copyUnmodified: true, ignore: ["css/**", "js/**", "src/**"]}),
      new Webpack.LoaderOptionsPlugin({minimize: true,debug: false}),
      new WebpackManifestPlugin({
        fileName: `${__dirname}/public/dist/manifest.json`
      })
    ];

    return plugins
  },

  moduleOptions: function() {
    return {
      rules: [
        {
          test: /\.s[ac]ss$/,
          use: [
            MiniCssExtractPlugin.loader,
            { loader: "css-loader", options: {sourceMap: true}},
            // the quietDeps is to suppress deprecation warnings from bootstrap as it uses
            // old operations like "Using / for division"
            { loader: "sass-loader", options: {sourceMap: true, sassOptions: { quietDeps: true }}}
          ]
        },
        { test: /\.tsx?$/, use: "ts-loader", exclude: /node_modules/},
        { test: /\.jsx?$/,loader: "babel-loader",exclude: /node_modules/ },
        { test: /\.(woff|woff2|ttf|svg)(\?v=\d+\.\d+\.\d+)?$/,use: "url-loader"},
        { test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,use: "file-loader" },
        { test: require.resolve("jquery"),use: "expose-loader?exposes=jQuery!expose-loader?exposes=$"},
        { test: /\.go$/, use: "gopherjs-loader"}
      ]
    }
  },

  buildConfig: function(){
	if (process.env.DEPLOY_ENV == "development" || process.env.DEPLOY_ENV == "test") {
		env = "development"
	} else {
		env = "production"
	}

    var config = {
      mode: env,
      entry: configurator.entries(),
      output: {
        filename: "[name].[fullhash].js",
		// Prefix our entries in manifest.json with "assets/" to match where Go expects to find them
        publicPath: "assets",
		path: `${__dirname}/public/dist/assets`,
        clean: true,
        library: {
			name: 'App',
			type: 'var',
		},
      },
      plugins: configurator.plugins(),
      module: configurator.moduleOptions(),
      resolve: {
        extensions: ['.ts', '.js', '.json']
      }
    }

    if (env === "development") {
      config.plugins.push(new LiveReloadPlugin({appendScriptTag: true}))
      return config
    }

    const terser = new TerserPlugin({
      terserOptions: {
        mangle: {keep_fnames: true},
        output: {comments: false},
        compress: {}
      }
    })

    config.optimization = {
      minimizer: [terser]
    }

    return config
  }
}

module.exports = configurator.buildConfig()
