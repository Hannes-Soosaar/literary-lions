# literary-lions

a web forum that allows users to communicate, associate categories with posts, like/dislike posts & comments, and filter posts.

Login

Landing page--

Has user
Has user roles 
Has user posts
Has book
Has Post
Has Thread
Has Password
Has 

Necessities:
go get golang.org/x/crypto/bcrypt

~/your_project_folder$ docker build -t lions:1.0 .

docker run -p 8080:8080 lions:1.0

command to show all dangling images:
docker images -f "dangling=true"
to remove any dangling images  run:
docker system prune