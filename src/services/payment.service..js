import { upload } from "@testing-library/user-event/dist/upload";
import axios from "axios";

const PAYMENT_URL = "https://payment-4dcnj7fm6a-uc.a.run.app/api/"

axios.defaults.withCredentials = true

let axiosConfig = {
    headers: {
        'Content-Type': 'text/plain',
       // "Access-Control-Allow-Origin": "http://localhost:3000"
    },
    withCredentials : true,
}

const makePayment = (payment) => {
    return axios.post(PAYMENT_URL + "payment", payment,
        axiosConfig)
            .then((response) => {
        return response;
      });
};

const PaymentService = {
    makePayment
}

  
export default PaymentService;

