import AuthService from "../services/auth.service";
import React, { useState, useRef, useEffect } from "react";
import TutoringService from "../services/tutoring.service";

const Booking = () => {
    const currentUser = AuthService.getCurrentUser();
    const userType = currentUser.user_type;
    const [appList, setAppList] = useState([]);
    let appListItem = "";

    useEffect(() => {
        TutoringService.getApplications().then(
          (response) => {
            setAppList(response)
            if (userType === "Tutor") {
                appListItem = appList.map((val) => (
                    <li>
                      <a>{val.student_id}</a></li>
                ))
                console.log(appListItem)
                ;  
            }
          },
          (error) => {
            if (error.response.status == 404){
              console.log(error)
            }
          }
        );
    }, []); // run useEffect only once, when the component is mounted

   

    

    return(
        <div className="auth-inner">
            <header className="TutorHeader">
                <h4>Handle your bookings here!</h4>
                <div id="example-collapse-text">
                    <ul>{appListItem}</ul>
                </div>
            </header>
        </div>
    )
}

export default Booking; 