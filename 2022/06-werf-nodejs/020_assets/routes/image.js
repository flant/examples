const express = require('express');
const router = express.Router();

/* Image page. */
router.get('/', express.static('dist/image.html'));

module.exports = router;
