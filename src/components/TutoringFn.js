import AuthService from "../services/auth.service";
import React, { useState, useEffect } from "react";
import { Modal, Button} from 'react-bootstrap';
import DropdownMultiselect from "react-multiselect-dropdown-bootstrap";
import { Dropdown, DropdownButton } from 'react-bootstrap';
import TutoringService from "../services/tutoring.service";
import SubjectServices from "../services/subject.service";
import Card from 'react-bootstrap/Card';


const Tutoring = () => {

    const currentUser = AuthService.getCurrentUser();
    const [selectedSubjects, setSelectedSubjects] = useState({
      "PSLE" : [],
      "O-Level": [],
      "A-Level": []
    });
    const [tutorList, setTutorList] = useState([]);
    const [listItemsTutors, setListItemTutors] = useState("")
    const [clickedTutor, setClickedTutor] = useState(null)
    const [chosenSubject, setChosenSubject] = useState("");

    const handleSelect = (eventKey) => {
      setChosenSubject(eventKey);
    };

    const [show, setShow] = useState(false);

    const handleClose = () => {
      setClickedTutor(null);
      setChosenSubject("");
      setShow(false);
    }

    const handleShow = (tutor) => {
      setClickedTutor(tutor)
      setShow(true);
    }


    useEffect(() => {
      // displaying the Tutor Card 
      setListItemTutors(Array.isArray(tutorList) ? tutorList.map((tutor, index) =>
        <div key={"tutor_"+index}>
          <br />
          <Card style={{ width: '100%' }}>
            <Card.Body>
              <Card.Title>{index+1 + ". " + tutor.Name}</Card.Title>
              <Card.Text>
                <div>
                  <p>
                    <strong>Availability</strong> <br />
                    {Object.entries(tutor.Availability).map(([day, schedule], index) => [
                      day + ': ' + schedule.start.slice(0, 2) + schedule.start.slice(3) + ' to ' + schedule.end.slice(0, 2) + schedule.end.slice(3), 
                      <br key={index} />
                    ])}
                  </p>
                  <p>
                    <strong>Hourly Rates</strong> <br />
                    $ {tutor.HourlyRate}
                    <br />
                  </p>
                  <p>
                    <strong>Matched Subjects</strong> <br />
                    {tutor.MatchedSubjectList.map((subject, index) => [subject, <br key={index} />])}
                  </p>
                </div>
              <hr className="hr"></hr>
              </Card.Text>
              <Button variant="primary" onClick={() => handleShow(tutor)}>Book A Session</Button>
            </Card.Body>
          </Card>
        </div>
        ) : []);

    }, [tutorList]);


    const subjects = SubjectServices.getAllSubjects()
    const PSLESubjects = subjects["PSLE"].sort()
    const OlevelSubjects = subjects["O-Level"].sort()
    const AlevelSubjects = subjects["A-Level"].sort()

    const PSLEArray = [];
    for (let i = 0; i < PSLESubjects.length; i++) {
        var value = PSLESubjects[i]
        PSLEArray.push({key: value, label: value })
    } 

    const OLevelArray = [];
    for (let i = 0; i < OlevelSubjects.length; i++) {
        var value = OlevelSubjects[i] 
        OLevelArray.push({key: value, label: value })
    } 

    const ALevelArray = [];
    for (let i = 0; i < AlevelSubjects.length; i++) {
        var value = AlevelSubjects[i] 
        ALevelArray.push({key: value, label: value })
    } 
    
    const OnChangePSLE = (subject) => {
        let updated = selectedSubjects
        updated["PSLE"] = subject
        setSelectedSubjects(selectedSubjects => ({
          ...updated
        }));
      }
    
    const onChangeOlevel = (subject) => {
        let updated = selectedSubjects
        updated["O-Level"] = subject
        setSelectedSubjects(selectedSubjects => ({
            ...updated
        }));
    }

    const onChangeAlevel = (subject) => {
        let updated = selectedSubjects
        updated["A-Level"] = subject
        setSelectedSubjects(selectedSubjects => ({
            ...updated
        }));
    }


    const handleSearchTutor = (e) => {
      e.preventDefault();
      
      if (selectedSubjects["PSLE"].length === 0 && selectedSubjects["O-Level"].length === 0 && selectedSubjects["A-Level"].length === 0) {
        alert("Please select at least one subject!")
        return
      }

      TutoringService.matchTutors(selectedSubjects).then(
        (response) => {
          console.log(response)
          if (response.status === 202) {
            setTutorList(response.data);
          } else {
            console.log("response status: " + response.status);
          }
        },
        (error) => {
          if (error.response.status == 404){
            setTutorList({})
            alert("No tutor matched!")
          }
        }
      );
    }

    const handleBookSession = (e) => {
      e.preventDefault();
      if (chosenSubject === "") {
        alert("Please choose one subject!")
        return
      } 

      // TO-DO
      console.log("sending request") 
      console.log("Chosen subject: ", chosenSubject)
      console.log("tutor", clickedTutor.Name)
    }

    return (
    <div className="auth-inner">
        <div>
          <h3>Search for tutors!</h3>
            <div>
              <hr className="hr"></hr>
              <h4>Choose the subjects you are looking for.</h4>
            
              <hr className="hr"></hr>
              
              <div className="mb-3">
                  <label htmlFor="options">PSLE</label>
                  <DropdownMultiselect options={PSLEArray} name="pslesubjects" 
                  handleOnChange={(selected) => {
                    OnChangePSLE(selected);
                  }}/>
              </div>

              <div className="mb-3">
                  <label htmlFor="options">O-Level</label>
                  <DropdownMultiselect options={OLevelArray} name="olevelsubjects"
                  handleOnChange={(selected) => {
                    onChangeOlevel(selected);
                  }}/>
              </div>

              <div className="mb-3">
                  <label htmlFor="options">A-Level</label>
                  <DropdownMultiselect options={ALevelArray} name="alevelsubjects"
                  handleOnChange={(selected) => {
                    onChangeAlevel(selected);
                  }}/>
              </div>
              <hr className="hr"></hr>

              <div className="d-grid">
                <button onClick={handleSearchTutor} className="btn btn-success btn-block">Search</button>
              </div>
            </div>
  
        </div>

        <div>
        <hr className="hr"></hr>
            {listItemsTutors}
        </div>
        {clickedTutor && (
          <Modal show={show} onHide={handleClose}>
          <Modal.Header closeButton>
            <Modal.Title>Book a Session</Modal.Title>
          </Modal.Header>
          <Modal.Body>
          <div>
            <p>
              <strong>Tutor Name</strong>
              <br />
              {clickedTutor.Name}
            </p>
            <p>
              <strong>Availability</strong> <br />
              {Object.entries(clickedTutor.Availability).map(([day, schedule], index) => [
                day + ': ' + schedule.start.slice(0, 2) + schedule.start.slice(3) + ' to ' + schedule.end.slice(0, 2) + schedule.end.slice(3), 
                <br key={index} />
              ])}
            </p>
            <p>
              <strong>Hourly Rates</strong> <br />
              $ {clickedTutor.HourlyRate}
              <br />
            </p>
            <p>
              <strong>Choose Subject</strong> <br />
              <Dropdown onSelect={handleSelect}>
                <Dropdown.Toggle variant="success" id="dropdown-basic">
                  {chosenSubject || "Select a subject"}
                </Dropdown.Toggle>

                <Dropdown.Menu>
                  {clickedTutor.MatchedSubjectList.map((subject) => (
                    <Dropdown.Item key={subject} eventKey={subject}>
                      {subject}
                    </Dropdown.Item>
                  ))}
                </Dropdown.Menu>
              </Dropdown>
            </p>
          </div>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleClose}>Cancel</Button>
            <Button variant="primary" onClick={handleBookSession}>Book</Button>
          </Modal.Footer>
        </Modal>

        )
      }
    </div>
  );
};

  export default Tutoring; 