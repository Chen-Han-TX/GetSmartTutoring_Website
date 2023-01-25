import React, { useState, useRef } from "react";
import Form from "react-validation/build/form";
import Input from "react-validation/build/input";
import CheckButton from "react-validation/build/button";
import { isEmail } from "validator";
import { useNavigate } from "react-router-dom";

import AuthService from "../services/auth.service";

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

function containsOnlyNumbers(str) {
  return /^\d+$/.test(str);
}

const validMobileNumber = (value) => {
  if (value.length !== 8 || !containsOnlyNumbers(value)) {
    return (
      <div className="invalid-feedback d-block">
        Please enter a valid 8 digit phone number
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


  const [email, setEmail] = useState("");
  const [first_name, setFirstName] = useState("");
  const [last_name, setLastName] = useState("");
  const [mobile_number, setMobileNumber] = useState("");
  const [ic_number, setICNumber] = useState("");
  const [car_lic_number, setCarLicNumber] = useState("");
  const [password, setPassword] = useState("");

  const [successful, setSuccessful] = useState(false);
  const [message, setMessage] = useState("");
  const navigate = useNavigate();


  const onChangeEmail = (e) => {
    const email = e.target.value;
    setEmail(email);
  };

  const onChangeFirstName = (e) => {
    const first_name = e.target.value;
    setFirstName(first_name);
  };

  const onChangeLastName = (e) => {
    const last_name = e.target.value;
    setLastName(last_name);
  };

  const onChangeMobileNumber = (e) => {
    const mobile_number = e.target.value;
    setMobileNumber(mobile_number);
  };

  const onChangeIcNumber = (e) => {
    const ic_number = e.target.value;
    setICNumber(ic_number);
  };

  const onChangeCarLicNumber = (e) => {
    const car_lic_number = e.target.value;
    setCarLicNumber(car_lic_number);
  };

  const onChangePassword = (e) => {
    const password = e.target.value;
    setPassword(password);
  };

  const handleRegister = (e) => {
    e.preventDefault();

    setMessage("");
    setSuccessful(false);

    form.current.validateAll();

    // If passed validation, call auth service to send the API request
    if (checkBtn.current.context._errors.length === 0) {
      AuthService.register_rider(email, password, first_name, last_name, mobile_number, ic_number, car_lic_number).then(
        () => {

          setMessage("Registered successfully!");
          setSuccessful(true);
          setTimeout(function () {
            navigate("/login");
            window.location.reload();
          }, 2000);
          

        },
        (error) => {
          var resMessage = ""
          if (error.response.status === 409) {
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
        <Form onSubmit={handleRegister} ref={form}>
          <h3>Register as Rider</h3>
          {!successful && (
            <div>
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
                <label htmlFor="first_name">First name</label>
                <Input
                  type="text"
                  className="form-control"
                  name="first_name"
                  value={first_name}
                  onChange={onChangeFirstName}
                  validations={[required]}
                />
              </div>

              <div className="mb-3">
                <label htmlFor="last_name">Last name</label>
                <Input
                  type="text"
                  className="form-control"
                  name="last_name"
                  value={last_name}
                  onChange={onChangeLastName}
                  validations={[required]}
                />
              </div>

              <div className="mb-3">
                <label htmlFor="mobile_number">Mobile number</label>
                <Input
                  type="text"
                  className="form-control"
                  name="mobile_number"
                  value={mobile_number}
                  onChange={onChangeMobileNumber}
                  validations={[required, validMobileNumber]}
                />
              </div>

              <div className="mb-3">
                <label htmlFor="ic_number">NRIC</label>
                <Input
                  type="text"
                  className="form-control"
                  name="ic_number"
                  value={ic_number}
                  onChange={onChangeIcNumber}
                  validations={[required]}
                />
              </div>
              
              <div className="mb-3">
                <label htmlFor="car_lic_number">Car License Number</label>
                <Input
                  type="text"
                  className="form-control"
                  name="car_lic_number"
                  value={car_lic_number}
                  onChange={onChangeCarLicNumber}
                  validations={[required]}
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

              <div className="d-grid">
                <button className="btn btn-primary btn-block">Register Rider</button>
              </div>
            </div>
          )}

          {message && (
            <div className="form-group">
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