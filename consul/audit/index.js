var express = require('express')
var app = express()

app.get('/', function (req, res) {
    res.send('Hello World!')
})
app.get('/health', function (req, res) {
    res.send('I am healthy!')
})

app.listen(8100, function () {
    console.log('Example app listening on port 8100!')
})