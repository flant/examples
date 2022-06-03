'use strict';

const fs = require('fs');
const path = require('path');
const Sequelize = require('sequelize');
const basename = path.basename(__filename);
const env = process.env.NODE_ENV || 'development';
const config = require(__dirname + '/../../config/database.json')[env];

module.exports = function init(logger) {
  const configWithInjections = {
    ...config,
    logging: logger,
  };

  let sequelize;
  if (config.use_env_variable) {
    sequelize = new Sequelize(
      process.env[config.use_env_variable],
      configWithInjections
    );
  } else {
    sequelize = new Sequelize(
      config.database,
      config.username,
      config.password,
      configWithInjections
    );
  }

  const db = {};

  fs.readdirSync(__dirname)
    .filter((file) => {
      return (
        file.indexOf('.') !== 0 && file !== basename && file.slice(-3) === '.js'
      );
    })
    .forEach((file) => {
      const model = require(path.join(__dirname, file))(
        sequelize,
        Sequelize.DataTypes
      );
      db[model.name] = model;
    });

  Object.keys(db).forEach((modelName) => {
    if (db[modelName].associate) {
      db[modelName].associate(db);
    }
  });

  db.sequelize = sequelize;
  db.Sequelize = Sequelize;

  return db;
};
