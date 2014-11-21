go-question-and-answer-ibm-watson-on-bluemix
=================================

Go Question and Answer using IBM Watson on Bluemix

## DEMO

[Go Question and Answer Watson on Bluemix](http://go-qa-watson-example.mybluemix.net/)

## General Requirements

* [Go](http://golang.org/)
* [Godep](https://github.com/tools/godep)
* [Martini](http://martini.codegangsta.io/)
* [IBM Bluemix](https://bluemix.net/)

## IBM Bluemix
Create your IBM Bluemix account and create a new app with custom buildpack. Once everything set up, a hub.jazz.net account should be created using Bluemix App Dashboard.

## Create App with CF
It's possible to create a new Bluemix app via CF with the following command:
  fvitullo/go-qa-example$ cf push <yourappname> -c https://github.com/cloudfoundry/go-buildpack.git

### Pushing code and Running

* Remember to edit Procfile, .godir and Manifest.yml with the proper names
* Run the following command 
  fvitullo/go-qa-example$ go install 
  fvitullo/go-qa-example$ godep save

Once a hub.jazz.net repository has been created and connected, it's possible to push the code in and go through the staging process automatically.

It's possible to look at the app logs using the following command:
  fvitullo/go-qa-example$ cf logs <app-name> --recent


