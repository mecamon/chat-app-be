//the creation of the user for development purposes
db.createUser(
  {
    user: "developer",
    pwd: "example",
    roles: [
    	{
        role: "readWrite",
        db: "chat"
    	}
    ]
  }
);