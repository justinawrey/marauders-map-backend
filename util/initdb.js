// assuming mongodb is running on localhost, default port 27017
db = db.getSiblingDB("marauders");

// put 60 random users around mcleod building
var mcloed_long = -123.249629;
var mcloed_lat = 49.261895;
var mass_entry = [];

for (var i = 0; i < 45; i++) {
    mass_entry.push({
        uuid: `uuiduser${i}`,
        name: `user${i}`,
        email: `user${i}@test.com`,
        photoURL: "sdf/sdf/sdf/",
        friends: [`user${(i + 1) % 60}`, `user${(i + 2) % 60}`, `user${(i + 3) % 60}`],
        location: {
            longitude: mcloed_long + (Math.random() / 100.0) - 0.005,
            latitude: mcloed_lat + (Math.random() / 100.0) - 0.005
        }
    })
}

db.users.insertMany(mass_entry);
db.users.insertMany([
    {
        uuid: "thisuserdoesexist",
        name: "jon doe",
        email: "jon@jon.com",
        photoURL: "sdf/sdf/sdf/",
        friends: ["oysterblue", "jonjonbinks", "shinnerninja"],
        location: {
            longitude: -123.249629,
            latitude: 49.261895
        }
    },
    {
        uuid: "shinnerninja",
        name: "Ingyu Shin",
        email: "steven.shin.95@gmail.com",
        photoURL: "www.google.ca",
        friends: ["oysterblue", "jonjonbinks"],
        location: {
            longitude: -123.249629,
            latitude: 49.261895
        }
    },
    {
        uuid: "jonjonbinks",
        name: "Jonathan Fleming",
        email: "jonathanfleming135@gmail.com",
        photoURL: "www.google.ca",
        friends: ["oysterblue", "shinnerninja"],
        location: {
            longitude: -123.249629,
            latitude: 49.261895
        }
    },
    {
        uuid: "oysterblue",
        name: "Justin Awrey",
        email: "awreyjustin@gmail.com",
        photoURL: "www.google.ca",
        friends: ["shinnerninja", "jonjonbinks"],
        location: {
            longitude: -123.249629,
            latitude: 49.261895
        }
    }
]);