package datamodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const (
	ImageTypeEnumSFSymbol string = "sf_symbol"
	ImageTypeEnumLocal    string = "local"
)

type imageTypeInterface interface {
	Check() UserPresentableErrorInterface
}

var (
	imageTypeRegistry = map[string]func(map[string]interface{}, *Image) (imageTypeInterface, error){
		ImageTypeEnumLocal:    unpackLocalImage,
		ImageTypeEnumSFSymbol: unpackSymbolImage,
	}
)

type Image struct {
	ImageType string
	Height    float64
	Fallback  *Image

	// Section types, stronly typed for easy consumption
	SymbolImageData *SymbolImage
	LocalImageData  *LocalImage

	// generalized interface for functions we need for any section type.
	// Typically a pointer to the one value above that is populated.
	imageData imageTypeInterface
}

type jsonImage struct {
	ImageType      string                 `json:"imageType"`
	Height         *float64               `json:"height,omitempty"`
	Fallback       *Image                 `json:"fallback,omitempty"`
	RawSectionData map[string]interface{} `json:"imageData,omitempty"`
}

func (i *Image) UnmarshalJSON(data []byte) error {
	var ji jsonImage
	err := json.Unmarshal(data, &ji)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an image.", err)
	}

	i.ImageType = ji.ImageType
	i.Fallback = ji.Fallback

	// Default to 50 points
	i.Height = 50.0
	if ji.Height != nil && *ji.Height > 0.0 {
		i.Height = *ji.Height
	}

	unpacker, ok := imageTypeRegistry[i.ImageType]
	if !ok {
		errString := fmt.Sprintf("Unsupported 'imageType' tag in config: \"%v\" is not a valid image type.", i.ImageType)
		if StrictDatamodelParsing {
			return NewUserPresentableError(errString)
		} else {
			// Forward compatibility: warn them the type is unrecognized in debug console, but could be newer config on older build so no hard error
			fmt.Printf("CriticalMoments: %v This image will be ignored. If unexpected, check the CM config file.", errString)
			i.imageData = &UnknownImage{}
		}
	} else {
		imageData, err := unpacker(ji.RawSectionData, i)
		if err != nil {
			return NewUserPresentableErrorWSource("Unknown issue parsing image section of modal page.", err)
		}
		i.imageData = imageData
	}

	return i.Check()
}

func (i *Image) Check() UserPresentableErrorInterface {
	if i.Height <= 0 {
		return NewUserPresentableError("Image height must be > 0")
	}

	if i.imageData == nil {
		return NewUserPresentableError("Invalid image in config -- no imageData map, which is required.")
	}
	if verr := i.imageData.Check(); verr != nil {
		return verr
	}

	if StrictDatamodelParsing && !slices.Contains(maps.Keys(imageTypeRegistry), i.ImageType) {
		return NewUserPresentableError(fmt.Sprintf("Image with invalid 'type' tag: %v", i.ImageType))
	}

	return nil
}

// Unknown image

type UnknownImage struct{}

func (u *UnknownImage) Check() UserPresentableErrorInterface {
	return nil
}

// Local image

type LocalImage struct {
	Path string
}

func unpackLocalImage(data map[string]interface{}, i *Image) (imageTypeInterface, error) {
	path, ok := data["path"].(string)
	if !ok || path == "" {
		return nil, NewUserPresentableError("Image of type local require a path.")
	}

	id := LocalImage{
		Path: path,
	}
	i.LocalImageData = &id

	return &id, nil
}

func (li *LocalImage) Check() UserPresentableErrorInterface {
	if li.Path == "" {
		return NewUserPresentableError("Local images must include a path.")
	}

	return nil
}

// System symbol image

const (
	SystemSymbolWeightEnumUltraLight string = "ultralight"
	SystemSymbolWeightEnumThin       string = "thin"
	SystemSymbolWeightEnumLight      string = "light"
	SystemSymbolWeightEnumRegular    string = "regular"
	SystemSymbolWeightEnumMedium     string = "medium"
	SystemSymbolWeightEnumSemiBold   string = "semibold"
	SystemSymbolWeightEnumBold       string = "bold"
	SystemSymbolWeightEnumHeavy      string = "heavy"
	SystemSymbolWeightEnumBlack      string = "black"

	SystemSymbolModeEnumMono         string = "mono"
	SystemSymbolModeEnumHierarchical string = "hierarchical"
	SystemSymbolModeEnumPalette      string = "palette"
)

var symbolWeights = []string{
	SystemSymbolWeightEnumUltraLight,
	SystemSymbolWeightEnumThin,
	SystemSymbolWeightEnumLight,
	SystemSymbolWeightEnumRegular,
	SystemSymbolWeightEnumMedium,
	SystemSymbolWeightEnumSemiBold,
	SystemSymbolWeightEnumBold,
	SystemSymbolWeightEnumHeavy,
	SystemSymbolWeightEnumBlack,
}

var symbolModes = []string{
	SystemSymbolModeEnumMono,
	SystemSymbolModeEnumHierarchical,
	SystemSymbolModeEnumPalette,
}

type SymbolImage struct {
	SymbolName string
	Weight     string
	Mode       string

	PrimaryColor   string // eg: "#ff0000"
	SecondaryColor string // eg: "#222222"
}

func unpackSymbolImage(data map[string]interface{}, i *Image) (imageTypeInterface, error) {
	symbolName, _ := data["symbolName"].(string)
	primaryColor, _ := data["primaryColor"].(string)
	secondaryColor, _ := data["secondaryColor"].(string)

	weight, _ := data["weight"].(string)
	if !StrictDatamodelParsing && weight != "" && !slices.Contains(symbolWeights, weight) {
		// Back-compat: default to regular
		weight = SystemSymbolWeightEnumRegular
	}

	mode, _ := data["mode"].(string)
	if !StrictDatamodelParsing && mode != "" && !slices.Contains(symbolModes, mode) {
		// Back-compat: default to monocromatic
		mode = SystemSymbolModeEnumMono
	}

	id := SymbolImage{
		SymbolName:     symbolName,
		Weight:         weight,
		Mode:           mode,
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
	}

	if err := id.Check(); err != nil {
		return nil, err
	}

	i.SymbolImageData = &id
	return &id, nil
}

func (si *SymbolImage) Check() UserPresentableErrorInterface {
	if si.SymbolName == "" {
		return NewUserPresentableError("Symbol images must include a symbolName.")
	}

	if si.Weight != "" && !slices.Contains(symbolWeights, si.Weight) {
		// Fallback to default if not strict
		if StrictDatamodelParsing {
			return NewUserPresentableError(fmt.Sprintf("Invalid SF Symbol image 'weight' tag in imageData in config: '%v'", si.Weight))
		}
	}

	if si.Mode != "" && !slices.Contains(symbolModes, si.Mode) {
		// Fallback to default if not strict
		if StrictDatamodelParsing {
			return NewUserPresentableError(fmt.Sprintf("invalid SF Symbold 'mode' tag in imageData in config: '%v'", si.Mode))
		}
	}

	colors := []string{si.PrimaryColor, si.SecondaryColor}
	for _, color := range colors {
		if !stringColorIsValidAllowEmpty(color) {
			return NewUserPresentableError(fmt.Sprintf("Color isn't a valid color. Should be in format '#ffffff' (lower case only). Found \"%v\".", color))
		}
	}

	return nil
}
