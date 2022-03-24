db = db.getSiblingDB('mongo_database');
db.createCollection('cats', {capped: false});
