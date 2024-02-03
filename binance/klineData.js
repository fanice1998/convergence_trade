const {request} = require('https')

function klineData(url, symbol, interval) {
    const options = {
        hostname: url,
        port: 443,
        path: `/api/v3/klines?symbol=${symbol}&interval=${interval}`,
        method: 'get',
    }

    const req = request(options, (res) => {
        console.log(`STATUS: ${res.statusCode}`)
        console.log(`HEADERS: ${JSON.stringify(res.headers)}`)

        res.setEncoding('utf-8')

        let data = '';

        res.on('data', (chunk) => data += chunk)

        res.on('end', () => console.log(data))
    })

    req.on('error', (error) => console.log(error))
    req.end()
    return "hi"
}

module.exports = {KlineData:klineData}