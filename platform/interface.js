
class platform {
    constructor(name) {
        this.name = name;
    }

    echo() {
        console.log(`${this.name} say hello`)
    }
}


module.exports = platform;