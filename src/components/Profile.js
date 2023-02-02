import React, { useState, useRef } from "react";
import AuthService from "../services/auth.service";
import Button from 'react-bootstrap/Button';
import Collapse from 'react-bootstrap/Collapse';


const Profile = () => {

  const currentUser = AuthService.getCurrentUser();
  const userType = currentUser.user_type;
  const [open, setOpen] = useState(false);
  const [open2, setOpen2] = useState(false);
  const [open3, setOpen3] = useState(false);
  var aoiDict = currentUser.area_of_interest
  var area_of_interest = []
  var listItems = "";

  var availDict = {};
  var availability = [];
  var listItemsAvail = "";

  var certificates = [];
  var listItemsCert = "";

  for (let k in aoiDict) {
    if (aoiDict[k] != null){
        for (var i=0; i<aoiDict[k].length; i++) {
            area_of_interest.push(k+"-"+aoiDict[k][i])
        }
    }
  }
  listItems = (area_of_interest.map(
    (interest) =>
    <li>{interest}</li>
  ));
  
  if (userType === "Tutor") {
    availDict = currentUser.availability
    for (let k in availDict) {
        availability.push(k+": "+ availDict[k]["start"] + " to " + availDict[k]["end"])
    }
    listItemsAvail = (availability.map(
      (avail) =>
      <li>{avail}</li>
    ));  

     
    certificates = currentUser.cert_of_evidence
    listItemsCert = (certificates.map(
      (link, index) =>
      <li>
        <a href={link}>{"File " + (index + 1)}</a></li>
    ));  
  }

  return (
    <div className="auth-inner">
        <header className="jumbotron">
          <h3>
          Your Profile
          </h3>
          <br></br>
        </header>
      <p>
        <strong>Name:</strong> {currentUser.name}  
      </p>
      <p>
        <strong>Email:</strong> {currentUser.email} 
      </p>
      {userType === "Student" && (
        <div>
        <p>
        <strong>School:</strong> {currentUser.school} 
        </p>

        <p>
        <Button
            onClick={() => setOpen(!open)}
            aria-controls="example-collapse-text"
            aria-expanded={open}> Area of Interests </Button>
        </p>
        <Collapse in={open}>
        <div id="example-collapse-text">
            <ul>{listItems}</ul>
        </div>
        </Collapse>

        </div>
      )}
      {userType === "Tutor" && (
        <div>
         <p>
        <strong>Hourly Rate:</strong> ${currentUser.hourly_rate} 
        </p>

        <p>
        <Button
             onClick={() => setOpen(!open)}
            aria-controls="example-collapse-text"
            aria-expanded={open}> Proficient Subjects </Button>
        </p>
        <Collapse in={open}>
        <div id="example-collapse-text">
            <ul>{listItems}</ul>
        </div>
      </Collapse>

        <p>
        <Button
            onClick={() => setOpen2(!open2)}
            aria-controls="example-collapse-text"
            aria-expanded={open2}> Availability </Button>
        </p>
        <Collapse in={open2}>
        <div id="example-collapse-text">
            <ul>{listItemsAvail}</ul>
        </div>
        </Collapse>

        <p>
        <Button
            onClick={() => setOpen3(!open3)}
            aria-controls="example-collapse-text"
            aria-expanded={open3}> Certificate of Evidence </Button>
        </p>
        <Collapse in={open3}>
        <div id="example-collapse-text">
            <ul>{listItemsCert}</ul>
        </div>
        </Collapse>


        </div>
      )}




    </div>
  )
};
      
export default Profile;