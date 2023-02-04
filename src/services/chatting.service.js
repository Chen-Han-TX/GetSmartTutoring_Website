import { upload } from "@testing-library/user-event/dist/upload";
import axios from "axios";

const CHATTING_URL = "http://localhost:5070/api/"

axios.defaults.withCredentials = true

let axiosConfig = {
    headers: {
        'Content-Type': 'text/plain',
       // "Access-Control-Allow-Origin": "http://localhost:3000"
    },
    withCredentials : true,
}

const createChatList = () => {
    return axios.post(CHATTING_URL + "createchatlist",  {},
    axiosConfig).then(
        (response) => {
            return response;
        }
    )
}


const sendMsg = (opp_user_id, content) => {
    return axios.post(CHATTING_URL + "sendmessages/" + opp_user_id, { "content": content},
        axiosConfig)
            .then((response) => {
        return response;
      });
};


const getChatList = () => {
    return axios.get(CHATTING_URL + "getlist", 
        axiosConfig)
            .then((response) => {
                localStorage.setItem("chatList", JSON.stringify(response.data));
        return response;
      });
};

const getCurrentChatList = () => {
    return JSON.parse(localStorage.getItem("chatList"));
};

const updateChatList = (chatId, message) => {
    var chatList = JSON.parse(localStorage.getItem("chatList"))
    for (var i = 0; i < chatList.length; i++) {
        if (chatList[i].chat_id === chatId) {
            chatList[i].messages.push(message)
            localStorage.setItem("chatList", JSON.stringify(chatList))
            return chatList
        }
    }
}


const ChattingServices = {
    createChatList,
    sendMsg, 
    getChatList,
    getCurrentChatList,
    updateChatList
}

  
export default ChattingServices;

