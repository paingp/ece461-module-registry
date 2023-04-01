'use strict';


/**
 * Create an access token.
 *
 * body AuthenticationRequest 
 * returns AuthenticationToken
 **/
exports.createAuthToken = function(body) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = "";
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 * Delete all versions of this package.
 *
 * name PackageName 
 * xAuthorization AuthenticationToken  (optional)
 * no response value expected for this operation
 **/
exports.packageByNameDelete = function(name,xAuthorization) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 * Return the history of this package (all versions).
 *
 * name PackageName 
 * xAuthorization AuthenticationToken  (optional)
 * returns List
 **/
exports.packageByNameGet = function(name,xAuthorization) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = [ {
  "Action" : "CREATE",
  "User" : {
    "name" : "Alfalfa",
    "isAdmin" : true
  },
  "PackageMetadata" : {
    "Version" : "1.2.3",
    "ID" : "ID",
    "Name" : "Name"
  },
  "Date" : "2023-03-23T23:11:15Z"
}, {
  "Action" : "CREATE",
  "User" : {
    "name" : "Alfalfa",
    "isAdmin" : true
  },
  "PackageMetadata" : {
    "Version" : "1.2.3",
    "ID" : "ID",
    "Name" : "Name"
  },
  "Date" : "2023-03-23T23:11:15Z"
} ];
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 * Get any packages fitting the regular expression.
 * Search for a package using regular expression over package names and READMEs. This is similar to search by name.
 *
 * body String 
 * regex PackageRegEx 
 * xAuthorization AuthenticationToken  (optional)
 * returns List
 **/
exports.packageByRegExGet = function(body,regex,xAuthorization) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = [ {
  "Version" : "1.2.3",
  "ID" : "ID",
  "Name" : "Name"
}, {
  "Version" : "1.2.3",
  "ID" : "ID",
  "Name" : "Name"
} ];
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 *
 * body PackageData 
 * xAuthorization AuthenticationToken 
 * returns Package
 **/
exports.packageCreate = function(body,xAuthorization) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = {
  "metadata" : {
    "Version" : "1.2.3",
    "ID" : "ID",
    "Name" : "Name"
  },
  "data" : {
    "Content" : "Content",
    "JSProgram" : "JSProgram",
    "URL" : "URL"
  }
};
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 * Delete this version of the package.
 *
 * id PackageID Package ID
 * xAuthorization AuthenticationToken  (optional)
 * no response value expected for this operation
 **/
exports.packageDelete = function(id,xAuthorization) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 *
 * id PackageID 
 * xAuthorization AuthenticationToken  (optional)
 * returns PackageRating
 **/
exports.packageRate = function(id,xAuthorization) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = {
  "GoodPinningPractice" : 7.061401241503109,
  "NetScore" : 0.8008281904610115,
  "ResponsiveMaintainer" : 5.637376656633329,
  "LicenseScore" : 2.3021358869347655,
  "RampUp" : 5.962133916683182,
  "GoodEngineeringProcess" : 9.301444243932576,
  "BusFactor" : 6.027456183070403,
  "Correctness" : 1.4658129805029452
};
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 * Interact with the package with this ID
 * Return this package.
 *
 * id PackageID ID of package to fetch
 * xAuthorization AuthenticationToken  (optional)
 * returns Package
 **/
exports.packageRetrieve = function(id,xAuthorization) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = {
  "metadata" : {
    "Version" : "1.2.3",
    "ID" : "ID",
    "Name" : "Name"
  },
  "data" : {
    "Content" : "Content",
    "JSProgram" : "JSProgram",
    "URL" : "URL"
  }
};
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 * Update this content of the package.
 * The name, version, and ID must match.  The package contents (from PackageData) will replace the previous contents.
 *
 * body Package 
 * id PackageID 
 * xAuthorization AuthenticationToken  (optional)
 * no response value expected for this operation
 **/
exports.packageUpdate = function(body,id,xAuthorization) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 * Get the packages from the registry.
 * Get any packages fitting the query. Search for packages satisfying the indicated query.  If you want to enumerate all packages, provide an array with a single PackageQuery whose name is \"*\".  The response is paginated; the response header includes the offset to use in the next query.
 *
 * body List 
 * offset EnumerateOffset Provide this for pagination. If not provided, returns the first page of results. (optional)
 * xAuthorization AuthenticationToken  (optional)
 * returns List
 **/
exports.packagesList = function(body,offset,xAuthorization) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = [ {
  "Version" : "1.2.3",
  "ID" : "ID",
  "Name" : "Name"
}, {
  "Version" : "1.2.3",
  "ID" : "ID",
  "Name" : "Name"
} ];
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 * Reset the registry
 * Reset the registry to a system default state.
 *
 * xAuthorization AuthenticationToken  (optional)
 * no response value expected for this operation
 **/
exports.registryReset = function(xAuthorization) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}

