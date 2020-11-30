### How to run examples:
- Run `docker-compose up` to start cadence service
- rnu your example with mode "worker"
- run your example with mode "trigger

__Notes:__
- If you want run commands on cadence cli, go to the cadence service's container, in cadence command connect to the ip which exists in /etc/hosts of that container(our mchine ip in container).
e.g.,
```bash
cadence --address 172.22.0.3:7933 help
```

- If it's the first time you are running examples, so you should first register you domain using command-line.  
e.g.,
```bash
cadence --address 172.22.0.3:7933 --do arrangerlab domain register -rd 1
```
