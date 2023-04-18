# Forum

This is a project for creating a web forum that allows communication between users, associating categories to posts, liking and disliking posts and comments, and filtering posts. The project will use SQLite as the database library for storing user data, posts, comments, etc. This readme file provides instructions for running the project and explains the features of the web forum.

## Main features
- User authentication
- Communication
- Likes and Dislikes
- Filtering

## Installation

Use the package manager [pip](https://pip.pypa.io/en/stable/) to install foobar.

Follow the instructions below to run the project:

1. Clone the repository to your local machine.
2. Make sure you have [Docker](https://docs.docker.com/get-docker/) installed on your machine.
3. Open the terminal and navigate to the project directory.
4. Run dockerize.sh from /sh_scripts/ to build image and run container.

```bash
bash sh_scripts/dockerize.sh 
```

5. Open a web browser and go to http://localhost:4000/ to access the forum.
6. Run dockerize_clear.sh from /sh_scripts/ to prune.

```bash
bash sh_scripts/dockerize_clear.sh
```
 

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Authors 
* [oyarukhi](https://01.kood.tech/git/oyarukhi)
* [dyskol](https://01.kood.tech/git/dyskol)
* [jegor_petsorin](https://01.kood.tech/git/jegor_petsorin)
* [erikje](https://01.kood.tech/git/erikje)

## License

[MIT](https://choosealicense.com/licenses/mit/)