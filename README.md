# psychic-sniffle

# docker deployment

setup .env file (see https://hub.docker.com/_/postgres for postgres setup!)
- PSQL_HOST
- PSQL_PORT
- PSQL_USER
- PSQL_PASSWORD
- PSQL_DBNAME
docker-compose pull
docker-compose up
curl http://localhost:8080
- successful response should be {"database_status":"connected","message":"psychic-sniffle server is running"}                                                                                                                      

# kubernetes deployment

setup demployment.yaml
minikube status
minikube start
kubectl apply -f deployment.yaml
kubectl apply -f go-app-deployment.yaml
kubectl expose deployment go-http-postgres --type=NodePort --port=8080
minikube service go-http-postgres
- successful response should be {"database_status":"connected","message":"psychic-sniffle server is running"}
kubectl port-forward svc/go-http-postgres 8080:8080 (to reach server at http://localhost:8080/ when running in Kubernetes)


# kubernetes commands w/ examples

- Monitor Pods: kubectl logs -l app=go-http-postgres -f

- Check Deployment & Replicas: kubectl get deployment go-http-postgres
    + go-http-postgres   3/3     3            3           123m
    + go-http-postgres-6466967844-mkz5t   1/1     Running   2 (20m ago)   21m
    + go-http-postgres-6466967844-pq5z2   1/1     Running   0             91s
    + go-http-postgres-6466967844-qs27r   1/1     Running   0             91s

- List the Pods: kubectl get pods -l app=go-http-postgres
    + go-http-postgres-6466967844-mkz5t   1/1     Running   2 (20m ago)   21m    10.244.0.9    minikube   <none>           <none>
    + go-http-postgres-6466967844-pq5z2   1/1     Running   0             110s   10.244.0.12   minikube   <none>           <none>
    + go-http-postgres-6466967844-qs27r   1/1     Running   0             110s   10.244.0.13   minikube   <none>           <none>
    + postgres-5d8755b89f-vpqrt           1/1     Running   0             21m    10.244.0.8    minikube   <none>           <none>

- Pod Details: kubectl describe pod <pod-name>
    Name:             go-http-postgres-6466967844-pq5z2
    Namespace:        default
    Priority:         0
    Service Account:  default
    Node:             minikube/192.168.49.2
    Start Time:       Sun, 06 Apr 2025 18:25:40 -0600
    Labels:           app=go-http-postgres
                    pod-template-hash=6466967844
    Annotations:      <none>
    Status:           Running
    IP:               10.244.0.12
    IPs:
    IP:           10.244.0.12
    Controlled By:  ReplicaSet/go-http-postgres-6466967844
    Containers:
    go-http-postgres:
        Container ID:   docker://368549c45408ad0db3555e23ed1de30eb7ba54c707a5020cf164df6d7e286c83
        Image:          peterjbishop/go-http-postgres:latest
        Image ID:       docker-pullable://peterjbishop/go-http-postgres@sha256:81e57864edf3f6d4702a97c993c22e8f338ba304b9e80a78be5242cfc1255b29
        Port:           8080/TCP
        Host Port:      0/TCP
        State:          Running
        Started:      Sun, 06 Apr 2025 18:25:41 -0600
        Ready:          True
        Restart Count:  0
        Environment:
        DB_HOST:      postgres
        DB_PORT:      5432
        DB_USER:      <set to the key 'PSQL_USER' in secret 'db-secret'>        Optional: false
        DB_PASSWORD:  <set to the key 'PSQL_PASSWORD' in secret 'db-secret'>    Optional: false
        DB_NAME:      <set to the key 'PSQL_DBNAME' of config map 'db-config'>  Optional: false
        Mounts:
        /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-fm5rr (ro)
    Conditions:
    Type                        Status
    PodReadyToStartContainers   True 
    Initialized                 True 
    Ready                       True 
    ContainersReady             True 
    PodScheduled                True 
    Volumes:
    kube-api-access-fm5rr:
        Type:                    Projected (a volume that contains injected data from multiple sources)
        TokenExpirationSeconds:  3607
        ConfigMapName:           kube-root-ca.crt
        ConfigMapOptional:       <nil>
        DownwardAPI:             true
    QoS Class:                   BestEffort
    Node-Selectors:              <none>
    Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                                node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
    Events:
    Type    Reason     Age    From               Message
    ----    ------     ----   ----               -------
    Normal  Scheduled  2m12s  default-scheduler  Successfully assigned default/go-http-postgres-6466967844-pq5z2 to minikube
    Normal  Pulling    2m12s  kubelet            Pulling image "peterjbishop/go-http-postgres:latest"
    Normal  Pulled     2m11s  kubelet            Successfully pulled image "peterjbishop/go-http-postgres:latest" in 635ms (636ms including waiting). Image size: 17367249 bytes.
    Normal  Created    2m11s  kubelet            Created container: go-http-postgres
    Normal  Started    2m11s  kubelet            Started container go-http-postgres

- Pod Distribution: kubectl get pods -o wide
    NAME                                READY   STATUS    RESTARTS        AGE     IP            NODE       NOMINATED NODE   READINESS GATES
    go-http-postgres-6466967844-mkz5t   1/1     Running   2 (3h43m ago)   3h44m   10.244.0.9    minikube   <none>           <none>
    go-http-postgres-6466967844-pq5z2   1/1     Running   0               3h24m   10.244.0.12   minikube   <none>           <none>
    go-http-postgres-6466967844-qs27r   1/1     Running   0               3h24m   10.244.0.13   minikube   <none>           <none>
    postgres-5d8755b89f-vpqrt           1/1     Running   0               3h44m   10.244.0.8    minikube   <none>           <none>