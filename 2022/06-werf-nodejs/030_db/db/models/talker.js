'use strict';
const { Model } = require('sequelize');
module.exports = (sequelize, DataTypes) => {
  class Talker extends Model {
    /**
     * Helper method for defining associations.
     * This method is not a part of Sequelize lifecycle.
     * The `models/index` file will call this method automatically.
     */
    static associate(models) {
      // Define association here.
    }
  }
  Talker.init(
    {
      answer: DataTypes.STRING,
      name: DataTypes.STRING,
    },
    {
      sequelize,
      modelName: 'Talker',
    }
  );
  return Talker;
};
