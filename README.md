availability-service
================
The availability service REST API provides the capability to provide a product identifier (SKU) and zipcode to determine if the item is available for purchase in the location provided.  The service will also leverage the API key to determine the catalog identifier which is leveraged to associate a product catalog with the customer.


## REST API

- All endpoints are served over HTTPS
- Responses are provided as JSON.
- Request is authenticated and verified via an API key and AWS V4 HMAC signature verification.  

### V1

#### Request

    https://[endpoint]/v1/availability/[sku]/[zipcode]/[count]

    - sku - the product sku in string format
    - zipcode - the location of the delviery requirement.
    - count - the requested quantity. optional, defaults to 1
    - catalogId - will be derived from provided API Key


#### Response
##### Data
```
    HTTP/1.1 200 OK
    Date: Mon, 01 Jan 2015 17:27:06 GMT
    Status: 200 OK
    {
      "sku": "XXXX",
      "zipcode": "40245"
      "availability": true/false
    }
```
##### Response Codes
```
200 - Success
400 - Bad Request - invalid inputs
401 - Not a valid access token
404 - No Products found for given input
500 - Internal Server Error
```

### V2

#### Request

    https://[endpoint]/v2/availability/[sku]/[zipcode]/[count]

    - sku - the product sku in string format
    - zipcode - the location of the delivery requirement.
    - count - the requested quantity. optional, defaults to 1
    - catalogId - will be derived from provided API Key


#### Response
##### Data
```
    HTTP/1.1 200 OK
    Date: Mon, 01 Jan 2015 17:27:06 GMT
    Status: 200 OK
  {
      "sku": "XXXX",
      "zipcode": "40245”,
      "availability": true/false,
      “count”: 1,
      “intransit_qty: 0 ,
      “intransit_date”: “2016/10/12”,
      “intransit” : true/false,
      “quantity” : 50
}
```
##### Response Codes
```
200 - Success
400 - Bad Request - invalid inputs
401 - Not a valid access token
404 - No Products found for given input
500 - Internal Server Error
```

#### Project Structure
* cfn - Location of cloud formation templates leveraged to construct the AWS Environment
* datamodel - JSON document describing the service data model.
* availability - Availability services for talking to PDAV and AWS for retrieving availability data.
* zipcode - Zipcode services for retrieving Zip code data from AWS
* Server.go - Gin Server powering API.
* pdav.json - JSON file representing the buffer for the PDAV call.
* pdavhd.json - JSON file representing the buffer for the PDAVHD call.
* Dockerfile-compiler - Docker file for compiling the go binary.
* Dockerfile - Docker container implementation.

===================

#References &amp; Tools

###[Swagger](https://github.appl.ge.com/geappliancesales/hd-swagger-hosted)
We have put together some swagger documentation that provides a developer friendly API documentation and Client to execute requests.

###[Hmaid](https://github.appl.ge.com/geappliancesales/hma-aid)
Hmaid is a service that was built to aid in the development of hmac clients. We often struggled with developing a new client and having it work the first time. Trouble-shooting was a burdensome process. This service allows the client to send it an hmac request and it will respond back to the client with the data/steps it signed the request against. This will allow the client to compare the server signing process vs their own and see where the difference lies. This service will be enabled on an as-needed basis for trouble-shooting

# [Jenkins](https://hd-jenkins.al.ge.com/)

Jenkins is used as our continuous integration/deployment pipeline. When a change is pushed to github, a github webhook fires off a build job to Jenkins. Jenkins handles running the tests, building the docker image and refreshing the AWS instances if necessary.

### Availability Service Jobs

#### build-availability-service
This job builds the docker image and publishes to the registry in both the east and west regions.  This job will run on a code commit or pull request merge to the master branch.  When completed an image will exist in both east and west registries tagged with the job number build as well as dev-latest.  It also moves the cfn associated with the application to the appropriate s3 directory.  This location is leveraged when running the update-or-create-cfn job.

#### update-or-create-cfn
Based on appropriate parameters this job will leverage the CFN updated per the specified environment and location in s3 and create/update the application stack.  This job does not update the application unless it is creating a fresh environment or there was a configuration change within the scope of the AWS stack.  To execute a code deployment see the restart-instances job below.

#### restart-instances
This job based on parameters will query the ASG to determine the appropriate ec2 instances associated with the group.  It will then cycle through each instance and run a script via ssh.  This script will clean the local docker image registry and re-download and deploy a version of the application.  You must specify the appropriate tag.

#### promote-images-to-prod
Will promote a tag to prod-latest.  It is recommended to deploy the prod latest tag to production environmentes to ensure the tag is indeed a production quality deployment.  You can specifiy any tag on the creation job though following guide lines for image management is recommended

#### add-cloudflare-dns
After a stack is created you will need to associate the elb endpoint with cloudflare for wedeliver.io.  This job takes a parameter to specify the elb endpoint as well as the environment.  The result will be <env>.wedeliver.io/v1/<service> pointing to the appropriate application.  For production to support cross region latency based load balancing a different process is followed.  

### AWS Console

Account #: 7080-6217-3806

The console is where all of the services can be viewed and interacted with.

### Authentication
Each of our services implement the [AWS4 hmac](http://docs.aws.amazon.com/general/latest/gr/signing_aws_api_requests.html) signing specification. Our services use an [auth-gem](https://github.appl.ge.com/geappliancesales/auth-gem) that handles the server side implementation of this. Clients will be provided an `access_key` and `secret_key` and will sign their requests with their secret per the aws4 spec and send alogn their access key with each request.

### Code Samples
We have created a few [sample clients](https://github.appl.ge.com/geappliancesales/service-clients) to help customers implement hmac signing

#Availability Service - Operations Guide

This section describes the components associated with supporting the Availability Service including software development processes, continuos integration and deployment, architecture, and maintenance.

===================

## Local Development Environment Setup

The availability order service leverages a GOLANG leveraging the [Grape API Framework](https://gin-gonic.github.io/gin/).

It is recommended that you install [GOLAN](https://golang.org/doc/install).  Follow the installation instructions referenced on the site.

Now pull the source code from Github.  The link is referenced above.

```
git clone <repo name>
```

Navigate to the cloned repository.

```
go run server.go
```


Your service should now be running on localhost port 8888.  You can test by leveraging one of the sample [client implementations](https://github.appl.ge.com/geappliancesales/service-clients) of AWS V4 HMAC Signing.

===================

## Committing Changes

Once you have tested and validated your changes as well as updated the [RSpec](http://rspec.info/) test if necessary to accommodate your changes it is time to commit.

Add the modified files to the local git repository.  It is not recommended to perform a add *.  This leads to unnecessary files imported with the project source.

```
git add <file_name> ...
```

Commit the changes with an appropriate message to the local repository.

```
git commit -m 'appropriate message'
```

Push the changes to the remote repository.  This will automatically launch a Jenkins job to deploy the changes to development.  See the details of the Jenkins processes in the section above.

```
git push
```

## Dependent Batch Processes

View status and logs of batch processes leveraging Mesos and Chronos
[Chronos] (http://hd-mesos.al.ge.com:4400/)
[Mesos](http://hd-mesos.al.ge.com:5050/)

#### Availability Service Batch
[Source Code](https://github.appl.ge.com/geappliancesales/availability-service-batch)
The availability batch process updates two DynamoDb tables.  It first pulls the list of active Zipcodes based on the customers identifier.  It updates that list appropriately.  It then pulls the catalog from Amazon S3 based on a B2B Content Delivery Feed (See Services Batch for details).  Based on the catalog sku list it leverages a file transferred from the Bull Exit to pull inventory at each distribution center.  It updates the DynamoDb table appropriately based on this information.  

#### Services Service Batch
[Source Code](https://github.appl.ge.com/geappliancesales/services-service-batch)
The availability service also has a dependency on the services service batch processes.  The services service will pull the B2B Content Delivery feed for all configured customers.  Based on this feed it will do its thing and then push the file to Amazon S3.  This file in S3 is what the availability batch is dependent on to get its catalog for the customer.

#### Local Development Environment Setup
Ensure you have installed the J2EE version of [Eclipse](https://eclipse.org/downloads/).  Navigate to the Git perspective and clone the repository.  Create a generic project from the source in the "Working Directory".  Navigate back to the Java EE perspective and change the configuration of the project to Maven.  Each Batch process has a Driver class located in the package com.gea.batch.  Review the class to determine the arguments required to run.

####Batch Jenkins Jobs
Each batch process has two Jenkins job associated with it.  A build job and a deploy job.  The build job will build the docker image and deploy a development version of the job when changes are pushed to Github.  The deploy change can be leveraged to deploy additional versions of the batch to support other environments.  Since batch processes are mainly associated with the DynamoDb environment you should align the deployments there - not necessarily with the service environments.

For Example:
##### build-homedelivery-availability-service-batch
##### deploy-homedelivery-availability-service-batch

You will notice the deploy job expects a version to be passed.  This version will need to be the latest image pushed to the docker registry.  This value is usually the latest build number of the build job.

## Dependent DynamoDb Tables
### homedelivery.api-keys.<env>
Leverages this table to pull customer information as well as the keys associated with authentication and authorization

### homedelivery.zipcode-service.<env>
Based on the customer information and zipcode passed it will pull the distribution center code (ADC) and validate that delivery is available for the zipcode.

### homedelivery.availability-service.<env>
Based on the catalog pulled from the customer information as well as the sku and ADC retrieved from the zipcode query it will validate the inventory level is greater than zero.

You can view the data in these tables by leveraging the AWS DynamoDb console.

## Cloud Watch Logs
To view the logs for the service navigate to the AWS console and select CloudWatch.  In the CloudWatch console in the left navigation select "Logs".  Now select "HomeDelivery-<env>" based on the environment needed.  Now select the service you wish to review the logs.  

## Notifications
Notifications are sent through Amazon Simple Notification Service.  
