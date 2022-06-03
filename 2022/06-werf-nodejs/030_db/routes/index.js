const express = require('express');
const path = require('path');
const router = express.Router();

/* GET home page. */
router.get('/', function (req, res, next) {
  res.sendFile(path.resolve(process.cwd(), 'dist', 'index.html'));
});

router.get('/err', function (req, res, next) {
  throw new Error('Hello from an unhandler error');
});

module.exports = router;
