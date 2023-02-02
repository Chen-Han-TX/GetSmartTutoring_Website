import React, { ChangeEvent, useState, useRef } from "react";
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

import { initializeApp } from "firebase/app";
import { getStorage } from "firebase/storage";
import {ref, uploadBytesResumable, getDownloadURL } from "firebase/storage"

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
  const firebaseConfig = {
    apiKey: "AIzaSyC1fuJ1HYNnv5o2_UauLPx7wwgQGj_Jcpc",
    authDomain: "eti-assignment-2.firebaseapp.com",
    databaseURL: "https://eti-assignment-2-default-rtdb.asia-southeast1.firebasedatabase.app",
    projectId: "eti-assignment-2",
    storageBucket: "eti-assignment-2.appspot.com",
    messagingSenderId: "903055257140",
    appId: "1:903055257140:web:676b8fbb0dfc642dbca461",
    measurementId: "G-986C89B33B"
  };

  const app = initializeApp(firebaseConfig);
  const storage = getStorage(app);
  
  const form = useRef();
  const checkBtn = useRef();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [hourlyRate, setHourlyRate] = useState(50);
  const [startTime, setStartTime] = useState("08:00");
  const [startForEndTime, setStartForEndTime] = useState("09:00");
  const [endTime, setEndTime] = useState("10:00");
  const [availability, setAvailability] = useState({});
  const [selectedSubjects, setSelectedSubjects] = useState({});


  // Maximum 5 files
  const MAX_COUNT = 5;
  const [uploadedFiles, setUploadedFiles] = useState([]);
  const [fileLimit, setFileLimit] = useState(false);
  const [uploadedURLs, setUploadedURLs] = useState([]);


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

  const handleUploadFiles = files => {
    const uploaded = [...uploadedFiles];
    let limitExceeded = false;
    files.some((file) => {
        if (uploaded.findIndex((f) => f.name === file.name) === -1) {
            uploaded.push(file);
            if (uploaded.length === MAX_COUNT) setFileLimit(true);
            if (uploaded.length > MAX_COUNT) {
                alert(`You can only add a maximum of ${MAX_COUNT} files`);
                setFileLimit(false);
                limitExceeded = true;
                return true;
            }
        }
    })
    if (!limitExceeded) setUploadedFiles(uploaded)
  }

  const handleFileEvent =  (e) => {
    const chosenFiles = Array.prototype.slice.call(e.target.files)
    handleUploadFiles(chosenFiles);
}

  const handleRegister = async (e) => {
    e.preventDefault();
    setMessage("");
    setSuccessful(false);

    form.current.validateAll();

    //validate
    if (uploadedFiles.length == 0) {
      alert("Please uploade at least 1 file!");
    }

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
      
      // upload the file to firebase storage
      for (var i=0; i<uploadedFiles.length;i++) {
        const storageRef = ref(storage, `/${email}/${uploadedFiles[i].name}`)
        const uploadTask = uploadBytesResumable(storageRef, uploadedFiles[i]);

        uploadTask.on(
          "state_changed",
          () => {},
          (err) => console.log(err),
          () => {
              getDownloadURL(uploadTask.snapshot.ref).then((url) => {
                const uploaded = uploadedURLs
                uploaded.push(url)
                setUploadedURLs(uploaded);
              });
            })
            if (uploadedURLs.length === uploadedFiles.length){

              alert("Files uploaded successfully!")
              await AuthService.register_tutor(name, email, password, hourlyRate, availability, selectedSubjects, uploadedURLs).then(
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
          }
        }
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
                  name="email"
                  className="form-control"
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
                  className="form-control"
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
                  }}
                    validations= {[required]}
                  />

                  <label>Start</label>
                  <TimePicker 
                      initialValue="08:00"
                      start="00:00"
                      end="23:59"
                      step={30} 
                      value={startTime}
                      onChange={onChangeStartTime}
                      validations={[required]}
                      />
                  <label>End</label>
                  <TimePicker 
                      initialValue="10:00"
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
                <input 
                className="form-control" 
                type="file" 
                multiple
                id="formFileMultiple" 
                onChange={handleFileEvent}
                disabled={fileLimit}
                 />
                <div className="uploaded-files-list">
                    {uploadedFiles.map(file => (
                      <div >
                          {file.name}
                      </div>))}
                  </div>
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