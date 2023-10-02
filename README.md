# Expense Tracker

A shared expense tracking application implemennted in Go. The application allows users to create groups, add other users to the group and track their shared expenses. Users can add expenses to the group, specifying splits of the expense and settle any outstanding debts among the group members. 

## Project structure

### api

The api directory contains the code for the api. The api is written using the Gin web framework in GO.

### web

The web directory contains the code for the web application. It is written in React + Typescript + Vite. 


## Getting Started
### Prerequisites
Before you begin, ensure you have the following software installed:

- Go (v1.19 or higher)
- Node.js (v18 or higher):


### API 

**Configure Environment Variables**

In the api directory, create a .env file with the following content and adjust as needed

```
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_NAME=expense_tracker
ALLOWED_ORIGINS=http://localhost:3000 
DOMAIN=localhost
```

Note: In production environment (GO_ENV=production) environment variables are loaded straight from the shell not the .env file

**Compile And Run Directly (Development)**

To run compile and run api directly from source code run the following command in the root of the project.

```
make dev-api
```

**Build And Run Binary (Production)**

To compile the source code into binary and run the binary run the following command in the root of the project

```
make build-api 
./api/bin/main 
```

### Web App

**Install Dependencies**

In the web directory, install the web application's dependencies by running:

```
npm install
```

**Configure Environment Variables**

In the web directory, create a .env file with the following content and adjust as needed

```
VITE_API_BASE_URL=http://localhost:3000
```

**Running With Development Server**

Start development server on port 3000

```
make dev-web
```

**Building**

Build application into static bundle

```
make build-web
```

## Getting Started With Docker

This approach starts the api server, web application server and database server without any additional configurations. By default the api server will be running on port 3001 and the we application server will be running on port 3000.

**Prerequisites**

Before you begin, ensure you have the following software installed:

- Docker 

- Docker Compose (if not already installed with docker)

**Build Docker Images**

Building docker images for api and web application.
```
make docker-build
```

**Run Docker Containers**

```
make docker-run
```


**Stop Docker Containers**

```
make docker-stop
```

## Deploy To AWS

The application can be deploy to AWS using the provided terraform files. The terraform file will provision the infastrcuture on AWS Cloud as shown in the diagram below. 

![AWS Clound](https://drive.google.com/uc?id=1z1jN-WEle6dqZlIwI32CvfPUMccdXLNj)

**Prerequisites**
- The Terraform CLI (1.5.0+) installed.
- The AWS CLI installed.
- AWS account and associated credentials that allow you to create resources.

**Setting AWS account credentials**

By default terraform will use the credeintails associated with default profile of the AWS CLI.

To explicitly specify the IAM role to be used to authenticate Terraform AWS provider set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables. 

```
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=..
```


**Initilaise Terraform**

In the root of the project run the following to initilaise Terraform.

```
terraform init
```

**Configure variables**

There are serveral Terraform variables that have to be configured before privisioning the infrasture. The variables that must be set are anotated with *. There variables can be assigned by a number of ways as outlined in the [Terraform Documentation](https://developer.hashicorp.com/terraform/language/values/variables#assigning-values-to-root-module-variables)

- `aws_region`: The AWS region where your resources will be created, default to us-east-1.

- `resource_prefix*`: A prefix that will be added to the names of all the resources created by this Terraform. 

- `common_tags`: A set of common tags that will be applied to all resources created by this configuration.

- `vpc_cidr_block`: The CIDR block for the Virtual Private Cloud (VPC) that will be created, default to 10.0.0.0/16.

- `api_instance_type`: The instance type for an API instance, default to t2.micro.

- `database_credentials`: Represents the database credentials (username and password) for accessing the database. The default credentials must not be used in a production environment. 

- `database_name`: The name of the database that will be created, default to shared_expense_tracker.

- `jwt_secret`: Represents the secret key used for JWT (JSON Web Token) authentication, default to secret. This must be changed in a production environment. 

- `key_pair_public_key_path*`: The path to the public key of an SSH key pair. 

- `api_docker_image*`: Represents the URL of the Docker image for the API web service. (Obtained from the next section)

- `api_autoscale_settings`: Settings for the API EC2 instance auto scaling group. 

**Building API docker image**

Provisioning the EC2 instance required the url to the docker image that will be deploy to it. To obtain the docker image URL, build the docker image with the target flatform linux/amd64 and push the built image to a container repostiory. 

Build docker image and tag your image according to our repository url. 
```
cd api
docker build . --platform linux/amd64 -t <repository url>    
```

Login to repository via docker CLI. 
```
docker login -u <username>
```

Push the built image to the pository. 
```
docker push <repository url>
```

**Provision Infrastructure**

To provision the infrastructure run the allowing command and type yes when prompted. This assume that have set your variables in a way that terraform will automatically load your variables. 
```
terraform apply
```
The command will return the web s3 bucket endpoint, `web_endpoint`, and the public DNS of the elastic IP associated with the EC2 instance, `api_endpoint`.

**Build And Upload the Web Application**

To build the application be must supply the `VITE_API_BASE_URL`, which is the url that points to the web api service EC2 instance. 

Run the following command to build the web application static assets
```
cd web
export VITE_API_BASE_URL=http://<api_endpoint>/api
npm run build
```
The build will produce the static assets in the `web/dist` folder 

To build upload the build asset to to s3 bucket run the following command (Assuming that you still in the `web` directory)

```
aws s3 cp ./dist/ s3://<bucket name> --recursive  
```

The application should now be accessible from `web_endpoint`, the ouput from the `terraform apply` command. 


**Cleaning Up**

To clean up the infrastructure run. 
```
terraform destroy
```

      

