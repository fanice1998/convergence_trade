const myInterface = require("../platform/interface.js")

function run(){
    const platform = new myInterface("binance")

    platform.echo()
}

module.exports = run;