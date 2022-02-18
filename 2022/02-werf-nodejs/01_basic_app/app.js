const express = require('express');
const path = require('path');
const cookieParser = require('cookie-parser');
const logger = require('morgan');

// [<snippet ping-route-import>]
const indexRouter = require('./routes/index');
const pingRouter = require('./routes/ping');
// [<endsnippet ping-route-import>]

const app = express();

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

// [<snippet ping-route-middleware>]
app.use('/', indexRouter);
app.use('/ping', pingRouter);
// [<endsnippet ping-route-middleware>]

module.exports = app;
