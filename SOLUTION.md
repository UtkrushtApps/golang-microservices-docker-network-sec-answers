# Solution Steps

1. 1. Redesign Docker Compose networks: define 'frontend' and 'backend' bridge networks with driver options (custom MTU), and assign services accordingly. Only the API joins both networks; db and Redis are backend only.

2. 2. For each service, remove unnecessary 'ports' except for API's published port; use 'expose' for db and cache to allow only internal container-to-container communication.

3. 3. Assign network aliases to services, so internal Docker DNS service discovery is robust. E.g., API is reachable as 'golang-api', Postgres as 'postgres', Redis as 'redis' within backend network.

4. 4. Apply least privilege: for each service, set security_opt 'no-new-privileges:true', drop all capabilities (cap_drop: - ALL), and specify a non-root user (USER in Dockerfile and Compose); for official images, pick the uid/gid for the service user.

5. 5. Edit the Go code: replace all hardcoded 'localhost' or 127.0.0.1 references for Postgres and Redis hosts with environment variables; ensure these use service names ('postgres', 'redis'). Only reference these env vars.

6. 6. In the Go app Dockerfile, add an unprivileged user and use it as USER (non-root) after copying in the built binary.

7. 7. Make sure only the API service is published on the host (via port 8080); remove any external exposure for db and cache. Use 'depends_on' for API on db/cache when needed.

8. 8. Build and run with the new Compose file; verify API-to-db and API-to-redis communication works via Docker DNS, and that database/cache are not exposed externally.

