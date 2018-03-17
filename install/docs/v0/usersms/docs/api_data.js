define({ "api": [
  {
    "type": "GET",
    "url": "/ratings/users/{forUserID}",
    "title": "GetRatingsOnUser",
    "name": "Get_Ratings_On_User",
    "version": "0.1.0",
    "group": "Service",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token containing auth token e.g. &quot;Bearer [value.of.jwt]&quot;</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Param": [
          {
            "group": "URL Param",
            "type": "String",
            "optional": true,
            "field": "forUserID",
            "description": "<p>Filter ratings by ratee's userID. At least one of forUserID or byUserID must be provided.</p>"
          }
        ],
        "URL Query": [
          {
            "group": "URL Query",
            "type": "String",
            "optional": true,
            "field": "byUserID",
            "description": "<p>Filter ratings by rater's userID. At least one of forUserID or byUserID must be provided.</p>"
          },
          {
            "group": "URL Query",
            "type": "String",
            "optional": true,
            "field": "forSection",
            "description": "<p>Filter ratings by section which ratee was rated.</p>"
          },
          {
            "group": "URL Query",
            "type": "Integer",
            "optional": true,
            "field": "offset",
            "defaultValue": "0",
            "description": "<p>Index from which to fetch(inclusive).</p>"
          },
          {
            "group": "URL Query",
            "type": "Integer",
            "optional": true,
            "field": "count",
            "defaultValue": "10",
            "description": "<p>Number of items to fetch.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service",
    "success": {
      "fields": {
        "200 JSON Response": [
          {
            "group": "200 JSON Response",
            "type": "Object[]",
            "optional": false,
            "field": "ratings",
            "description": "<p>List of ratings (values indented below).</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ratings.ID",
            "description": "<p>Unique identifier of this rating.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ratings.forUserID",
            "description": "<p>Ratee' userID.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ratings.byUserID",
            "description": "<p>Rater's userID.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ratings.comment",
            "description": ""
          },
          {
            "group": "200 JSON Response",
            "type": "Integer",
            "size": "1-5",
            "optional": false,
            "field": "ratings.rating",
            "description": "<p>Rating awarded by rater to ratee.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ratings.created",
            "description": "<p>ISO8601 date of rating creation.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ratings.lastUpdated",
            "description": "<p>Last ISO8601 date of update.</p>"
          }
        ]
      }
    }
  },
  {
    "type": "GET",
    "url": "/users/{userID}",
    "title": "GetUser",
    "name": "Get_user",
    "version": "0.1.0",
    "group": "Service",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token containing auth token e.g. &quot;Bearer [value.of.jwt]&quot;</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Param": [
          {
            "group": "URL Param",
            "type": "String",
            "optional": true,
            "field": "userID",
            "description": "<p>ID of the user to fetch.</p>"
          }
        ],
        "URL Query": [
          {
            "group": "URL Query",
            "type": "Integer",
            "optional": true,
            "field": "offsetUpdateDate",
            "description": "<p>Earliest ISO8601 date that the user should have been updated. If the userID exists but the update date is earlier than this value then a 404 will be returned.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service",
    "success": {
      "fields": {
        "200 JSON Response": [
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>User's ID.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": ""
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ICEPhone",
            "description": "<p>User's (In Case of Emergency) phone number.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "allowedValues": [
              "\"MALE\"",
              "\"FEMALE\"",
              "\"OTHER\""
            ],
            "optional": false,
            "field": "gender",
            "description": ""
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "avatarURL",
            "description": "<p>User's profile picture URL.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "bio",
            "description": "<p>Brief description of user.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "Float",
            "size": "1-5",
            "optional": false,
            "field": "rating",
            "description": "<p>Overall rating of user.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date of user profile creation.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>last ISO8601 date when this profile was updated.</p>"
          }
        ]
      }
    }
  },
  {
    "type": "POST",
    "url": "/ratings/users/{forUserID}",
    "title": "RateUser",
    "name": "Rate_a_user",
    "version": "0.1.0",
    "group": "Service",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token containing auth token e.g. &quot;Bearer [value.of.jwt]&quot;</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "URL Param": [
          {
            "group": "URL Param",
            "type": "String",
            "optional": true,
            "field": "forUserID",
            "description": "<p>ID of the user to rate (ratee).</p>"
          }
        ],
        "JSON Request Body": [
          {
            "group": "JSON Request Body",
            "type": "Integer",
            "size": "1-5",
            "optional": false,
            "field": "rating",
            "description": "<p>The rating awarded by rater to ratee.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": true,
            "field": "comment",
            "description": "<p>Comment provided by rater.</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200 Response": [
          {
            "group": "200 Response",
            "optional": false,
            "field": "nil",
            "description": "<p>an empty body</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  },
  {
    "type": "get",
    "url": "/status",
    "title": "Status",
    "name": "Status",
    "version": "0.1.0",
    "group": "Service",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200": [
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Micro-service name.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "version",
            "description": "<p>http://semver.org version.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "description",
            "description": "<p>Short description of the micro-service.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "canonicalName",
            "description": "<p>Canonical name of the micro-service.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  },
  {
    "type": "PUT",
    "url": "/users/{userID}",
    "title": "UpdateUser",
    "name": "Update_User_Profile",
    "version": "0.1.0",
    "group": "Service",
    "description": "<p>All declared JSON values are used, including empty strings, except null values.</p>",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          },
          {
            "group": "Header",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Bearer token containing auth token e.g. &quot;Bearer [value.of.jwt]&quot;</p>"
          }
        ]
      }
    },
    "parameter": {
      "fields": {
        "JSON Request Body": [
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": true,
            "field": "name",
            "description": "<p>New Name.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "String",
            "optional": true,
            "field": "ICEPhone",
            "description": "<p>New (In Case of Emergency) phone number.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "String",
            "allowedValues": [
              "\"MALE\"",
              "\"FEMALE\"",
              "\"OTHER\""
            ],
            "optional": true,
            "field": "gender",
            "description": "<p>New gender.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "Object",
            "optional": true,
            "field": "avatarURL",
            "description": "<p>New profile picture URL.</p>"
          },
          {
            "group": "JSON Request Body",
            "type": "Object",
            "optional": true,
            "field": "bio",
            "description": "<p>New brief description of user.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service",
    "success": {
      "fields": {
        "200 JSON Response": [
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>User's ID.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": ""
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "ICEPhone",
            "description": "<p>User's (In Case of Emergency) phone number.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "allowedValues": [
              "\"MALE\"",
              "\"FEMALE\"",
              "\"OTHER\""
            ],
            "optional": false,
            "field": "gender",
            "description": ""
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "avatarURL",
            "description": "<p>User's profile picture URL.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "bio",
            "description": "<p>Brief description of user.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "Float",
            "size": "1-5",
            "optional": false,
            "field": "rating",
            "description": "<p>Overall rating of user.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date of user profile creation.</p>"
          },
          {
            "group": "200 JSON Response",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>last ISO8601 date when this profile was updated.</p>"
          }
        ]
      }
    }
  }
] });
