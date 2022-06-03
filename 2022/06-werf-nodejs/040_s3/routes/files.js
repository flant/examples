//@ts-check
const express = require('express');
const router = express.Router();
const asyncHandler = require('express-async-handler');
const config = require('../config/minio.json');
const fileUpload = require('express-fileupload');
const {
  S3,
  GetObjectCommand,
  PutObjectCommand,
} = require('@aws-sdk/client-s3');
const { Readable } = require('stream');
const { ReadableWebToNodeStream } = require('readable-web-to-node-stream');

module.exports = (logger) => {
  const s3 = new S3({ ...config, logger });
  const key = 'thekey';

  router.get(
    '/download',
    asyncHandler(async (req, res) => {
      try {
        const cmd = new GetObjectCommand({
          Bucket: config.bucket,
          Key: key,
        });

        const file = await s3.send(cmd);
        const body = file.Body;
        if (!body) {
          throw new Error('absent object body');
        }

        if (body instanceof Readable) {
          body.pipe(res);
          return;
        }

        if (body instanceof ReadableStream) {
          new Readable(new ReadableWebToNodeStream(body)).pipe(res);
          return;
        }

        if (body instanceof Blob) {
          body.stream().pipe(res);
          return;
        }

        throw new Error('cannot handle S3 response body');
      } catch (e) {
        if (e.name === 'NoSuchKey') {
          res.status(404).send(`You haven't uploaded anything yet.\n`);
          return;
        }
        res.status(500).send(`Something went wrong: ${e.message}\n`);
      }
    })
  );

  router.post(
    '/upload',
    fileUpload(),
    asyncHandler(async (req, res) => {
      try {
        if (!req.files || !req.files.file) {
          res.status(400).send('You forgot to attach a file.\n');
          return;
        }

        let file = req.files.file;
        if (Array.isArray(file)) {
          file = file[0];
        }

        const cmd = new PutObjectCommand({
          Bucket: config.bucket,
          Key: key,
          //@ts-ignore
          Body: file.data,
        });
        await s3.send(cmd);
        res.status(200).send('File uploaded.\n');
      } catch (e) {
        res.status(500).send(`Something went wrong: ${e.message}\n`);
        return;
      }
    })
  );

  return router;
};
