# myeircode

Yet another learning exercise project 👨🏼‍🏫

## Outline: Description
A small web app that will return a predefined list of eircodes. The user should be able to interact with the web app to add and remove addresses in a _secure_ manner. 

## Key Features.
Features to be implement to provide functionality and learning oportunities. 

  ### Containerised:
  The app should be able to run in a container. 
  ### State:
  The only state should be the list of addresses, for the MVP this can be stored in a S3 bucket.
  ### WebUX
  Some way to pretty print the address list. 
  ### API
  The app should implement CRUD behavior
  ### Secure 1
  Read can be open, but write operations should be approved. Suggestion is a token exchange to a trusted e-mail account for approval.
  ### Secure 2
  The app should implement TLS using `LetsEncrypt`. - Bonus for implementing auto cert rotation. 

## Approach:

Will create a ~Pivotal~ VMware [tracker](https://www.pivotaltracker.com/n/projects/2478782) project 

