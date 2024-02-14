const {request} = require('https')

function klineData(url, symbol, interval) {
    const options = {
        hostname: url,
        port: 443,
        path: `/api/v3/klines?symbol=${symbol}&interval=${interval}`,
        method: 'get',
    }

    const req = request(options, (res) => {
        // console.log(`STATUS: ${res.statusCode}`)
        // console.log(`HEADERS: ${JSON.stringify(res.headers)}`)

        res.setEncoding('utf-8')

        let data = '';

        res.on('data', (chunk) => data += chunk)

        res.on('end', () => callback(symbol, data))
    })

    req.on('error', (error) => console.log(error))
    req.end()
    return "hi"
}

function callback(symbol, data) {
    // 取最後五個kline 資料
    data = JSON.parse(data).slice(-5)
    console.log(symbol)
    // 計算每個 data 內 open - close 的價位
    result = data.map( i => {
        result = i.slice(1,5)
        return result[0] - result[3]
    })
    sum = result.slice(0,4).reduce((sum, curr) => sum + Math.abs(curr), 0)
    if (Math.abs(result[4]) > sum/4) {
        return "ok"
    }else {
        return "not"
    }
}
module.exports = {KlineData:klineData}