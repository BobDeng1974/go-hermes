# go-hermes
Go-hermes is a Go app that exposes an HTTP API which receives requests with software metrics.

# Main components
1. HTTPS endpoint where all requests will go to (we need to spread the servers across the globe to achieve low latency).
2. Analyze data from requests. This depends on the type of data we received and each one will require different analysis. For example, if it's an app error analysis would be to extract stack trace.
3. Present data on a UI dashboard. This will allow users to get a visual understanding on what's under-performing and what's performing well.
4. Agents that push data to endpoint. These are going to be executable files installed on remote host (client installs on their machines), and will collect metrics and push them to our endpoint.
5. API Clients in different programming languages. This will allow users to create custom metrics that matters to them, and push them programmatically to our endpoint (for example profiling for a function in their app).

There are many interesting issues to tackle on a project like this:
- coming up with a distributed system architecture
- scalability and fault tolerance

The kind of payload we will work with, is going to be JSON as its more lightweight than XML, and all the incoming requests will be compressed, meaning less bandwidth usage (for customers and for us) and faster communication.

Kind of metrics that we can collect (below list is not exhaustive):
- app data (database, external requests, errors/exceptions)
- server monitoring (disks, memory, cpu)
- mobile data (usage data, errors, device information)

We also need to think about app events such as: deployments, or software updates and allow customer to compare how that change impacted the servers/application performance. That will allow them to easily decide if a rollback is needed.

# Try it out
* [Install Go](https://golang.org/dl/).
* Clone the repo: ```$ git clone git@github.com:go-hermes/go-hermes.git```
* [Install MongoDB](https://www.mongodb.com/).
* Create a `gohermes` database in MySQL:
```sql
CREATE DATABASE IF NOT EXISTS `gohermes` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `gohermes`;
```
* Create `server` table in MySQL:
```sql
CREATE TABLE `server` (
  `id` int(11) NOT NULL,
  `userId` int(11) NOT NULL,
  `hostname` varchar(40) NOT NULL,
  `os` varchar(20) NOT NULL,
  `lastMetricDate` datetime DEFAULT NULL,
  `creationDate` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


ALTER TABLE `server`
  ADD PRIMARY KEY (`id`),
  ADD KEY `userId` (`userId`),
  ADD KEY `userId_2` (`userId`,`hostname`);

ALTER TABLE `server` MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
```
# Build
* Copy `build.sh.dist` to `build.sh` and modify it to match your credentials.
* Run `$ ./build.sh`

### Note regarding requests
Endpoint runs on https and for development self-signed keys are used, and not used
in production anywhere (therefore `-k` flag in `curl` is used in examples below).
If you want to make this service publicly accessible, [please generate your own keys](https://github.com/golang/go/blob/master/src/crypto/tls/generate_cert.go).

* Open a new terminal and register a new user:
```
$ curl -k -H "Content-Type: application/json" -d '{"username":"myUser", "email": "user@example.com"}' https://localhost:8080/user/create
```

You should see something like:
```
{"message":"User created successfully!","Metadata":{"id":1,"username":"myUser","email":"user@example.com","creationDate":"2016-11-05T16:09:33.340621976Z"}}
```

Try sending data with invalid email, passing an id, send request with existing email/username.

## /server/create endpoint
Open another terminal and try creating a new server:
```
$ curl -k -H "Content-Type: application/json" -d '{"hostname":"usersetup", "user":{"id":1}, "os":{"name":"ubuntu"}}' https://localhost:8080/server/create
```

Try sending data with non-existing customer id, or try adding an existing server.
