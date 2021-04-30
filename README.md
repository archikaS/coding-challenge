# coding-challenge

# ASSUPMTIONS
1) It is assumed that Endpoint is hit by authenticated user.
2) Pagination is not applied for this code assuming limited data set.
3) Involvement of less dependency in project. For query params validation, a function is used instead of middleware as it would have required using context package for passing variables. 

# Project Structure
a) Router file - the uri for api endpoint is specified in this file.
b) Controller - to do RESTful things for a particular type of data or   resource.
c) Service - the endpoint for an api where all the logics and operation is performed.