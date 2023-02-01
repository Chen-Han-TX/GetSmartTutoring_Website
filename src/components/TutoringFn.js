import AuthService from "../services/auth.service";
import React, { useState, useRef } from "react";
import DropdownMultiselect from "react-multiselect-dropdown-bootstrap";
import Form from "react-validation/build/form";
import CheckButton from "react-validation/build/button";
import Input from "react-validation/build/input";
import { useNavigate } from "react-router-dom";
import TutoringService from "../services/tutoring.service";
import SubjectServices from "../services/subject.service";

const Tutoring = () => {
    const currentUser = AuthService.getCurrentUser();
    const form = useRef();
    const checkBtn = useRef();
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [school, setSchool] = useState("");
    const [selectedSubjects, setSelectedSubjects] = useState({});

    const searchTutor = (subjList) => {
      TutoringService.matchTutors(subjList).then(
        (response) => {
          if (response.status == 202) {
            console.log(response)
          } else {
            console.log("bruh, " + response);
          }
        },
        (error) => {
          console.log(error);
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

    const onChangeEmail = (e) => {
        const email = e.target.value;
        setEmail(email);
    };
    
    const onChangeName = (e) => {
        const name = e.target.value;
        setName(name);
    };


    const onChangePassword = (e) => {
        const password = e.target.value;
        setPassword(password);
    };

    const onChangeSchool = (e) => {
        const school = e.target.value;
        setSchool(school);
    };

    const OnChangePSLE = (subject) => {
        let updated = selectedSubjects
        updated["PSLE"] = subject
        setSelectedSubjects(selectedSubjects => ({
          ...updated
        }))
        searchTutor(selectedSubjects);
        ;
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

    const handleRegister = (e) => {
        e.preventDefault();
    
        setMessage("");
        setSuccessful(false);
    
        form.current.validateAll();
    
        // If passed validation, call auth service to send the API request
        if (checkBtn.current.context._errors.length === 0) {
          
          AuthService.register_student(name, email, password, school, selectedSubjects).then(
            (response) => {
              if (response.status == 200) {
                setMessage("Applied successfully!");
                setSuccessful(true);
                setTimeout(function () {
                  navigate("/login");
                  window.location.reload();
                }, 2000);
              } else {
                setMessage("Error, try again!");
                setSuccessful(false);
                setTimeout(function () {
                  setMessage("");
                }, 5000);
              }
            },
            (error) => {
              var resMessage = ""
              if (error.response.status === 406) {
                resMessage = "This email has been registered!"
              } else {
                resMessage =
                (error.response &&
                  error.response.data &&
                  error.response.data.message) ||
                error.message ||
                error.toString();
              }
              setMessage(resMessage);
              setSuccessful(false);
            }
          );
        }
      };
  
  
    return (
    <div className="auth-inner">
        <Form onSubmit={handleRegister} ref={form}>
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
        </div>
  );
};

  export default Tutoring; 