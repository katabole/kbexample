const Webpack = require("webpack");
const Glob = require("glob");
const CopyWebpackPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const { WebpackManifestPlugin}  = require("webpack-manifest-plugin");
const TerserPlugin = require("terser-webpack-plugin");
const LiveReloadPlugin = require('webpack-livereload-plugin');

const configurator = {
  entries: function(){
    var entries = {
      application: [
        './node_modules/jquery-ujs/src/rails.js',
        './public/assets/css/application.scss',
      ],
    }

    Glob.sync("./assets/*/*.*").forEach((entry) => {
      if (entry === './public/assets/css/application.scss') {
        return
      }

      let key = entry.replace(/(\.\/assets\/(src|js|css|go)\/)|\.(ts|js|s[ac]ss|go)/g, '')
      if(key.startsWith("_") || (/(ts|js|s[ac]ss|go)$/i).test(entry) == false) {
        return
      }

      if( entries[key] == null) {
        entries[key] = [entry]
        return
      }

      entries[key].push(entry)
    })
    return entries
  },

  plugins() {
    var plugins = [
      new Webpack.ProvidePlugin({$: "jquery",jQuery: "jquery"}),
      new MiniCssExtractPlugin({filename: "[name].[contenthash].css"}),
      new CopyWebpackPlugin([{from: "./public/assets",to: ""}], {copyUnmodified: true,ignore: ["css/**", "js/**", "src/**"] }),
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
        { test: require.resolve("jquery"),use: "expose-loader?exposes=jQuery!expose-loader?$"},
        { test: /\.go$/, use: "gopherjs-loader"}
      ]
    }
  },

  buildConfig: function(){
    // NOTE: If you are having issues with this not being set "properly", make
    // sure your GO_ENV is set properly as `buffalo build` overrides NODE_ENV
    // with whatever GO_ENV is set to or "development".
    const env = process.env.NODE_ENV || "development";

    var config = {
      mode: env,
      entry: configurator.entries(),
      output: {
        filename: "[name].[fullhash].js",
		// Prefix our entries in manifest.json with "assets/" to match where Go expects to find them
        publicPath: "assets",
		path: `${__dirname}/public/dist/assets`,
        clean: true,
      },
      plugins: configurator.plugins(),
      module: configurator.moduleOptions(),
      resolve: {
        extensions: ['.ts', '.js', '.json']
      }
    }

    if( env === "development" ){
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
