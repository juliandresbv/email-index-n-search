# Email Index 'n' Search

## Table of Contents

* Description
* Pre-requisites
* Quick Start
* Considerations
* Documentation

## Description

This project contains the source code for the project of "email-index-n-search". A set of apps that allow you to index and search emails persisted on ZincSearch ("a search engine that does full text indexing").

## Pre-requisites

Before running the application, you need to have the following dependencies installed on your system:

* Docker: https://www.docker.com/
* ZincSearch (Docker Image): https://zincsearch-docs.zinc.dev/installation/
* Go (>= v1.21.1): https://go.dev/
* NodeJS (>= v18): https://nodejs.org

## Quick Start

### Initial setup

Follow the steps below to quickly set up the apps:

0. Clone the repository:

```bash
git clone https://github.com/juliandresbv/email-index-n-search
cd email-index-n-search/
```

1. On the root directory of the repository, run ZincSearch on Docker:

```bash
docker-compose up -d
```

or

```bash
docker compose up -d
```

2. Open three terminal windows/tabs at the root of the project in each one:
    
    1. On tab/window #1, run the following commands to install the dependencies for `indexer` app:

        ```bash
        cd indexer/
        go mod download -x
        ```

    2. On tab/window #2, run the following commands to install the dependencies for `server` app:

        ```bash
        cd server/
        go mod download -x
        ```

    3. On tab/window #3, run the following commands to install the dependencies for `emails-search-app` app:

        ```bash
        cd emails-search-app/
        npm install
        ```

3. Set the environment variables on a .env file for every app following the form of the .env.example file on each app's directory.

### Running the apps

Follow the steps below to quickly set up the apps:

1. Indexer:

    0. About the download and decompression dataset:

        The `indexer` app downloads and decompresses the dataset automatically. However, if you want to do it manually, you can do so by downloadloading the dataset from the following link: https://www.cs.cmu.edu/~./enron/enron_mail_20110402.tgz. Then, placing it on the `indexer/data` and decompressing it.

    1. Open/re-use a terminal window/tab at the `indexer/` directory and follow the steps below:

        1. Run app:

            ```bash
            go run main.go
            ```

        > Note: To run the app performing profiling, add the `-profile.mode=<mode>` flag to the command above. The available modes are: `cpu`, `mem`, `goroutine`, `thread`. The profiling results will be saved on the `indexer/profiling-results` directory.


2. Server:

    1. Open/re-use a terminal window/tab at the `indexer/` directory and follow the steps below:

        1. Run app:

            ```bash
            go run main.go
            ```

3. Emails Search App:

    1. Open/re-use a terminal window/tab at the `emails-search-app/` directory and follow the steps below:

        1. Run app:

            ```bash
            npm run dev
            ```

## Considerations

### Profiling

// TODO: Add profiling results, and images

### Optimizations

// TODO: Add content and sections 

#### Sequential vs Concurrent

#### JSON file size

#### DTOs (Marshaling and Unmarshaling) vs Bytes

## Documentation

### API documentation

// TODO: Specify API documentation via Swagger
