//@ts-check
const express = require('express');
const router = express.Router();
const asyncHandler = require('express-async-handler');

module.exports = (db) => {
  router.get(
    '/say',
    asyncHandler(async (req, res) => {
      try {
        const talker = await db.sequelize.transaction(async (t) => {
          return await db.Talker.findOne({ where: { id: 1 }, transaction: t });
        });

        if (!talker) {
          res.send(`I have nothing to say.\n`);
          return;
        }

        res.send(`${talker.answer}, ${talker.name}!\n`);
      } catch (e) {
        res.status(500).send(`Something went wrong: ${e.message}\n`);
      }
    })
  );

  router.get(
    '/remember',
    asyncHandler(async (req, res) => {
      const { answer, name } = req.query;
      if (!answer) {
        res.status(422).send('You forgot the answer :(\n');
        return;
      }

      if (!name) {
        res.status(422).send('You forgot the name :(\n');
        return;
      }

      try {
        await db.sequelize.transaction(async (t) => {
          const [talker] = await db.Talker.findOrCreate({
            where: { id: 1 },
            transaction: t,
          });
          talker.set({ answer, name });
          await talker.save({ transaction: t });
        });

        res.send(`Got it.\n`);
      } catch (e) {
        res.status(500).send(`Something went wrong: ${e.message}\n`);
      }
    })
  );

  return router;
};
