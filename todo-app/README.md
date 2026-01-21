**Using Namespace: project**

The `todo-frontend` and `todo-backend` are deployed in separate pods.
Because the application utilizes **Server-Side Rendering (SSR)**, the
frontend communicates with the backend via the internal Kubernetes
Service DNS: `http://todo-backend-svc:80`. This service maps port
**80** to the backend's container port **3000**.

![Architecture](./arch.png)


# Ex-2.6
Moved the `BACKEND` and `FRONTEND` ports to `ConfigMap`
