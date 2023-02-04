import AuthService from "../services/auth.service";
import React, { useState, useRef, useEffect } from "react";
import TutoringService from "../services/tutoring.service";
import ChattingServices from "../services/chatting.service";
import PaymentService from "../services/payment.service.";

import Card from 'react-bootstrap/Card';
import {Button, Modal} from 'react-bootstrap';

const Booking = () => {
    const currentUser = AuthService.getCurrentUser();
    const userType = currentUser.user_type;
    const [appList, setAppList] = useState([]);
    const [listItemsApps, setListItemsApps] = useState("")

    const [appToPay, setAppToPay] = useState(null);

    // For confirm payment modal
    const [show, setShow] = useState(false);

    const handleClose = () => {
      setShow(false)
      setAppToPay(null)
    }
    const handleShow = (app) => {
      setShow(true) 
      setAppToPay(app);
    }

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
              ChattingServices.createChatList().then(
                (response) => {
                  console.log(response)
                  if (response.data.length > 0) {
                    alert("A new chatroom has been created for you!")
                    window.location.href = "/chat"
                  } else {
                    window.location.reload(true)
                  }
                },
                (error) => {
                  if (error.response.status == 404){
                    console.log(error)
                    alert(error)
                  }
                }
              );

            },
            (error) => {
              if (error.response.status == 404){
                console.log(error)
                alert(error)
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

      const makePayment = app => {
        const payment = {
          "amount": app.hourly_rate * app.session_length,
          "tutor_id": app.tutor_id,
          "student_id": app.student_id,
          "session_id": app.session_id
        }
        PaymentService.makePayment(payment).then(
          (response) => {
            alert("Payment has been made!")
            window.location.reload(true)
          },
          (error) => {
            if (error.response.status == 404){
              console.log(error)
            }
          }
        );
      }

      useEffect(() => {
        // displaying the Application Cards
        console.log(listItemsApps)
        setListItemsApps(Array.isArray(appList) ? appList.map((app, index) =>
        <Card key={"app_"+index} style={{ width: '100%' }}>
        <Card.Body>
          {userType==="Tutor" && (
            <Card.Title>{index+1 + ". " + app.student_name}</Card.Title>
          )}
          {userType==="Student" && (
            <Card.Title>{index+1 + ". " + app.tutor_name}</Card.Title>
          )}
            <Card.Text>
            Subject requested: {app.subject}
            </Card.Text>
            <Card.Text>
            Session length: {app.session_length} hours
            </Card.Text>
            <Card.Text>
            Hourly Rate: $ {app.hourly_rate}
            </Card.Text>

            <Card.Text>
            Status: {app.application_status}
            </Card.Text>
            {app.application_status === "Pending" && userType==="Tutor" && (
              <>
                <Button variant="success" onClick={() => handleAccept(app)}>Accept</Button>
                <Button variant="danger" onClick={() => handleReject(app)}>Reject</Button>
              </>
            )} 
             {app.application_status === "Accepted" && userType === "Student" && (
                <>

                  <Button variant="success" onClick={() => handleShow(app)} >
                    Pay $ {app.session_length * app.hourly_rate}
                  </Button>
                </>
              )}
              
              {app.application_status === "Rejected"  && (
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
              {userType==="Tutor" && (
                <h4>Handle your bookings here!</h4>
              )}{userType==="Student" && (
                <h4>View your booking statuses!</h4>
              )}
                
                <div id="example-collapse-text">
                    {listItemsApps}
                </div>
            </header>
            <Modal show={show} onHide={handleClose}>
              <Modal.Header closeButton>
                <Modal.Title>Confirmation</Modal.Title>
              </Modal.Header>
              { appToPay && (
                  <Modal.Body>Are you confirmed to pay a total amount of ${appToPay.hourly_rate * appToPay.session_length} to {appToPay.tutor_name}?</Modal.Body>
              )}
              <Modal.Footer>
                <Button variant="secondary" onClick={handleClose}>
                  Cancel
                </Button>
                <Button variant="primary" onClick={() => makePayment(appToPay)}>
                  Confirm
                </Button>
              </Modal.Footer>
            </Modal>
        </div>
    )
}

export default Booking; 