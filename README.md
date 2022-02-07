# Getting Started

1. Transfer repos and change branches
2. Build images of microservices
3. Deploy images of microservices

# Requirements
- Git
- Go
- Docker
  
# Recommended Software
- VSCode
- Fork
- Studio 3T

# SETUP

1. Run `cloneMicroServices.bat` (FirstTimeOnly) and `setBranchAndPush.bat`
2. Build images `docker-compose -f dev-build.yaml build`
3. Active docker swarm `docker swarm init` (if not active)
4. Deploy images  `docker stack deploy -c dev-deploy.yaml micro-service`

# Helpfull comands

- Destroy swarm `docker stack rm micro-service`