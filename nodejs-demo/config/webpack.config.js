
const path = require("path");
const base = path.join(__dirname, "..");
const HtmlWebPackPlugin = require("html-webpack-plugin");

module.exports = {
  mode: "production",
  entry: path.resolve(base, "src", "index.js"),
  output: {
    filename: "bundle.js",
    path: path.resolve(base, "dist")
  },
  resolve: {
    extensions: [".js", ".jsx"]
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: { // babel 转义的配置选项
            babelrc: false,
            presets: [
              require.resolve('@babel/preset-react'),
              [require.resolve('@babel/preset-env'), { modules: false }],
            ],
            cacheDirectory: true,
          },
        },
      },
    ],
  },
  plugins: [
    new HtmlWebPackPlugin({
      template: "src/index.html",
      filename: "index.html"
    })
  ]
};

