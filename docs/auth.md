# Authentication

## Paths:

 - > ```/register  - POST```

 - - Gives a jwt token that its metadata contains tokenuuid, username, userid and expire

 - > ```/whoami    - GET```

 - - Returns username and userid. Bearer header should be added

 - > ```/login     - GET```

 - - TODO

 - > ```/todo      - POST```

 - - This is returns what you write in data. Should be post with bearer header. 

 - - It doesnt save any thing to db.