package appcore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore"
	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

func GoPing() string {
	return "AppcorePong->" + cmcore.CmCorePing()
}

type Appcore struct {
	// Library binding/delegate
	libBindings LibBindings

	// Primary configuration
	configUrlString string
	config          *datamodel.PrimaryConfig
}

var sharedAppcore Appcore = newAppcore()

func SharedAppcore() *Appcore {
	return &sharedAppcore
}
func newAppcore() Appcore {
	return Appcore{}
}

const filePrefix = "file://"
const httpsPrefix = "https://"

func (ac *Appcore) SetConfigUrl(configUrl string) error {
	if !strings.HasPrefix(configUrl, filePrefix) && !strings.HasPrefix(configUrl, httpsPrefix) {
		return errors.New("Config URL must start with https:// or file://")
	}
	ac.configUrlString = configUrl

	return nil
}

func (ac *Appcore) RegisterLibraryBindings(lb LibBindings) {
	ac.libBindings = lb
}

// Might not be "error"
func (ac *Appcore) Start() error {
	if ac.configUrlString == "" {
		return errors.New("A config URL must be provided before starting critical moments")
	}
	if ac.libBindings == nil {
		return errors.New("The SDK must register LibBindings before calling start")
	}

	// Load file:// urls load sync
	if strings.HasPrefix(ac.configUrlString, filePrefix) {
		// Strip file:// prefix
		filePath := ac.configUrlString[len(filePrefix):]
		configFileData, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		var pc datamodel.PrimaryConfig
		err = json.Unmarshal(configFileData, &pc)
		if err != nil {
			return err
		}
		ac.config = &pc
		err = ac.postConfigSetup()
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(ac.configUrlString, httpsPrefix) {
		// TODO: Load http urls async and call postConfigSetup
	}

	return nil
}

func (ac *Appcore) postConfigSetup() error {
	if ac.config.DefaultTheme != nil {
		err := ac.libBindings.SetDefaultTheme(ac.config.DefaultTheme)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: method considered WIP, not tested, expect a full re-write for conditions so saving for later
func (ac *Appcore) SendEvent(e string) {
	actions := ac.config.ActionsForEvent(e)
	for _, action := range actions {
		err := dispatchActionToLib(&action, ac.libBindings)
		if err != nil {
			fmt.Printf("CriticalMoments: there was an issue performing action for event \"%v\". Error: %v", e, err)
		}
	}
}

// TODO: test
func (ac *Appcore) PerformNamedAction(actionName string) error {
	action := ac.config.ActionWithName(actionName)
	if action != nil {
		err := dispatchActionToLib(action, ac.libBindings)
		if err != nil {
			fmt.Printf("CriticalMoments: there was an issue performing action named \"%v\". Error: %v", actionName, err)
			return err
		}
		return nil
	}
	return errors.New(fmt.Sprintf("No action found named %v", actionName))
}

func (ac *Appcore) ThemeForName(themeName string) *datamodel.Theme {
	return ac.config.ThemeWithName(themeName)
}
