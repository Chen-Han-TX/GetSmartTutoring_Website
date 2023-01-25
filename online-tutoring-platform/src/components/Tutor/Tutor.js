import AuthService from "../services/auth.service";
import RideServices from "../services/ride.service";
import React, { useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import Form from "react-validation/build/form";
import Input from "react-validation/build/input";
import CheckButton from "react-validation/build/button";


const Tutor = () => {

  return (
    <div className="auth-inner">
      <header className="passengerPage">
      <h3>Welcome!</h3> 
      <h3> Passenger {currentUser.first_name}</h3>
      </header>
      <br></br>

      {
        !currentRide ? (
          <div className="">
            <h4>Book a ride now!</h4>
       
            <Form onSubmit={handleNewRide} ref={form}>
              <div className="mb-3">
                <label htmlFor="pickup_code">Pickup Postal Code</label>
                <Input
                  type="text"
                  className="form-control"
                  name="pickup_code"
                  value={pickup_code}
                  onChange={onChangePickUpCode}
                  validations={[required, validPostalCode]}
                />
              </div>
    
              <div className="mb-3">
                <label htmlFor="dropoff_code">Dropoff Postal Code</label>
                <Input
                  type="text"
                  className="form-control"
                  name="dropoff_code"
                  value={dropoff_code}
                  onChange={onChangeDropoffCode}
                  validations={[required, validPostalCode]}
                />
              </div>
              
              <div className="d-grid">
              <button className="btn btn-primary btn-block" disabled={loading}>
                {loading && (
                  <span className="spinner-border spinner-border-sm"></span>
                )}
                <span>Search For Rider</span>
              </button>
            </div>

              {message && (
                <div className="form-group">
                  <div className="alert alert-danger" role="alert">
                    {message}
                  </div>
                </div>
              )}
              <CheckButton style={{ display: "none" }} ref={checkBtn} />
            </Form>

          </div>
        ) : (

          <div className="showRidePending">
            <h4>Current Ride Information</h4>
            <p>
              <strong>Ride ID:</strong> {currentRide.ride_id}
            </p>

            {currentRide.ride_status === "Riding" && (
              <div>
                <p>
                <strong>Rider Name: </strong> {currentRide.rider_name}
                </p>
                <p>
                <strong>Rider Phone No.: </strong> {currentRide.rider_phone}
                </p>
              </div>
            )}

            <p>
              <strong>Pickup Postal Code: </strong> {currentRide.pickup_code}
            </p>
            <p>
              <strong>Dropoff Postal Code: </strong> {currentRide.dropoff_code}
            </p>
            <p>
              <strong>Ride Status: </strong> {currentRide.ride_status}
            </p>
            {currentRide.ride_status === "Pending" && (
              <div>
            <p>
              Waiting for rider...
            </p>
            <CancelRideButton ride_id={currentRide.ride_id} />
            </div>
            
            )}
            </div>
        )
      }

    </div>
  );
};

export default Tutor; 