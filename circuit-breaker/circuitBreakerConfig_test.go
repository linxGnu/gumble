package cbreaker

import (
	"testing"
	"time"
)

func TestCircuitBreakerConfigValidation(t *testing.T) {
	invalidConf := &CircuitBreakerConfig{}

	if invalidConf.failureRateThreshold = 0; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else if invalidConf.failureRateThreshold = 1.1; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else {
		invalidConf.failureRateThreshold = 0.5
	}

	if invalidConf.trialRequestInterval = 0; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else if invalidConf.trialRequestInterval = -1; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else {
		invalidConf.trialRequestInterval = 1
	}

	if invalidConf.circuitOpenWindow = 0; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else if invalidConf.circuitOpenWindow = -1; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else {
		invalidConf.circuitOpenWindow = 1
	}

	if invalidConf.counterSlidingWindow = 0; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else if invalidConf.counterSlidingWindow = -1; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else {
		invalidConf.counterSlidingWindow = 1
	}

	if invalidConf.counterUpdateInterval = 0; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else if invalidConf.counterUpdateInterval = -1; invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	} else {
		invalidConf.counterUpdateInterval = 1
	}

	if invalidConf.Validate() == nil {
		t.Error("Validation function failed")
	}

	if invalidConf.counterSlidingWindow = 2; invalidConf.Validate() != nil {
		t.Error("Validation function failed")
	}
}

func TestCircuitBreakerConfig(t *testing.T) {
	name := &Name{Name: "dummy-breaker"}
	validConfig := &CircuitBreakerConfig{
		name:                    name,
		failureRateThreshold:    0.7,
		minimumRequestThreshold: 19,
		trialRequestInterval:    time.Second,
		circuitOpenWindow:       time.Second * 2,
		counterSlidingWindow:    time.Second * 3,
		counterUpdateInterval:   time.Second * 4,
		listeners:               make(CircuitBreakerListeners, 2, 10),
	}

	if validConfig.GetName() != name ||
		validConfig.GetFailureRateThreshold() != 0.7 ||
		validConfig.GetMinimumRequestThreshold() != 19 ||
		validConfig.GetTrialRequestInterval() != time.Second ||
		validConfig.GetCircuitOpenWindow() != 2*time.Second ||
		validConfig.GetCounterSlidingWindow() != 3*time.Second ||
		validConfig.GetCounterUpdateInterval() != 4*time.Second ||
		len(validConfig.Getlisteners()) != 2 {
		t.Errorf("Invalid CircuitBreakerConfig")
	}
}
