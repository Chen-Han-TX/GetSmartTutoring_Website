import AuthService from "../services/auth.service";
import React, { useState, useRef, useEffect } from "react";
import DropdownMultiselect from "react-multiselect-dropdown-bootstrap";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import Input from "react-validation/build/input";
import { useNavigate } from "react-router-dom";
import TutoringService from "../services/tutoring.service";
import SubjectServices from "../services/subject.service";

const Tutoring = () => {
    const currentUser = AuthService.getCurrentUser();
    const userType = currentUser.user_type;
    const form = useRef();
    const checkBtn = useRef();
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [school, setSchool] = useState("");
    const [selectedSubjects, setSelectedSubjects] = useState({});
    const [tutorList, setTutorList] = useState({});

    const searchTutor = (subjList) => {
      TutoringService.matchTutors(subjList).then(
        (response) => {
          if (response.status == 202) {
            OnChangeTutorList(response.data);
          } else {
            console.log("bruh, " + response.status);
          }
        },
        (error) => {
          if (error.response.status == 404){
            OnChangeTutorList({})
          }
        }
      );
    }


   

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

    

    const [successful, setSuccessful] = useState(false);
    const [message, setMessage] = useState("");
    
    const navigate = useNavigate();

    
    const OnChangePSLE = (subject) => {
        let updated = selectedSubjects
        updated["PSLE"] = subject
        setSelectedSubjects(selectedSubjects => ({
          ...updated
        }));
        searchTutor(selectedSubjects);
      }
    
    const onChangeOlevel = (subject) => {
        let updated = selectedSubjects
        updated["O-Level"] = subject
        setSelectedSubjects(selectedSubjects => ({
            ...updated
        }));
        searchTutor(selectedSubjects);
    }

    const onChangeAlevel = (subject) => {
        let updated = selectedSubjects
        updated["A-Level"] = subject
        setSelectedSubjects(selectedSubjects => ({
            ...updated
        }));
        searchTutor(selectedSubjects);
    }

    const OnChangeTutorList = (newTutorList) => {
      setTutorList(newTutorList)
    }

    // const results = [];

    // tutorList.forEach((tutor, index) => {
    //   results.push(
    //     <div key={index}>
    //       <h2>name: {tutor.name}</h2>
    //       <h2>country: {tutor.email}</h2>
  
    //       <hr />
    //     </div>,
    //   );
    // });
  
    return (
    <div className="auth-inner">
        <Form ref={form}>
          <h3>Search for tutors!</h3>
          {!successful && (

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

            </div>
          )}
        </Form>
        
        {userType === "Student" && tutorList!=null && tutorList.length!=0 &&(
          <div>
           {/* {tutorList.map(tutor => {
            return (
              
                <p>{tutor.name}</p>
              
            )
           })} */}

           
          </div>
        )}

        </div>
  );
};

  export default Tutoring; 