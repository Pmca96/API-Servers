docker build -t user_server .
docker service update API_user-server --image=user_server:local --force