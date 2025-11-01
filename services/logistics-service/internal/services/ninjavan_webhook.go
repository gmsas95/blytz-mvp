package services

import (
	"time"
)

// NinjaVanWebhookPayload represents the structure of Ninja Van webhook payloads
type NinjaVanWebhookPayload struct {
	TrackingID        string `json:"tracking_id"`
	ShipperOrderRefNo string `json:"shipper_order_ref_no"`
	Timestamp         string `json:"timestamp"`
	Event             string `json:"event"`
	Status            string `json:"status"`
	IsParcelOnRTSLeg  bool   `json:"is_parcel_on_rts_leg"`

	// Pickup information
	PickedUpInformation      *NinjaVanPickedUpInfo      `json:"picked_up_information,omitempty"`
	ArrivedAtPUDOInformation *NinjaVanArrivedAtPUDOInfo `json:"arrived_at_pudo_information,omitempty"`
	PickupException          *NinjaVanPickupException   `json:"pickup_exception,omitempty"`

	// Hub information
	ArrivedAtOriginHubInfo        *NinjaVanHubInfo `json:"arrived_at_origin_hub_information,omitempty"`
	ArrivedAtTransitHubInfo       *NinjaVanHubInfo `json:"arrived_at_transit_hub_information,omitempty"`
	ArrivedAtDestinationHubInfo   *NinjaVanHubInfo `json:"arrived_at_destination_hub_information,omitempty"`
	InTransitToNextSortingHubInfo *NinjaVanHubInfo `json:"in_transit_to_next_sorting_hub_information,omitempty"`

	// Delivery information
	OnVehicleInformation *NinjaVanOnVehicleInfo     `json:"on_vehicle_information,omitempty"`
	DeliveryInformation  *NinjaVanDeliveryInfo      `json:"delivery_information,omitempty"`
	DeliveryException    *NinjaVanDeliveryException `json:"delivery_exception,omitempty"`

	// Parcel information
	ParcelMeasurementsInformation *NinjaVanParcelMeasurements `json:"parcel_measurements_information,omitempty"`

	// Cancellation information
	CancellationInformation *NinjaVanCancellationInfo `json:"cancellation_information,omitempty"`

	// Return information
	RecoveryInformation *NinjaVanRecoveryInfo `json:"recovery_information,omitempty"`
	RTSReason           string                `json:"rts_reason,omitempty"`

	// International transit information
	InternationalTransitInformation *NinjaVanInternationalTransit `json:"international_transit_information,omitempty"`
}

type NinjaVanPickedUpInfo struct {
	State string             `json:"state"`
	Proof *NinjaVanProofInfo `json:"proof,omitempty"`
}

type NinjaVanArrivedAtPUDOInfo struct {
	State string             `json:"state"`
	Proof *NinjaVanProofInfo `json:"proof,omitempty"`
}

type NinjaVanPickupException struct {
	State           string             `json:"state"`
	FailureReason   string             `json:"failure_reason"`
	RescheduledDate string             `json:"rescheduled_date"`
	IsLiable        bool               `json:"is_liable"`
	Proof           *NinjaVanProofInfo `json:"proof,omitempty"`
}

type NinjaVanHubInfo struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Hub     string `json:"hub"`
}

type NinjaVanOnVehicleInfo struct {
	AllowDoorstepDropoff bool `json:"allow_doorstep_dropoff"`
}

type NinjaVanDeliveryInfo struct {
	State           string             `json:"state"`
	LeftInSafePlace bool               `json:"left_in_safe_place"`
	Proof           *NinjaVanProofInfo `json:"proof,omitempty"`
}

type NinjaVanDeliveryException struct {
	State           string             `json:"state"`
	FailureReason   string             `json:"failure_reason"`
	RescheduledDate string             `json:"rescheduled_date"`
	IsLiable        bool               `json:"is_liable"`
	Proof           *NinjaVanProofInfo `json:"proof,omitempty"`
}

type NinjaVanParcelMeasurements struct {
	Weight     float64            `json:"weight"`
	Dimensions NinjaVanDimensions `json:"dimensions"`
}

type NinjaVanCancellationInfo struct {
	Reason string `json:"reason"`
}

type NinjaVanRecoveryInfo struct {
	State string `json:"state"`
}

type NinjaVanInternationalTransit struct {
	State               string                `json:"state"`
	Remark              string                `json:"remark"`
	Stage               string                `json:"stage"`
	LinehaulInformation *NinjaVanLinehaulInfo `json:"linehaul_information,omitempty"`
}

type NinjaVanLinehaulInfo struct {
	TransportType            string `json:"transport_type"`
	MasterAWB                string `json:"master_awb"`
	OriginPort               string `json:"origin_port"`
	OriginCountry            string `json:"origin_country"`
	DestinationPort          string `json:"destination_port"`
	DestinationCountry       string `json:"destination_country"`
	EstimatedTimeOfDeparture string `json:"estimated_time_of_departure"`
	EstimatedTimeOfArrival   string `json:"estimated_time_of_arrival"`
	VesselNo                 string `json:"vessel_no"`
}

type NinjaVanProofInfo struct {
	SignatureURI string            `json:"signature_uri"`
	ImageURIs    []string          `json:"image_uris"`
	SignedBy     *NinjaVanSignedBy `json:"signed_by,omitempty"`
}

type NinjaVanSignedBy struct {
	Name         string `json:"name"`
	Contact      string `json:"contact"`
	Relationship string `json:"relationship"`
}
