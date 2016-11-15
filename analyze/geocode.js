
var MongoClient = require("mongodb").MongoClient;

MongoClient.connect("mongodb://localhost:27017/engl452", function(err, db) {
  console.log("err", err);
  console.log("Connected correctly to server");
 
  db.close();
});