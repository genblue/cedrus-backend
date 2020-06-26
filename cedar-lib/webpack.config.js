function buildConfig(env) {
    return require("./config/webpack." + Object.keys(env)[0] + ".js");
}

module.exports = buildConfig;
