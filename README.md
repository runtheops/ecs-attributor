# ecs-attributor
ECS container instance attributes registrator

Usage: 
  ecs-attributor -a attribute1=value1,attribute2=value2
  
Dead simple util to cover inability of AWS ECS Agent (work is in progress, but feature isn't released at the moment) to register custom ECS container intance attributes from within the instance itself. Attributor relies of instance IAM role, ECS Agent metadata and user-provided key=value attribute pairs

Options:
  -a=""            comma-separated list of key=value attribute pairs
