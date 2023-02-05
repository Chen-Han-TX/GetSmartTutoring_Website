import { upload } from "@testing-library/user-event/dist/upload";
import axios from "axios";

const TUTORING_URL = "https://tutoring-4dcnj7fm6a-uc.a.run.app/api/tutoring/"

axios.defaults.withCredentials = true

let axiosConfig = {
    headers: {
        'Content-Type': 'text/plain',
       // "Access-Control-Allow-Origin": "http://localhost:3000"
    },
    withCredentials : true,
}

const matchTutors = (searchedSubjects) => {
    return axios.post(TUTORING_URL + "matchtutors", searchedSubjects,
        axiosConfig)
            .then((response) => {
        return response;
      });
};

const applyForTutor = (application) => {
    return axios.post(TUTORING_URL + "apply", application,
        axiosConfig)
            .then((response) => {
        return response;
      });
};

const getApplications = (user_id, user_type) => {
    return axios.get(TUTORING_URL + "getapplications/" + user_id + "/" + user_type, 
        axiosConfig)
            .then((response) => {
        return response.data;
      });
};

const handleApplications = (user_type, application) => {
    return axios.post(TUTORING_URL + "handleapplications/" + user_type, application,
        axiosConfig)
            .then((response) => {
        return response;
      });
};

const TutoringService = {
    matchTutors,
    applyForTutor,
    getApplications,
    handleApplications
}

  
export default TutoringService;

