const { createCanvas } = require('canvas')
const fs = require('fs')

async function DrawCandlistick() {
    const canvas = createCanvas(800,600)
    const ctx = canvas.getContext('2d')

    ctx.fillStyle = 'green'
    ctx.fillRect(100, 100, 50, 200)

    ctx.fillStyle = 'red'
    ctx.fillRect(100, 150,50,50)

    const buffer = canvas.toBuffer('image/png')
    fs.writeFileSync('candlestick.png', buffer)
}

module.exports = {DrawCandlistick:DrawCandlistick}