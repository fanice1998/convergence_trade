const assert = require('assert')
const  klineData = require('../binance/klineData')

describe('klineData', () => {
    it('should request binance klineData with symbol', () => {
        // assert.strictEqual(klineData.klineData("hello"), "hello")
        const result = klineData.klineData("https://api3.binance.com", "ETHUSDT","1h")
        assert.ok(result !== null && result !== undefined && typeof result == "string")
    });
})