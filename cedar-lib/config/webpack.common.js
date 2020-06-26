const path = require('path');

module.exports = () => {
    return {
        entry: path.resolve(__dirname, '../index.js'),
        output: {
            library: 'CedrusLib',
            filename: 'cedrus-lib.js',
            path: path.resolve(__dirname, '../../assets/static/assets/js/'),
        },
        node: {
            global: true
        }
    };
};
