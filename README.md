# literary-lions

a web forum that allows users to communicate, associate categories with posts, like/dislike posts & comments, and filter posts.

## basic operations

-logon from the top right profile profile icon 
-the user 

Login

Landing page--

#Technical details

Necessities:


# Running app from a Docker container

~/your_project_folder$ docker build -t lions:1.0 .

docker run -p 8080:8080 lions:1.0

command to show all dangling images:
docker images -f "dangling=true"
to remove any dangling images  run:
docker system prune