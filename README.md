# tsla lab
Exercise performed in GO language for demonstration purposes. The prject exercises the requirement's defined in an external document. This project will not reveal the requirements publically to protect the exercise criteria and the source of the exercise.

This project will however demonstrate a generic use of asynchronous go routines (agents) and a concurrent safe listener that is managed by a controller. The console will send commands to the controller to execute the control flow. 

## CLI options
------
<fill in here>

## Interactive console options
* Starting the agents 
* Stopping the agents, some or all
* List the current state of running agents and their types (good instrumentation candidate, think prometheus)

at anytime after the program starts, type "h" or "help" to get the commands menu

# Build Instructions
Instructions here assume you want the source code to evaluate, as well as run the program. 
Checkout the project from GitHub
```git clone github.com/dfense/gRPCtemplate```
Run the program
 ```go run github.com/dfense/tslab```
to build this project, just clone the project, and run  
```go build github.com/dfense/tslab```  

**need to build for alternate OS or ARM , just add the standard GO environment variables before you build**

* Note to Reviewers
Special thanks for being allowed to participate in the Code Challange and considerations for the available position on this team at TESLA Energy division. 

The Code Challenge was clever in many ways. It was brief and yet concise in the core intent of the test. It was also abstract and allowed a lot of design 'discussion' to have occurred had this been a normal project with multiple members. The part I would have allocated much more time to, would be upfront discussions of the problem we were trying to solve. This surely would have led to a more aligned perhaps simpler, perhaps more complex, more expandable implementation. 

My assumptions at what was best to demnonstrate on the code submitted, were multiple different design patterns and tools within the GO language. The use of channels are really the highlight of the code base. Separation of the different behaivor types, and consistent code format quality. 

I wanted to state, i took a chance in exercising many additional options above the standard requirement. Not to try and impress with clevernes, but make a more concrete choice of discussion topics around the review process (assumptions that we have the opportunity to discuss) My goal was still to meet all the original basic requirments, and expand.

The other approach that I struggled not choosing, was to keep it bare bones the simplest choice to still reach the goal. I apologize if i added additional work on reviewers. There is beauty in simplicity. I keep a printed poster in my cubicle at work. A quote from a French novelist in the early-mid 1900s. 

`Perfection is achieved, not when there is nothing more to add, but when there is nothing left to take away.`

With that stated, here are a few design improvments, variances, refactors that might be worthy of making notes around. 

## Design Discussions: 
* more interface abstraction. interfaces allow for great testability and expansion. Also in GO lang assist greatly in trapping yourself into the ultimate evil of cyclic import errors 
* test driven up front understanding intent. my minimal test just show the awareness they are critical to good code and maintainability
* separation of control vs observation in the listener (aggregator)
* use of channels were well demonstrated, but there are more choices to have discussed. The action of a steady stream of events flowing back to aggregator of channel, vs perhaps a different approach with a callback via a subscribe to thing instead (and more)
* use of richer monitoring software and TimeSeries database tools. Prometheus, InfuxDB, etc are excellent references for such tools. I began a few instrumentation points, but more can be identified
* stronger build script for running generators, vetters, consistency checking, etc. I have a favorite Makefile i continue improving upon that i find value with in my toolbox
* model the Actors. Ability to much better simulate behaivor on each and every things.Thing Actor
* suspicion the things.Thing structs would more realistically be network nodes, would allow for interesting uService design, and networking options. an original suspicion a gRPC clients/servers were of interest, but perhaps beyond the scope of this exercise
* more CLI options
* appending, overwriting, rolling event files
* code reuse, find common reusable func. 
* utilize Composition reuse, candidate is things.Thing implementations
* introduce an HTTPS service symetric with Console behaivor
* build tagging with git hash, date, version at compile time
* network logging - cloudwatch, syslog, use of log library extensions
* better naming conventions once problem was better understood. more GO idomatic naming
* ...

# Screenshots of console
```go run ... fill in rest after i have code written``` 