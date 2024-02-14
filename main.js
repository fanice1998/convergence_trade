const binanceData = require('./binance/klineData.js');
// const tools = require("./tools/tools");

symbols = ["ETHUSDT", "OPUSDT", "BTCUSDT", "SUIUSDT"]
for (const symbol of symbols) {
    binanceData.KlineData("api3.binance.com", symbol, "1h")
}
// binanceData.KlineData("api3.binance.com", "ETHUSDT","1h");
// tools.DrawCandlistick()