### Hexa arranger is a wrapper and helper for Uber Cadence

### Install
```bash
go get github.com/kamva/hexa-arranger
```


##### Notes about error handling:
- Read [the docs](https://github.com/temporalio/sdk-go/blob/master/temporal/error.go) about error types in temporal.

- solution to handle errors:  
    temporal don't have interceptor or activities yet.
    For workflows, it has, maybe later use it.
    for now should report and convert errors when return 
    error from either activity or workflow.
    I think a simple decorator for each workflow
    and activity can be a good idea. when temporal
    implemented interceptor for activities, we 
    can use the workflow and activity interceptors
    to implement it without the decorator.

error behaviour:

```
-> when activity handler returns error:
    - report error
    - convert hexa error to application error

-> when workflow get activity|child_workflow error:
    - nothing. it can convert it to hexa error if needed.

-> when workflow itself returns error:
    - report error 
    - convert hexa error to application error

-> when app get workflow error:
    -> convert error to hexa (simply error by call to a method on returned error).
    -> report error (don't need to implement, our app report every error in the edge)
    
-> returned error from workflow to app:
    
```


#### TODO
- [ ] Provide OpenTracing.
- [ ] Write tests.
