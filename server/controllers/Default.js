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

module.exports.packageByNameDelete = function packageByNameDelete (req, res, next, xAuthorization, name) {
  Default.packageByNameDelete(xAuthorization, name)
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

module.exports.packageByRegExGet = function packageByRegExGet (req, res, next, body, xAuthorization) {
  Default.packageByRegExGet(body, xAuthorization)
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

module.exports.packageDelete = function packageDelete (req, res, next, xAuthorization, id) {
  Default.packageDelete(xAuthorization, id)
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

module.exports.packageRetrieve = function packageRetrieve (req, res, next, xAuthorization, id) {
  Default.packageRetrieve(xAuthorization, id)
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

module.exports.RegistryReset = function RegistryReset (req, res, next, xAuthorization) {

  let auth_code = "ABC"
  xAuthorization = req.headers['x-authorization'];

  if (xAuthorization == auth_code){
    res.status(200).send("Registry is reset.");
    
  } else if (xAuthorization == "") {
    res.status(400).send("There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.");
  } else {
    res.status(401).send("You do not have permission to reset the registry.");
  }

  Default.RegistryReset(xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
