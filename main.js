const binanceData = require('./binance/klineData.js');
const tools = require("./tools");

binanceData.KlineData("api3.binance.com", "ETHUSDT","1h");
tools.drawCandlistick()