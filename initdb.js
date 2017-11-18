// assuming mongodb is running on localhost, default port 27017
db = db.getSiblingDB("marauders");
db.users.insertMany([
    {
        _id: ObjectId("5a00b64d652276486a496816"),
        uuid: "thisuserdoesexist",
        name: "jon doe",
        email: "jon@jon.com",
        photoURL: "sdf/sdf/sdf/",
        friends: ["13", "14", "12345"],
        location: {
            longitude: -140.25176239999999,
            latitude: 77.2616841
        }
    },
    {
        _id: ObjectId("5a02a0a6e7e282cd149ad4ef"),
        uuid: "shinnerninja",
        name: "Ingyu Shin",
        email: "steven.shin.95@gmail.com",
        photoURL: "www.google.ca",
        friends: ["oysterblue", "jonjonbinks"],
        location: {
            longitude: -123.25176239999999,
            latitude: 80.2616841
        }
    },
    {
        _id: ObjectId("5a02a0bbe7e282cd149ad4f0"),
        uuid: "jonjonbinks",
        name: "Jonathan Fleming",
        email: "jonathanfleming135@gmail.com",
        photoURL: "www.google.ca",
        friends: ["oysterblue", "shinnerninja"],
        location: {
            longitude: -129.25176239999999,
            latitude: -60.2616841
        }
    },
    {
        _id: ObjectId("5a02a0d0e7e282cd149ad4f1"),
        uuid: "oysterblue",
        name: "Justin Awrey",
        email: "awreyjustin@gmail.com",
        photoURL: "www.google.ca",
        friends: ["shinnerninja", "jonjonbinks"],
        location: {
            longitude: -130.25176239999999,
            latitude: 30.2616841
        }
    }
]);
