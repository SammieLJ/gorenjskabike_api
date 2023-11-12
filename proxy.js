const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();

app.use('/api', (req, res, next) => {
  // Set CORS headers
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');
  res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept');
  next();
}, createProxyMiddleware({ target: 'https://trzic.musiclab.si', changeOrigin: true }));

app.listen(3000, () => {
  console.log('Proxy server is running on http://localhost:3000');
});
