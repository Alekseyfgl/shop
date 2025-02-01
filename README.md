```bash
docker build -t go .
```

### Running the Container with Required Environment Variables

```bash

docker run \
  -p 3000:3000 \
  -e SERV_PORT=3000 \
  -e POSTGRES_URI="user=gen_user password=tp+.^7}9)k72_8 host=194.87.76.134 port=5432 dbname=default_db" \
  -e MONGO_URI="" \
  -e JWT_KEY="vkldfgklfd" \
  -e SUPER_ADMIN_LOGIN="admin" \
  -e SUPER_ADMIN_PASSWORD="admin" \
  --name shop-cnt1 \
  shop:1.0.0
```

* ```-p``` 3000:3000 maps the container's port 3000 to the host's port 3000.
* The ```-e``` flags specify the environment variables for the container.
* The ```-d``` flag runs the container in detached mode (in the background).
* ```--name``` my-go-app-cnt assigns a custom name to the container for easier management.

### Restarting a Stopped or Crashed Container

```bash

docker start my-go-app-cnt
```