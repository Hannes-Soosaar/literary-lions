# literary-lions app

A web forum that allows users to communicate, associate categories with posts, like/dislike posts & comments, filter and search for posts.

## basic operations

    - clicking on the logo at the top left of the nav bar will take you to all categories
    - logon from the top right profile profile icon.
    - the user can like/dislike a post and comments once.
    - the user can not like and dislike a post/comment at the same time
    - the user can undo the like/dislike a post/comment by clicking on the post/comment again
    - after the user changing profile data they will be logged out.

# Technical details

    -Developed with Sqlite3 and Golang V1.18
    -For dependencies please check the go.mod file

## Starting the app in Dev mode

  To run the app in dev mode execute the following commands in the terminal.

    1. Initialize the Lions database

        $ make liondb

    2. Start the server

        $ make run 

# Running the built app from a Docker container

To running the app from a Docker container on linux (ubuntu).

  1. install the latest version of Docker from the developer https://docs.docker.com/desktop/install/linux-install/

  2. create the image by entering to the terminal in the folder where the Docker file named "Dockerfile" is:
    
        $ docker build -t lions:1.0 .

  3. to run the created image on port 8080 run the following command in the terminal
    
     $  docker run -p 8080:8080 lions:1.0

The app is now running and can be accessed with your browser navigating to localhost:8080

# Credits and Bugs

The app was create by Hannes Soosaar and  Karl - Hendrik Kahn as a 2024 Kood/Jõhvi Golang project if you have any question or find any issues/bugs with the literary-lions app please write an email to
hsoosaar@gmail.com or karlhendrikkahn@gmail.com.