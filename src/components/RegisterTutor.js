import React, { useState, useRef, useEffect } from "react";
import * as ReactDOM from 'react-dom';
import Form from "react-validation/build/form";
import Input from "react-validation/build/input";
import CheckButton from "react-validation/build/button";
import DropdownMultiselect from "react-multiselect-dropdown-bootstrap";
import { isEmail } from "validator";

import AuthService from "../services/auth.service";
import SubjectServices from "../services/subject.service";
import { useNavigate } from "react-router-dom";

import TimePicker from 'react-bootstrap-time-picker';

import { MDBFile } from 'mdb-react-ui-kit';

const required = (value) => {
  if (!value) {
    return (
      <div className="invalid-feedback d-block">
        This field is required!
      </div>
    );
  }
};

const validEmail = (value) => {
  if (!isEmail(value)) {
    return (
      <div className="invalid-feedback d-block">
        This email cannot be registered!
      </div>
    );
  }
};


const vpassword = (value) => {
  if (value.length < 8 || value.length > 40) {
    return (
      <div className="invalid-feedback d-block">
        The password must be between 8 and 40 characters.
      </div>
    );
  }
};


const RegisterTutor = () => {
  const form = useRef();
  const checkBtn = useRef();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [hourlyRate, setHourlyRate] = useState(50);
  const [startTime, setStartTime] = useState("");
  const [startForEndTime, setStartForEndTime] = useState("11:00");
  const [endTime, setEndTime] = useState("");
  const [availability, setAvailability] = useState({});

  var uploadedDocuments = [];


  var selectedSubjects = {
    "PSLE": [], "O-Level": [], "A-Level": []
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

  const weekdayArray = [
    {key: "Monday", label: "Monday"},
    {key: "Tuesday", label: "Tuesday"},
    {key: "Wednesday", label: "Wednesday"},
    {key: "Thursday", label: "Thursday"},
    {key: "Friday", label: "Friday"},
    {key: "Saturday", label: "Saturday"},
    {key: "Sunday", label: "Sunday"},
  ];


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

  const onChangeHourlyRate = (e) => {
    const hourlyRate = e.target.value;
    setHourlyRate(hourlyRate)
    document.getElementById("hourlyRate").textContent = "$"+hourlyRate
  }

  const onChangeStartTime = (e) => {
    const startTime = new Date(e * 1000).toISOString().substring(11, 16)

    // Add 30 mins 
    const startForEndTime = new Date((e + 1800) * 1000).toISOString().substring(11, 16)

    setStartTime(startTime)
    setStartForEndTime(startForEndTime)
  }

  const onChangeEndTime = (e) => {
    const endTime = new Date(e * 1000).toISOString().substring(11, 16)
    setEndTime(endTime)
  }

  const OnChangePSLE = (subject) => {
    selectedSubjects["PSLE"] = subject
  }

  const onChangeOlevel = (subject) => {
    selectedSubjects["O-Level"] = subject
  }

  const onChangeAlevel = (subject) => {
    selectedSubjects["A-Level"] = subject
  }

  const onChangeWeekDays = (weekdays) => {
    let updatedValue = {};
    for (var i=0; i<weekdays.length; i++) {
      updatedValue[weekdays[i]] = {
        "start": "", "end": ""
      }
    }
    setAvailability(availability => ({
          ...updatedValue
        }));
  }

  const handleRegister = (e) => {
    e.preventDefault();

    setMessage("");
    setSuccessful(false);

    form.current.validateAll();

    // If passed validation, call auth service to send the API request
    if (checkBtn.current.context._errors.length === 0) {

      // Consolidate the availbility
      let updatedValue = availability
      var weekdays = Object.keys(availability)
      for (var i=0; i<weekdays.length; i++) {
        updatedValue[weekdays[i]]["start"] = startTime
        updatedValue[weekdays[i]]["end"] = endTime
      }
      setAvailability(availability => ({
            ...availability,
            ...updatedValue
          }));
          
      /*
      AuthService.register_tutor(name, email, password, selectedSubjects, uploadedDocuments).then(
        (response) => {
          if (response.status === 200) {
            setMessage("Registered Successfully!");
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
      */


    }
  };

  return (
        <Form onSubmit={handleRegister} ref={form}>
          <h3>Register Tutor</h3>
          {!successful && (

            <div>
              <hr className="hr" />
              <h4>Basic Details</h4>
              <div className="mb-3">
                <label htmlFor="name">Name</label>
                <Input
                  type="text"
                  className="form-control"
                  name="name"
                  value={name}
                  onChange={onChangeName}
                  validations={[required]}
                />
              </div>

              <div className="mb-3">
                <label htmlFor="email">Email</label>
                <Input
                  type="text"
                  className="form-control"
                  name="email"
                  value={email}
                  onChange={onChangeEmail}
                  validations={[required, validEmail]}
                />
              </div>

              <div className="mb-3">
                <label htmlFor="password">Password</label>
                <Input
                  type="password"
                  className="form-control"
                  name="password"
                  value={password}
                  onChange={onChangePassword}
                  validations={[required, vpassword]}
                />
              </div>

              <div className="mb-3">
                <label htmlFor="hourlyrate">Indicated Hourly Rate<div id="hourlyRate">$50</div></label>
                <Input
                  type="range"
                  min="10"
                  max="100"
                  step={5}
                  value={hourlyRate}
                  onChange={onChangeHourlyRate}
                  validations={[required]}
                   />
              </div>

              <hr className="hr" />
              <h4>Availabilities</h4>
              <div className="mb-3">
                <div className="mb-4" id="addAvail">

                  <DropdownMultiselect 
                    options={weekdayArray} 
                    name="weekdayArray"
                    handleOnChange={(selected) => {
                    onChangeWeekDays(selected);
                  }}/>

                  <label>Start</label>
                  <TimePicker 
                      start="00:00"
                      end="23:59"
                      step={30} 
                      value={startTime}
                      onChange={onChangeStartTime}
                      validations={[required]}
                      />
                  <label>End</label>
                  <TimePicker 
                      start={startForEndTime}
                      end="23:59"
                      step={30} 
                      value={endTime}
                      onChange={onChangeEndTime}
                      validations={[required]}
                      />
                </div>
              </div>

              <hr className="hr"></hr>
              <h4>Proficient Subjects</h4>
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

              <hr className="hr" />
              <h4>Certificate of Evidence</h4>
              <div className="alert alert-primary" role="alert"> 
                <p className="mb-0"> 
                  You may submit any document proving your ability as a Tutor for subjects indicated above
                </p>
              </div>

              <div className="mb-3">
                <input className="form-control" type="file" id="formFileMultiple" multiple />
              </div>
              

              <div className="d-grid">
                <button className="btn btn-success btn-block">Register</button>
              </div>
            </div>
          )}

          {message && (
            <div className="mb-3">
              <div
                className={
                  successful ? "alert alert-success" : "alert alert-danger"
                }
                role="alert"
              >
                {message}
              </div>
            </div>
          )}
          <CheckButton style={{ display: "none" }} ref={checkBtn} />
        </Form>
  );
};

export default RegisterTutor