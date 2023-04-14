'use strict';

var utils = require('../utils/writer.js');
var Default = require('../service/DefaultService');

module.exports.createAuthToken = function createAuthToken (req, res, next, body) {
  Default.createAuthToken(body)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packageByNameDelete = function packageByNameDelete (req, res, next, name, xAuthorization) {
  Default.packageByNameDelete(name, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packageByNameGet = function packageByNameGet (req, res, next, name, xAuthorization) {
  Default.packageByNameGet(name, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packageByRegExGet = function packageByRegExGet (req, res, next, body, regex, xAuthorization) {
  Default.packageByRegExGet(body, regex, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packageCreate = function packageCreate (req, res, next, body, xAuthorization) {
  Default.packageCreate(body, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packageDelete = function packageDelete (req, res, next, id, xAuthorization) {
  Default.packageDelete(id, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packageRate = function packageRate (req, res, next, id, xAuthorization) {
  Default.packageRate(id, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packageRetrieve = function packageRetrieve (req, res, next, id, xAuthorization) {
  Default.packageRetrieve(id, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packageUpdate = function packageUpdate (req, res, next, body, id, xAuthorization) {
  Default.packageUpdate(body, id, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.packagesList = function packagesList (req, res, next, body, offset, xAuthorization) {
  Default.packagesList(body, offset, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.registryReset = function registryReset (req, res, next, xAuthorization) {
  Default.registryReset(xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
