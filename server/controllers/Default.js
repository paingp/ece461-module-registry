'use strict';

var utils = require('../utils/writer.js');
var Default = require('../service/DefaultService');

let auth_code = "ABC"

module.exports.CreateAuthToken = function CreateAuthToken (req, res, next, body) {

  res.status(501).send("This system does not support authentication.");

  // Default.CreateAuthToken(body)
  //   .then(function (response) {
  //     utils.writeJson(res, response);
  //   })
  //   .catch(function (response) {
  //     utils.writeJson(res, response);
  //   });
};

module.exports.PackageByNameDelete = function PackageByNameDelete (req, res, next, xAuthorization, name) {
  
  xAuthorization = req.headers['x-authorization'];

  if (xAuthorization == auth_code){
    res.status(200).send("Package is deleted.");
  }
  
  var name = req.openapi.pathParams.name;
  console.log(name);
  
  Default.PackageByNameDelete(xAuthorization, name)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.PackageByNameGet = function PackageByNameGet (req, res, next, name, xAuthorization) {
  
  xAuthorization = req.headers['x-authorization'];

  if (xAuthorization == auth_code){
    res.status(200).send("Return the package history.");
  }
  
  var name = req.openapi.pathParams.name;
  console.log(name);
  
  
  Default.PackageByNameGet(name, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.PackageByRegExGet = function PackageByRegExGet (req, res, next, body, xAuthorization) {

  console.log(req);

  xAuthorization = req.headers['x-authorization'];

  if (xAuthorization == auth_code){
    res.status(200).send("Return a list of packages.");
  }
  
  var regex_string = req.body.RegEx;
  
  Default.PackageByRegExGet(body, xAuthorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.PackageCreate = function PackageCreate (req, res, next, body, xAuthorization) {
  
  xAuthorization = req.headers['x-authorization'];

  if (xAuthorization == auth_code){
    res.status(201).send("Check the ID in the returned metadata for the official ID.");
  }
  
  // Default.PackageCreate(body, xAuthorization)
  //   .then(function (response) {
  //     utils.writeJson(res, response);
  //   })
  //   .catch(function (response) {
  //     utils.writeJson(res, response);
  //   });
};

module.exports.PackageDelete = function PackageDelete (req, res, next, xAuthorization, id) {
  
  xAuthorization = req.headers['x-authorization'];

  var id = req.openapi.pathParams.id;
  console.log(id);

  if (xAuthorization == auth_code){
    res.status(200).send("Return the package. Content is required.");
  }
  
  // Default.PackageDelete(xAuthorization, id)
  //   .then(function (response) {
  //     utils.writeJson(res, response);
  //   })
  //   .catch(function (response) {
  //     utils.writeJson(res, response);
  //   });
};

module.exports.PackageRate = function PackageRate (req, res, next, id, xAuthorization) {
  
  xAuthorization = req.headers['x-authorization'];

  var id = req.openapi.pathParams.id;
  console.log(id);

  if (xAuthorization == auth_code){
    res.status(200).send("Return the rating. Only use this if each metric was computed successfully.");
  }
  
  // Default.PackageRate(id, xAuthorization)
  //   .then(function (response) {
  //     utils.writeJson(res, response);
  //   })
  //   .catch(function (response) {
  //     utils.writeJson(res, response);
  //   });
};

module.exports.PackageRetrieve = function PackageRetrieve (req, res, next, xAuthorization, id) {
  
  xAuthorization = req.headers['x-authorization'];

  var id = req.openapi.pathParams.id;

  if (xAuthorization == auth_code){
    res.status(200).send("Return the package. Content is required.");
  }
  
  Default.PackageRetrieve(xAuthorization, id)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.PackageUpdate = function PackageUpdate (req, res, next, body, id, xAuthorization) {
  
  xAuthorization = req.headers['x-authorization'];
  var id = req.openapi.pathParams.id;

  console.log(id);

  if (xAuthorization == auth_code){
    res.status(200).send("Return the package. Content is required.");
  }
  
  // Default.PackageUpdate(body, id, xAuthorization)
  //   .then(function (response) {
  //     utils.writeJson(res, response);
  //   })
  //   .catch(function (response) {
  //     utils.writeJson(res, response);
  //   });
};

module.exports.PackagesList = function PackagesList (req, res, next, body, offset, xAuthorization) {
  
  xAuthorization = req.headers['x-authorization'];

  var offset = req.query.offset;
  
  if (xAuthorization == auth_code){
    res.status(200).send("test output");
    console.log("working auth code");
  }

  // Need to eerror check 

  // Default.PackagesList(body, offset, xAuthorization)
  //   .then(function (response) {
  //     utils.writeJson(res, response);
  //   })
  //   .catch(function (response) {
  //     utils.writeJson(res, response);
  //   });
};

module.exports.RegistryReset = function RegistryReset (req, res, next, xAuthorization) {

  
  xAuthorization = req.headers['x-authorization'];

  if (xAuthorization == auth_code){
    res.status(200).send("Registry is reset.");
    // Call functionality 
  } else if (xAuthorization == "") {
    res.status(400).send("There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.");
  } else {
    res.status(401).send("You do not have permission to reset the registry.");
  }

  // Default.RegistryReset(xAuthorization)
  //   .then(function (response) {
  //     utils.writeJson(res, response);
  //   })
  //   .catch(function (response) {
  //     utils.writeJson(res, response);
  //   });
};
