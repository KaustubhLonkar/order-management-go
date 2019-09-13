# order-management-go

# Problem statement

Implement an order management system using micro service architecture and deploy minikube environment using docker and kubernetes tools.
Add products to Inventory using REST API
Take Order Based Product available in Inventory and Generate invoice using REST API
Post Order for Shipping using message broker like Kafka
Non-Function Requirements:
1. Documentation using Swagger
2. Performance Metrics


The code to be compiled and deployable in minikube.

# The application is running locally.Below are the steps to run the application locally:-
    Open a command prompt/shell.
    Navigate to server.
    Run the command 'go run server.go'
    The command prompt shall display the logs which show logging details and the URLS supported.These content on command prompt would keep on growing as long as you send requests.
    In order to send the request uses the forms in view.
    The below lets you add products AND place order respectively.

    https://localhost:8888/
    https://localhost:8888/orderForm.html

# Created REST APIs to add products to Inventory and take Order Based Product available in Inventory and below URLS are supported:-

    https://localhost:8888/addProduct
    https://localhost:8888/getProducts
    https://localhost:8888/placeOrder
    https://localhost:8888/metrics (Prometheus)

# Mysql Database is used and the schema is added which can be found in mysqldump.
# Added view to add/retrieve data so that you can initiate rest calls through view.
# Used Apache kafka so that the place order API could commit the order details in database and when successful will post shipping/order details in Kafka.
# I have added a prometheus dependency which shall capture the logs "https://localhost:8888/metrics"
# DockerFile creates the image of application.
# docker-compose.yaml contains the configs/links related to REST API image,zoo keeper,kafka and mysql

# As far as deployment is concerned,it could be done in 2 ways.
    Specify your docker image in deployment.yaml and that deploys the image to minikube.Else if you do not have a dockerfile and prefer to use docker-compose.yaml file then you would have to use Kompose for deployimg to kubernetes.Just in case if kompose is being used to deploy in kubernetes (Docker Compose + Kubernetes) we have to follow the below steps:-

    Install Compose:- http://kompose.io/installation/
    Create a directory in your filesystem and name it kompose (for ease of identification)
    $ cd kompose
    For Linux :
    Download Kompose:
    $curl -L https://github.com/kubernetes/kompose/releases/download/v1.16.0/kompose-linux-amd64 -o kompose
    After downloading

    $chmod +x kompose
    $sudo mv ./kompose /usr/local/bin/kompose

    Start minikube
    $ minikube start

    get your docker-compose.yaml file



    Now converting this Kompose file we can get all required deployment and service yaml(s)from this.
    $ kompose convert                           
    INFO Kubernetes file "docker-compose.yaml" created  

    now do a ls to see all the yaml(s) are created or not.
    $ ls -l

    Now to create the deployments and services in kubernetes we can Kompose up
    $ kompose up

    Now check whether your deployments are ready or not :
    run :
    $ kubectl get deployment,svc,pods,pvc

    check minikube dashboard
    $ minikube dashboard

    After the deployments and image pull is complete check for the url to view the output
    $ minikube service frontend — url
