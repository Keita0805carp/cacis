# Sample
- red:   GPIO20
- green: GPIO21

## Build/Push

```bash
docker build -f ./red.Dockerfile -t localhost:32000/led:red .
docker build -f ./green.Dockerfile -t localhost:32000/led:green .

docker push localhost:32000/led:red
docker push localhost:32000/led:green
```

## Apply

```bash
kubectl apply -f red-ds.yaml
kubectl apply -f green-ds.yaml
```
