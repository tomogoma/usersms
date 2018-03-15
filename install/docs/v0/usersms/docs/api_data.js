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
            "type": "Integer",
            "optional": true,
            "field": "offset",
            "defaultValue": "0",
            "description": "<p>Index from which to fetch ratings (inclusive).</p>"
          },
          {
            "group": "URL Query",
            "type": "Integer",
            "optional": true,
            "field": "count",
            "defaultValue": "10",
            "description": "<p>Number of ratings to fetch.</p>"
          },
          {
            "group": "URL Query",
            "type": "String",
            "optional": true,
            "field": "byUserID",
            "description": "<p>Filter ratings by rater's userID. At least one of forUserID or byUserID must be provided.</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200": [
          {
            "group": "200",
            "type": "JSON",
            "optional": false,
            "field": "body",
            "description": "<p>Array of <a href=\"#api-Objects-Rating\">ratings</a>.</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
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
    "success": {
      "fields": {
        "200": [
          {
            "group": "200",
            "type": "JSON",
            "optional": false,
            "field": "body",
            "description": "<p>Updated <a href=\"#api-Objects-User\">user</a> object</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  },
  {
    "type": "POST",
    "url": "/ratings/users/{userID}",
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
            "field": "userID",
            "description": "<p>ID of the user to rate (ratee).</p>"
          }
        ],
        "JSON Body": [
          {
            "group": "JSON Body",
            "type": "String",
            "optional": false,
            "field": "byUserID",
            "description": "<p>ID of the user rating (rater).</p>"
          },
          {
            "group": "JSON Body",
            "type": "Integer",
            "size": "0-5",
            "optional": false,
            "field": "rating",
            "description": "<p>The rating awarded by rater to ratee.</p>"
          },
          {
            "group": "JSON Body",
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
        "200": [
          {
            "group": "200",
            "optional": false,
            "field": "200",
            "description": ""
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
        "JSON Body": [
          {
            "group": "JSON Body",
            "type": "String",
            "optional": true,
            "field": "name",
            "description": "<p>New Name</p>"
          },
          {
            "group": "JSON Body",
            "type": "String",
            "optional": true,
            "field": "ICEPhone",
            "description": "<p>New (In Case of Emergency) phone number</p>"
          },
          {
            "group": "JSON Body",
            "type": "String",
            "allowedValues": [
              "\"MALE\"",
              "\"FEMALE\"",
              "\"OTHER\""
            ],
            "optional": true,
            "field": "gender",
            "description": "<p>New gender</p>"
          },
          {
            "group": "JSON Body",
            "type": "Object",
            "optional": true,
            "field": "avatarURL",
            "description": "<p>New profile picture URL</p>"
          },
          {
            "group": "JSON Body",
            "type": "Object",
            "optional": true,
            "field": "bio",
            "description": "<p>New brief description of user</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200": [
          {
            "group": "200",
            "type": "JSON",
            "optional": false,
            "field": "body",
            "description": "<p>Updated <a href=\"#api-Objects-User\">user</a> object</p>"
          }
        ]
      }
    },
    "filename": "pkg/handler/http/handler.go",
    "groupTitle": "Service"
  }
] });