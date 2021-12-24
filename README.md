![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/enriquetejeda/github-webhook-docker-stop)
![GitHub top language](https://img.shields.io/github/languages/top/enriquetejeda/github-webhook-docker-stop)
![GitHub last commit](https://img.shields.io/github/last-commit/enriquetejeda/github-webhook-docker-stop)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Github Webhook Docker Stop

A solution for interact with docker api local with a simple web server for listen events from github webhook 

## How works?

The golang app interact with the docker api for search and stop the containers with a specific label (`projectName` & `pullRequestNumber`).

*This labeling is mandatory for the containers that you want to detain.*

## Requirements

* Docker Engine :heart:
* Go 1.17

## Getting Started

### Docker :heart:

You only run this command in your terminal:

```
docker run -p :8080 \
-v /var/run/docker.sock:/var/run/docker.sock \
-e GITHUB_CLIENT_SECRET=YOUR_WEBHOOK_SECRET \
ghcr.io/enriquetejeda/github-webhook-docker-stop:latest
```

### Standalone

1. Rename the `.env.example` to `.env` and configure the values
2. Compile with the command `make build`
3. Run the command `make run`

## Development

### Building the binary

I provided a makefile for do this job, only run this command:
```
make build 
```
### Building the container

I provided a makefile for do this job, only run this command:
```
make build-docker
```
### Environment Variables 

| Name  | Description  | Default | Required |
| -- | -- | -- | -- |
| HOST | The address for init the webhook webserver  | `0.0.0.0` | *no* |
| PORT | The port for the webserver | `8080` | *no* |
| DOCKER_HOST | The docker api endpoint | `unix:///var/run/docker.sock` | *no* |
| GITHUB_CLIENT_SECRET | If you configure a secret for the webhook put here | `123456` | *no* |

## How contribute? :rocket:

Please feel free to contribute to this project, please fork the repository and make a pull request!. :heart:

## Share the Love :heart:

Like this project? Please give it a ★ on [this GitHub](https://github.com/EnriqueTejeda/github-webhook-docker-stop)! (it helps me a lot).

## License

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 

See [LICENSE](LICENSE) for full details.

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.

