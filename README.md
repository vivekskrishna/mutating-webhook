# Extending Kubernetes — Part 2 — Mutating Webhook

![enter image description here](https://raw.githubusercontent.com/krvarma/mutating-webhook/master/images/extending_kubernetes_mutating_webhook.png?token=AA46XG3KLIZ2ROIABPDJMDK7UEGXE)

Extending Kubernetes is a series of articles that explore different options available to extend the Kubernetes system's functionality. The series discusses some of the methods to extend the functionality of Kubernetes.

In this part 2 of the series, we will discuss how to develop a Mutating Webhook.

# Admission Controllers

Kubernetes Admission Controllers are component that intercepts API requests and can accept or reject the API requests. Admission controllers can modify, approve, or deny the incoming API requests. There are many admission controllers are there in the Kubernetes system. Two of them are of particular interest to us, the mutating webhook admission controller and the validating webhook admission controller.

# What are the admission webhooks?

Admission Webhooks are HTTP callbacks that Admission Controllers calls when there is an API request. Admission Webhook returns responses to the API Requests. There are two types of Admission Webhooks in Kubernetes, Mutating Admission Webhooks and Validating Admission Webhooks.

The Admission Controller calls mutating webhooks while in the mutating phase. Mutating Webhooks can modify the incoming objects. Mutating Webhooks can do it by sending a patch in the response. Examples of mutating webhooks are, adding additional labels and annotations, injecting sidecar containers, etc.

Validating Admission Webhooks called in the validating phase and can only accept or reject a request. It cannot modify the object. Examples of validating webhooks allow access to only authorized namespaces, allowing/denying the incoming API requests based on corporate policy, etc.

Here is a diagram representing the admission controller phases.

![Kubernetes API Request](https://raw.githubusercontent.com/krvarma/mutating-webhook/master/images/Kubernetes%20API%20Request.jpeg?token=AA46XG2E5VY2MYRT7UNHKAC7R26GC)

In this article, we will explore mutating webhooks. In the next part, we will explore validating webhook.

I searched a lot on Google about the Mutating Webhooks; several resources explain how to create mutating webhooks. One particular example that I consider is this [sample code](https://github.com/kubernetes-sigs/controller-runtime/tree/master/examples/builtins). The sample is from the controller-runtime repository. I took the [mutatingwebhook.go](https://github.com/kubernetes-sigs/controller-runtime/blob/master/examples/builtins/mutatingwebhook.go) and [main.go](https://github.com/kubernetes-sigs/controller-runtime/blob/master/examples/builtins/main.go) as a starting point.

This article will create a webhook that inserts a sidecar container that handles the application logs. We all know the application log is an integral part of the application development. Many specialized tools are there that aggregate logs and send them for processing. We choose one tool and may decide to change it at a later stage. To make it easy to replace the logging tool at any stage, we usually write a wrapper around the logging module and hide all the log complexity. In this article, we will develop a simple module that handles application logs. We will expose an HTTP API that will accept a string and log to the console. In an actual implementation, we may use any production-grade tools like FluentBit, Open Telemetry, etc. However, for the sake of simplicity, we will print it to the console.

The mutating webhook that we will develop will automatically inject the logging server container into our Kubernetes application. It will enable the application developer to concentrate on the logic and connect to the logging module and log application events.

As explained earlier, the Admission Controller allows Mutating Webhook to modify the incoming API request. We can insert a container spec in the incoming API request. The API request will finally become a multi-container pod.

# This project is an example to change a service Type from LoadBalancer to ClusterIP based on presence of a label in the service definition.