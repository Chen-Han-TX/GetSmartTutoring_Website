# Get Smart Tutoring (GST)

Get Smart Tutoring, or _**GST**_ for short, is a tutoring service web application that utilizes React and microservices, containerized and hosted on Google Cloud Kubernetes. 

# Architecture Diagram

# Microservice Design Considerations
### 1. Authentication Management
- User Registration

  ![image](https://user-images.githubusercontent.com/73156798/216923148-bbe91320-07d9-4b09-b81a-f4f6755a455c.png)
  
- User Login with credentials, set non-critical user data into local storage

  ![image](https://user-images.githubusercontent.com/73156798/216917763-783e3d68-a48f-4cb5-82cd-28171714d17f.png)
  
- User Logout
### 2. Tutoring Functionalities
- Student search for tutors based on subjects of interest

  ![image](https://user-images.githubusercontent.com/73156798/216924054-70069754-eeae-4a8e-a387-3237c4570290.png)


- Student apply for tutoring

  ![image](https://user-images.githubusercontent.com/73156798/216925103-f392a524-1385-4979-a603-a33e24cd69ec.png)


- Student/Tutor get all applications sent/received

  ![image](https://user-images.githubusercontent.com/73156798/216925473-07b032f0-8ad8-49d9-9946-c4ea027f303a.png)


- Tutor accept/deny tutoring applications

  ![image](https://user-images.githubusercontent.com/73156798/216925695-76be51b4-cc82-4ee6-9bc8-f3c2f0c4db78.png)


### 3. Chatroom 
- Create chatroom once tutor accepts application
- User get a list of created chats
- User send message in a specific chatroom

  ![image](https://user-images.githubusercontent.com/73156798/216926693-eb2b41ce-7ff3-4787-b502-5e78766a0661.png)

### 4. Payment
- User make payment for tutoring session

  ![image](https://user-images.githubusercontent.com/73156798/216927211-5ab436a1-7dae-4df8-aeaa-5d9678dedb57.png)

### 5. Subjects
- Get a list of subjects, with the ability to specify between `all`, `psle`, `olevel` and `alevel`

# Cloud Native Design, Architecture Rationale and Resiliency

- Kubernetes replicas being set to 3; In the case that any of the pods go down, another one will be reinstated in place of the previous.
- Cloud Scalability through Google Cloud, where auto-scaling is set to listen to CPU and memory usage to tune the number of resources being deployed.
# Development Tools & Methods
Get Smart Tutoring was developed using the following tools:
- Developed using **Microservice architecture**, specifically for the Golang backend.
- **Back-end**: [Golang server](https://go.dev/) using [Mux Router](https://github.com/gorilla/mux)
- **Front-end**: [React.js](https://reactjs.org/) & 
             [Bootstrap](https://getbootstrap.com/)
- **CORS Middleware Tool**: [Moesif Origin & CORS Changer](https://chrome.google.com/webstore/detail/moesif-origin-cors-change/digfbfaphojjndkpccljibejjbppifbc) (Used for testing)
- **CORS Middleware Package**: [rs/cors - GO Package](https://github.com/rs/cors)
- **Database**: [Firebase](https://firebase.google.com)
- **Cloud deployment**: [Google Cloud](https://cloud.google.com)

# Deployment
You can try out our latest deployment of the application on ()

# Credits
- This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).
- More stuff
