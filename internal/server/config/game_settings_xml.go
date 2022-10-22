package config

import (
	"encoding/xml"
	"io"
	"os"
	"path"
)

type GameConfig struct {
	XMLName               xml.Name `xml:"gameSettings"`
	Text                  string   `xml:",chardata"`
	Revision              string   `xml:"revision,attr"`
	ModsDirectoryOverride struct {
		Text      string `xml:",chardata"`
		Active    string `xml:"active,attr"`
		Directory string `xml:"directory,attr"`
	} `xml:"modsDirectoryOverride"`
	DefaultMultiplayerPort    string `xml:"defaultMultiplayerPort"`
	MotorStopTimerDuration    string `xml:"motorStopTimerDuration"`
	HorseAbandonTimerDuration string `xml:"horseAbandonTimerDuration"`
	InvertYLook               string `xml:"invertYLook"`
	IsHeadTrackingEnabled     string `xml:"isHeadTrackingEnabled"`
	ForceFeedback             string `xml:"forceFeedback"`
	IsGamepadEnabled          string `xml:"isGamepadEnabled"`
	CameraSensitivity         string `xml:"cameraSensitivity"`
	VehicleArmSensitivity     string `xml:"vehicleArmSensitivity"`
	RealBeaconLightBrightness string `xml:"realBeaconLightBrightness"`
	SteeringBackSpeed         string `xml:"steeringBackSpeed"`
	SteeringSensitivity       string `xml:"steeringSensitivity"`
	InputHelpMode             string `xml:"inputHelpMode"`
	EasyArmControl            string `xml:"easyArmControl"`
	GyroscopeSteering         string `xml:"gyroscopeSteering"`
	Hints                     string `xml:"hints"`
	CameraTilting             string `xml:"cameraTilting"`
	FrameLimit                string `xml:"frameLimit"`
	ShowAllMods               string `xml:"showAllMods"`
	OnlinePresenceName        string `xml:"onlinePresenceName"`
	Player                    struct {
		Text                string `xml:",chardata"`
		LastPlayerStyleMale string `xml:"lastPlayerStyleMale,attr"`
	} `xml:"player"`
	MpLanguage string `xml:"mpLanguage"`
	CreateGame struct {
		Text             string `xml:",chardata"`
		Password         string `xml:"password,attr"`
		Name             string `xml:"name,attr"`
		Port             string `xml:"port,attr"`
		UseUpnp          string `xml:"useUpnp,attr"`
		AutoAccept       string `xml:"autoAccept,attr"`
		AllowOnlyFriends string `xml:"allowOnlyFriends,attr"`
		AllowCrossPlay   string `xml:"allowCrossPlay,attr"`
		Capacity         string `xml:"capacity,attr"`
		Bandwidth        string `xml:"bandwidth,attr"`
	} `xml:"createGame"`
	Volume struct {
		Text        string `xml:",chardata"`
		Music       string `xml:"music"`
		Vehicle     string `xml:"vehicle"`
		Environment string `xml:"environment"`
		Radio       string `xml:"radio"`
		Gui         string `xml:"gui"`
		Voice       string `xml:"voice"`
		VoiceInput  string `xml:"voiceInput"`
	} `xml:"volume"`
	SoundPlayer struct {
		Text         string `xml:",chardata"`
		AllowStreams string `xml:"allowStreams,attr"`
	} `xml:"soundPlayer"`
	RadioIsActive    string `xml:"radioIsActive"`
	RadioVehicleOnly string `xml:"radioVehicleOnly"`
	Voice            struct {
		Text             string `xml:",chardata"`
		Mode             string `xml:"mode,attr"`
		InputSensitivity string `xml:"inputSensitivity,attr"`
	} `xml:"voice"`
	Units struct {
		Text       string `xml:",chardata"`
		Money      string `xml:"money"`
		Miles      string `xml:"miles"`
		Fahrenheit string `xml:"fahrenheit"`
		Acre       string `xml:"acre"`
	} `xml:"units"`
	IsTrainTabbable        string `xml:"isTrainTabbable"`
	ShowTriggerMarker      string `xml:"showTriggerMarker"`
	ShowHelpTrigger        string `xml:"showHelpTrigger"`
	ShowFieldInfo          string `xml:"showFieldInfo"`
	ShowHelpIcons          string `xml:"showHelpIcons"`
	ShowHelpMenu           string `xml:"showHelpMenu"`
	ResetCamera            string `xml:"resetCamera"`
	ActiveSuspensionCamera string `xml:"activeSuspensionCamera"`
	CameraCheckCollision   string `xml:"cameraCheckCollision"`
	UseWorldCamera         string `xml:"useWorldCamera"`
	IngameMapState         string `xml:"ingameMapState"`
	IngameMapFilters       string `xml:"ingameMapFilters"`
	DirectionChangeMode    string `xml:"directionChangeMode"`
	GearShiftMode          string `xml:"gearShiftMode"`
	HudSpeedGauge          string `xml:"hudSpeedGauge"`
	ShownFreemodeWarning   string `xml:"shownFreemodeWarning"`
	ShowMultiplayerNames   string `xml:"showMultiplayerNames"`
	IngameMapGrowthFilter  string `xml:"ingameMapGrowthFilter"`
	IngameMapSoilFilter    string `xml:"ingameMapSoilFilter"`
	IngameMapFruitFilter   string `xml:"ingameMapFruitFilter"`
	UseColorblindMode      string `xml:"useColorblindMode"`
	MaxNumMirrors          string `xml:"maxNumMirrors"`
	LightsProfile          string `xml:"lightsProfile"`
	FovY                   string `xml:"fovY"`
	UiScale                string `xml:"uiScale"`
	RealBeaconLights       string `xml:"realBeaconLights"`
	CameraBobbing          string `xml:"cameraBobbing"`
}

func GetGameConfig(directory string, file string) (GameConfig, error) {
	x, err := os.Open(path.Join(directory, file))
	if err != nil {
		return GameConfig{}, nil
	}

	b, err := io.ReadAll(x)
	if err != nil {
		return GameConfig{}, nil
	}

	var p GameConfig
	err = xml.Unmarshal(b, &p)
	if err != nil {
		return GameConfig{}, nil
	}

	return p, nil
}
