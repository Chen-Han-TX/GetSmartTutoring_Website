import AuthService from "../services/auth.service";
import React, { useState, useRef, useEffect } from "react";
import TutoringService from "../services/tutoring.service";
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';
const Booking = () => {
    const currentUser = AuthService.getCurrentUser();
    const userType = currentUser.user_type;
    const [appList, setAppList] = useState([]);
    const [listItemsApps, setListItemsApps] = useState("")

    useEffect(() => {
        TutoringService.getApplications().then(
            (response) => {
              setAppList(response)
              console.log(response)
            },
            (error) => {
              if (error.response.status == 404){
                console.log(error)
              }
            }
          );

      }, []);

      const handleAccept = app => {
        app.application_status = "Accepted"
        TutoringService.handleApplications(app).then(
            (response) => {
              alert("Tutoring application has been accepted!")
              window.location.reload(true)

            },
            (error) => {
              if (error.response.status == 404){
                console.log(error)
              }
            }
          );
      };
      
      const handleReject = app => {
        app.application_status = "Rejected"
        TutoringService.handleApplications(app).then(
            (response) => {
              alert("Tutoring application has been rejected!")
              window.location.reload(true)
            },
            (error) => {
              if (error.response.status == 404){
                console.log(error)
              }
            }
          );
      };

      useEffect(() => {
        // displaying the Application Cards
        setListItemsApps(Array.isArray(appList) ? appList.map((app, index) =>
        <Card key={"app_"+index} style={{ width: '100%' }}>
        <Card.Body>
            <Card.Title>{index+1 + ". " + app.student_name}</Card.Title>
            <Card.Text>
            Subject requested: {app.subject}
            </Card.Text>
            <Card.Text>
            Session length: {app.session_length} hours
            </Card.Text>
            <Card.Text>
            Status: {app.application_status}
            </Card.Text>
            {app.application_status === "Pending" && (
              <>
                <Button variant="success" onClick={() => handleAccept(app)}>Accept</Button>
                <Button variant="danger" onClick={() => handleReject(app)}>Reject</Button>
              </>
            )} {app.application_status === "Accepted" && (
                <>
                  <Button variant="success" disabled>
                    Accepted
                  </Button>
                </>
              )}{app.application_status === "Rejected" && (
                <>
                  <Button variant="danger" disabled>
                    Rejected
                  </Button>
                </>
              )}
        </Card.Body>
        </Card>
        ) : []);
      }, [appList]);

   

    

    return(
        <div className="auth-inner">
            <header className="TutorHeader">
                <h4>Handle your bookings here!</h4>
                <div id="example-collapse-text">
                    {listItemsApps}
                </div>
            </header>
        </div>
    )
}

export default Booking; 